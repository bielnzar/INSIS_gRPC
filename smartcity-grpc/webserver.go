package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	pb "smartcity/smartcity"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
    // Atur mode release untuk Gin
    gin.SetMode(gin.ReleaseMode)

    // Koneksi ke server gRPC
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Failed to connect to gRPC server: %v", err)
    }
    defer conn.Close()
    client := pb.NewSmartCityServiceClient(conn)

    // Inisialisasi router Gin
    r := gin.Default()

    // Atur trusted proxies
    if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
        log.Fatalf("Failed to set trusted proxies: %v", err)
    }

    // Rute untuk file statis (index.html)
    r.Static("/static", "./static")
    r.GET("/", func(c *gin.Context) {
        c.File("./static/index.html")
    })

    // Unary RPC
    r.POST("/unary", func(c *gin.Context) {
        var req struct {
            RoadID string `json:"road_id"`
        }
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        resp, err := client.GetTrafficStatus(context.Background(), &pb.TrafficRequest{RoadId: req.RoadID})
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{
            "road_id":         resp.RoadId,
            "congestion_level": resp.CongestionLevel,
            "vehicle_count":   resp.VehicleCount,
        })
    })

    // Server Streaming RPC
    r.GET("/server-stream", func(c *gin.Context) {
        zoneID := c.Query("zone_id")
        ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
        if err != nil {
            log.Printf("WebSocket upgrade error: %v", err)
            return
        }
        defer ws.Close()

        stream, err := client.StreamAirQuality(context.Background(), &pb.AirQualityRequest{ZoneId: zoneID})
        if err != nil {
            ws.WriteJSON(gin.H{"error": err.Error()})
            return
        }

        for {
            data, err := stream.Recv()
            if err == io.EOF {
                break
            }
            if err != nil {
                ws.WriteJSON(gin.H{"error": err.Error()})
                break
            }
            ws.WriteJSON(gin.H{
                "zone_id":        data.ZoneId,
                "pollution_level": data.PollutionLevel,
                "timestamp":      data.Timestamp,
            })
            time.Sleep(100 * time.Millisecond)
        }
    })

    // Client Streaming RPC
    r.GET("/client-stream", func(c *gin.Context) {
        ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
        if err != nil {
            log.Printf("WebSocket upgrade error: %v", err)
            return
        }
        defer ws.Close()

        stream, err := client.SetTrafficLights(context.Background())
        if err != nil {
            ws.WriteJSON(gin.H{"error": err.Error()})
            return
        }

        for {
            _, message, err := ws.ReadMessage()
            if err != nil {
                stream.CloseAndRecv()
                return
            }
            var cmd struct {
                IntersectionID string `json:"intersection_id"`
                Action         string `json:"action"`
                Done           bool   `json:"done"`
            }
            if err := json.Unmarshal(message, &cmd); err != nil {
                ws.WriteJSON(gin.H{"error": "Format perintah tidak valid"})
                continue
            }

            if cmd.Done {
                resp, err := stream.CloseAndRecv()
                if err != nil {
                    ws.WriteJSON(gin.H{"error": err.Error()})
                    return
                }
                ws.WriteJSON(gin.H{"message": resp.Message})
                return
            }

            if err := stream.Send(&pb.TrafficLightCommand{
                IntersectionId: cmd.IntersectionID,
                Action:         cmd.Action,
            }); err != nil {
                ws.WriteJSON(gin.H{"error": err.Error()})
                return
            }
        }
    })

    // Bidirectional Streaming RPC
    r.GET("/bidi-stream", func(c *gin.Context) {
        ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
        if err != nil {
            log.Printf("WebSocket upgrade error: %v", err)
            return
        }
        defer ws.Close()

        stream, err := client.EmergencyControl(context.Background())
        if err != nil {
            ws.WriteJSON(gin.H{"error": err.Error()})
            return
        }

        // Goroutine untuk menerima feedback dari server
        go func() {
            for {
                feedback, err := stream.Recv()
                if err == io.EOF {
                    ws.WriteJSON(gin.H{"message": "Stream ditutup"})
                    ws.Close()
                    return
                }
                if err != nil {
                    ws.WriteJSON(gin.H{"error": err.Error()})
                    ws.Close()
                    return
                }
                ws.WriteJSON(gin.H{
                    "unit_id":  feedback.UnitId,
                    "status":   feedback.Status,
                    "location": feedback.Location,
                })
            }
        }()

        // Menerima perintah dari WebSocket
        for {
            _, message, err := ws.ReadMessage()
            if err != nil {
                stream.CloseSend()
                return
            }
            var cmd struct {
                UnitID  string `json:"unit_id"`
                Command string `json:"command"`
            }
            if err := json.Unmarshal(message, &cmd); err != nil {
                ws.WriteJSON(gin.H{"error": "Format perintah tidak valid"})
                continue
            }
            if err := stream.Send(&pb.EmergencyCommand{
                UnitId:  cmd.UnitID,
                Command: cmd.Command,
            }); err != nil {
                ws.WriteJSON(gin.H{"error": err.Error()})
                return
            }
        }
    })

    log.Println("Web server started on :8081")
    if err := r.Run(":8081"); err != nil {
        log.Fatalf("Failed to start web server: %v", err)
    }
}