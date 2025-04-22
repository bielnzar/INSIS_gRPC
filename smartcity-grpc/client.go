<<<<<<< HEAD
package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	pb "smartcity/smartcity"

	"google.golang.org/grpc"
)

func main() {
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer conn.Close()
    client := pb.NewSmartCityServiceClient(conn)

    scanner := bufio.NewScanner(os.Stdin)

    for {
        fmt.Println("\nPilih operasi:")
        fmt.Println("1. Get Traffic Status (Unary RPC)")
        fmt.Println("2. Stream Air Quality (Server Streaming RPC)")
        fmt.Println("3. Set Traffic Lights (Client Streaming RPC)")
        fmt.Println("4. Emergency Control (Bidirectional Streaming RPC)")
        fmt.Println("5. Keluar")
        fmt.Print("Masukkan pilihan (1-5): ")

        scanner.Scan()
        choice := strings.TrimSpace(scanner.Text())

        switch choice {
        case "1":
            fmt.Print("Masukkan ID jalan (contoh: Road A): ")
            scanner.Scan()
            roadID := strings.TrimSpace(scanner.Text())
            trafficReq := &pb.TrafficRequest{RoadId: roadID}
            trafficStatus, err := client.GetTrafficStatus(context.Background(), trafficReq)
            if err != nil {
                log.Printf("Error getting traffic status: %v", err)
                continue
            }
            fmt.Printf("Traffic Status: %s - %s, Vehicles: %d\n", trafficStatus.RoadId, trafficStatus.CongestionLevel, trafficStatus.VehicleCount)

        case "2":
            fmt.Print("Masukkan ID zona (contoh: Downtown): ")
            scanner.Scan()
            zoneID := strings.TrimSpace(scanner.Text())
            airReq := &pb.AirQualityRequest{ZoneId: zoneID}
            stream, err := client.StreamAirQuality(context.Background(), airReq)
            if err != nil {
                log.Printf("Error streaming air quality: %v", err)
                continue
            }
            fmt.Println("Menerima data kualitas udara...")
            for {
                data, err := stream.Recv()
                if err == io.EOF {
                    break
                }
                if err != nil {
                    log.Printf("Error receiving stream: %v", err)
                    break
                }
                fmt.Printf("Air Quality in %s: %.2f at %s\n", data.ZoneId, data.PollutionLevel, data.Timestamp)
            }

        case "3":
            fmt.Println("Masukkan perintah lampu lalu lintas (format: intersection_id,action). Ketik 'done' untuk selesai.")
            var commands []*pb.TrafficLightCommand
            for {
                fmt.Print("Perintah: ")
                scanner.Scan()
                input := strings.TrimSpace(scanner.Text())
                if input == "done" {
                    break
                }
                parts := strings.SplitN(input, ",", 2)
                if len(parts) != 2 {
                    fmt.Println("Format salah. Gunakan: intersection_id,action")
                    continue
                }
                cmd := &pb.TrafficLightCommand{
                    IntersectionId: strings.TrimSpace(parts[0]),
                    Action:         strings.TrimSpace(parts[1]),
                }
                commands = append(commands, cmd)
            }
            if len(commands) == 0 {
                fmt.Println("Tidak ada perintah yang dikirim.")
                continue
            }
            cmdStream, err := client.SetTrafficLights(context.Background())
            if err != nil {
                log.Printf("Error setting traffic lights: %v", err)
                continue
            }
            for _, cmd := range commands {
                if err := cmdStream.Send(cmd); err != nil {
                    log.Printf("Error sending command: %v", err)
                    continue
                }
            }
            resp, err := cmdStream.CloseAndRecv()
            if err != nil {
                log.Printf("Error closing stream: %v", err)
                continue
            }
            fmt.Println(resp.Message)

        case "4":
            fmt.Println("Masukkan perintah darurat (format: unit_id,command).")
            fmt.Println("Contoh: Ambulance 01,Move to Zone B")
            fmt.Println("Perintah lain: check_status, priority_mode")
            fmt.Println("Ketik 'done' untuk selesai.")
            emerStream, err := client.EmergencyControl(context.Background())
            if err != nil {
                log.Printf("Error in emergency control: %v", err)
                continue
            }
            firstCommand := true // Untuk menunda tabel hingga perintah pertama
            go func() {
                for {
                    feedback, err := emerStream.Recv()
                    if err == io.EOF {
                        break
                    }
                    if err != nil {
                        log.Printf("Error receiving feedback: %v", err)
                        break
                    }
                    if firstCommand {
                        fmt.Println("\n--- Emergency Feedback ---")
                        fmt.Printf("| %-15s | %-40s | %-15s |\n", "Unit ID", "Status", "Location")
                        fmt.Println("|----------------|------------------------------------------|-----------------|")
                        firstCommand = false
                    }
                    status := feedback.Status
                    if len(status) > 40 {
                        status = status[:37] + "..."
                    }
                    fmt.Printf("| %-15s | %-40s | %-15s |\n", feedback.UnitId, status, feedback.Location)
                }
            }()
            for {
                fmt.Print("Perintah darurat: ")
                scanner.Scan()
                input := strings.TrimSpace(scanner.Text())
                if input == "done" {
                    break
                }
                parts := strings.SplitN(input, ",", 2)
                if len(parts) != 2 {
                    fmt.Println("Format salah. Gunakan: unit_id,command")
                    continue
                }
                cmd := &pb.EmergencyCommand{
                    UnitId:  strings.TrimSpace(parts[0]),
                    Command: strings.TrimSpace(parts[1]),
                }
                if err := emerStream.Send(cmd); err != nil {
                    log.Printf("Error sending emergency command: %v", err)
                    continue
                }
                time.Sleep(500 * time.Millisecond)
            }
            emerStream.CloseSend()
            time.Sleep(2 * time.Second)
            fmt.Println("-----------------------")

        case "5":
            fmt.Println("Keluar dari program.")
            return

        default:
            fmt.Println("Pilihan tidak valid. Silakan pilih antara 1-5.")
        }
    }
=======
package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	pb "smartcity/smartcity"

	"google.golang.org/grpc"
)

