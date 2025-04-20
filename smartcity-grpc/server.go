package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"

	pb "smartcity/smartcity"

	"google.golang.org/grpc"
)

type unitState struct {
    position   string // Lokasi saat ini: "Pusat", "Perjalanan", atau tujuan
    mission    string // Tujuan misi (misalnya, "Lokasi Bencana")
    travelTime int    // Waktu tersisa dalam detik
    active     bool   // Menandakan apakah unit sedang diproses
}

type server struct {
    pb.UnimplementedSmartCityServiceServer
    units map[string]*unitState // Melacak status unit
    mu    sync.Mutex
}

func (s *server) GetTrafficStatus(ctx context.Context, req *pb.TrafficRequest) (*pb.TrafficStatus, error) {
    congestion := []string{"Rendah", "Sedang", "Tinggi"}
    return &pb.TrafficStatus{
        RoadId:          req.RoadId,
        CongestionLevel: congestion[rand.Intn(3)],
        VehicleCount:    rand.Int31n(100),
    }, nil
}

func (s *server) StreamAirQuality(req *pb.AirQualityRequest, stream pb.SmartCityService_StreamAirQualityServer) error {
    for i := 0; i < 5; i++ {
        data := &pb.AirQualityData{
            ZoneId:         req.ZoneId,
            PollutionLevel: rand.Float32() * 100,
            Timestamp:      time.Now().Format(time.RFC3339),
        }
        if err := stream.Send(data); err != nil {
            return err
        }
        time.Sleep(1 * time.Second)
    }
    return nil
}

func (s *server) SetTrafficLights(stream pb.SmartCityService_SetTrafficLightsServer) error {
    for {
        cmd, err := stream.Recv()
        if err != nil {
            if err == io.EOF {
                return stream.SendAndClose(&pb.CommandResponse{Message: "Semua perintah telah diproses"})
            }
            return err
        }
        fmt.Printf("Menerima perintah untuk %s: %s\n", cmd.IntersectionId, cmd.Action)
    }
}

func (s *server) EmergencyControl(stream pb.SmartCityService_EmergencyControlServer) error {
    unitID := ""
    stopChan := make(chan struct{})
    s.mu.Lock()
    if s.units == nil {
        s.units = make(map[string]*unitState)
    }
    s.mu.Unlock()

    // Goroutine untuk mengirim pembaruan posisi
    go func() {
        idleSent := false
        for {
            select {
            case <-stopChan:
                return
            default:
                s.mu.Lock()
                if state, exists := s.units[unitID]; exists && state.active {
                    if state.mission == "" && !idleSent {
                        feedback := &pb.EmergencyFeedback{
                            UnitId:   unitID,
                            Status:   "Misi saat ini: Idle",
                            Location: state.position,
                        }
                        if err := stream.Send(feedback); err != nil {
                            s.mu.Unlock()
                            return
                        }
                        idleSent = true
                    } else if state.mission != "" {
                        idleSent = false
                        if state.position == "Pusat" {
                            // Pindah ke status Perjalanan
                            state.position = "Perjalanan"
                            feedback := &pb.EmergencyFeedback{
                                UnitId:   unitID,
                                Status:   fmt.Sprintf("Unit dalam perjalanan menuju %s, waktu tersisa: %ds", state.mission, state.travelTime),
                                Location: state.position,
                            }
                            if err := stream.Send(feedback); err != nil {
                                s.mu.Unlock()
                                return
                            }
                        } else if state.position == "Perjalanan" {
                            // Update status selama perjalanan
                            feedback := &pb.EmergencyFeedback{
                                UnitId:   unitID,
                                Status:   fmt.Sprintf("Unit dalam perjalanan menuju %s, waktu tersisa: %ds", state.mission, state.travelTime),
                                Location: state.position,
                            }
                            if err := stream.Send(feedback); err != nil {
                                s.mu.Unlock()
                                return
                            }
                            if state.travelTime > 0 {
                                state.travelTime--
                                if state.travelTime == 0 {
                                    // Sampai di tujuan
                                    state.position = state.mission
                                    state.mission = ""
                                    stream.Send(&pb.EmergencyFeedback{
                                        UnitId:   unitID,
                                        Status:   "Unit telah tiba di tujuan",
                                        Location: state.position,
                                    })
                                    idleSent = false
                                }
                            }
                        }
                    }
                }
                s.mu.Unlock()
                time.Sleep(1 * time.Second)
            }
        }
    }()

    for {
        cmd, err := stream.Recv()
        if err != nil {
            if err == io.EOF {
                close(stopChan)
                return nil
            }
            close(stopChan)
            return err
        }
        unitID = cmd.UnitId
        fmt.Printf("Menerima perintah darurat untuk %s: %s\n", cmd.UnitId, cmd.Command)

        s.mu.Lock()
        if _, exists := s.units[unitID]; !exists {
            s.units[unitID] = &unitState{position: "Pusat", mission: "", active: true}
        }
        state := s.units[unitID]
        state.active = true
        s.mu.Unlock()

        switch strings.ToLower(cmd.Command) {
        case "check_status":
            mission := state.mission
            if mission == "" {
                mission = "Idle"
            }
            status := "Misi saat ini: " + mission
            if state.position == "Perjalanan" {
                status = fmt.Sprintf("Unit dalam perjalanan menuju %s, waktu tersisa: %ds", state.mission, state.travelTime)
            }
            feedback := &pb.EmergencyFeedback{
                UnitId:   unitID,
                Status:   status,
                Location: state.position,
            }
            if err := stream.Send(feedback); err != nil {
                close(stopChan)
                return err
            }
        case "priority_mode":
            if state.mission == "" {
                feedback := &pb.EmergencyFeedback{
                    UnitId:   unitID,
                    Status:   "Tidak dapat mengaktifkan mode prioritas: Tidak ada misi aktif",
                    Location: state.position,
                }
                if err := stream.Send(feedback); err != nil {
                    close(stopChan)
                    return err
                }
            } else {
                feedback := &pb.EmergencyFeedback{
                    UnitId:   unitID,
                    Status:   "Mode prioritas diaktifkan, waktu tempuh dipercepat",
                    Location: state.position,
                }
                state.travelTime = state.travelTime / 2
                if err := stream.Send(feedback); err != nil {
                    close(stopChan)
                    return err
                }
            }
        default:
            // Parse perintah untuk mengekstrak tujuan
            destination := cmd.Command
            if strings.HasPrefix(strings.ToLower(cmd.Command), "kirim bantuan ke ") {
                destination = strings.TrimSpace(cmd.Command[17:])
            } else if strings.HasPrefix(strings.ToLower(cmd.Command), "evakuasi ") {
                destination = strings.TrimSpace(cmd.Command[9:])
            }
            travelTime := 60 // Waktu perjalanan 60 detik
            s.mu.Lock()
            state.mission = destination
            state.travelTime = travelTime
            state.position = "Pusat" // Mulai dari Pusat
            s.mu.Unlock()
            feedback := &pb.EmergencyFeedback{
                UnitId:   unitID,
                Status:   fmt.Sprintf("Unit memulai perjalanan ke %s, perkiraan waktu: %ds", destination, travelTime),
                Location: state.position,
            }
            if err := stream.Send(feedback); err != nil {
                close(stopChan)
                return err
            }
        }
    }
}

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Gagal mendengarkan: %v", err)
    }
    s := grpc.NewServer()
    pb.RegisterSmartCityServiceServer(s, &server{})
    log.Println("Server dimulai di port 50051")
    if err := s.Serve(lis); err != nil {
        log.Fatalf("Gagal melayani: %v", err)
    }
}