func main() {
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer conn.Close()
    client := pb.NewSmartCityServiceClient(conn)

    scanner := bufio.NewScanner(os.Stdin)

    for {
        fmt.Println("\nPilih operasi:")
        fmt.Println("1. Get Traffic Status (Unary RPC)")
        fmt.Println("2. Stream Air Quality (Server Streaming RPC)")
        fmt.Println("3. Set Traffic Lights (Client Streaming RPC)")
        fmt.Println("4. Emergency Control (Bidirectional Streaming RPC)")
        fmt.Println("5. Keluar")
        fmt.Print("Masukkan pilihan (1-5): ")

        scanner.Scan()
        choice := strings.TrimSpace(scanner.Text())

        switch choice {
        case "1":
            fmt.Print("Masukkan ID jalan (contoh: Road A): ")
            scanner.Scan()
            roadID := strings.TrimSpace(scanner.Text())
            trafficReq := &pb.TrafficRequest{RoadId: roadID}
            trafficStatus, err := client.GetTrafficStatus(context.Background(), trafficReq)
            if err != nil {
                log.Printf("Error getting traffic status: %v", err)
                continue
            }
            fmt.Printf("Traffic Status: %s - %s, Vehicles: %d\n", trafficStatus.RoadId, trafficStatus.CongestionLevel, trafficStatus.VehicleCount)

        case "2":
            fmt.Print("Masukkan ID zona (contoh: Downtown): ")
            scanner.Scan()
            zoneID := strings.TrimSpace(scanner.Text())
            airReq := &pb.AirQualityRequest{ZoneId: zoneID}
            stream, err := client.StreamAirQuality(context.Background(), airReq)
            if err != nil {
                log.Printf("Error streaming air quality: %v", err)
                continue
            }
            fmt.Println("Menerima data kualitas udara...")
            for {
                data, err := stream.Recv()
                if err == io.EOF {
                    break
                }
                if err != nil {
                    log.Printf("Error receiving stream: %v", err)
                    break
                }
                fmt.Printf("Air Quality in %s: %.2f at %s\n", data.ZoneId, data.PollutionLevel, data.Timestamp)
            }

        case "3":
            fmt.Println("Masukkan perintah lampu lalu lintas (format: intersection_id,action). Ketik 'done' untuk selesai.")
            var commands []*pb.TrafficLightCommand
            for {
                fmt.Print("Perintah: ")
                scanner.Scan()
                input := strings.TrimSpace(scanner.Text())
                if input == "done" {
                    break
                }
                parts := strings.SplitN(input, ",", 2)
                if len(parts) != 2 {
                    fmt.Println("Format salah. Gunakan: intersection_id,action")
                    continue
                }
                cmd := &pb.TrafficLightCommand{
                    IntersectionId: strings.TrimSpace(parts[0]),
                    Action:         strings.TrimSpace(parts[1]),
                }
                commands = append(commands, cmd)
            }
            if len(commands) == 0 {
                fmt.Println("Tidak ada perintah yang dikirim.")
                continue
            }
            cmdStream, err := client.SetTrafficLights(context.Background())
            if err != nil {
                log.Printf("Error setting traffic lights: %v", err)
                continue
            }
            for _, cmd := range commands {
                if err := cmdStream.Send(cmd); err != nil {
                    log.Printf("Error sending command: %v", err)
                    continue
                }
            }
            resp, err := cmdStream.CloseAndRecv()
            if err != nil {
                log.Printf("Error closing stream: %v", err)
                continue
            }
            fmt.Println(resp.Message)

        case "4":
            fmt.Println("Masukkan perintah darurat (format: unit_id,command).")
            fmt.Println("Contoh: Ambulance 01,Move to Zone B")
            fmt.Println("Perintah lain: check_status, priority_mode")
            fmt.Println("Ketik 'done' untuk selesai.")
            emerStream, err := client.EmergencyControl(context.Background())
            if err != nil {
                log.Printf("Error in emergency control: %v", err)
                continue
            }
            firstCommand := true // Untuk menunda tabel hingga perintah pertama
            go func() {
                for {
                    feedback, err := emerStream.Recv()
                    if err == io.EOF {
                        break
                    }
                    if err != nil {
                        log.Printf("Error receiving feedback: %v", err)
                        break
                    }
                    if firstCommand {
                        fmt.Println("\n--- Emergency Feedback ---")
                        fmt.Printf("| %-15s | %-40s | %-15s |\n", "Unit ID", "Status", "Location")
                        fmt.Println("|----------------|------------------------------------------|-----------------|")
                        firstCommand = false
                    }
                    status := feedback.Status
                    if len(status) > 40 {
                        status = status[:37] + "..."
                    }
                    fmt.Printf("| %-15s | %-40s | %-15s |\n", feedback.UnitId, status, feedback.Location)
                }
            }()
            for {
                fmt.Print("Perintah darurat: ")
                scanner.Scan()
                input := strings.TrimSpace(scanner.Text())
                if input == "done" {
                    break
                }
                parts := strings.SplitN(input, ",", 2)
                if len(parts) != 2 {
                    fmt.Println("Format salah. Gunakan: unit_id,command")
                    continue
                }
                cmd := &pb.EmergencyCommand{
                    UnitId:  strings.TrimSpace(parts[0]),
                    Command: strings.TrimSpace(parts[1]),
                }
                if err := emerStream.Send(cmd); err != nil {
                    log.Printf("Error sending emergency command: %v", err)
                    continue
                }
                time.Sleep(500 * time.Millisecond)
            }
            emerStream.CloseSend()
            time.Sleep(2 * time.Second)
            fmt.Println("-----------------------")

        case "5":
            fmt.Println("Keluar dari program.")
            return

        default:
            fmt.Println("Pilihan tidak valid. Silakan pilih antara 1-5.")
        }
    }
>>>>>>> e6016d1 (Menambahkan fitur streaming dan memperbaiki tampilan UI)
}