package main

import (
	"log"
	"os"
	"context"
	"time"

	"github.com/krishnaGauss/SoCode/proto/SoCode/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: agent <server-address>")
	}

	serverAddr := os.Args[1]

	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := proto.NewLogServiceClient(conn)

	// Send some sample logs
	logs := []*proto.LogRequest{
		{
			Timestamp: timestamppb.Now(),
			Level:     "INFO",
			Message:   "Application started successfully",
			Source:    "app.log",
			Service:   "web-service",
			Host:      "server-01",
			Tags:      map[string]string{"env": "production", "version": "1.0.0"},
		},
		{
			Timestamp: timestamppb.Now(),
			Level:     "ERROR",
			Message:   "Database connection failed",
			Source:    "app.log",
			Service:   "web-service",
			Host:      "server-01",
			Tags:      map[string]string{"env": "production", "component": "database"},
		},
		{
			Timestamp: timestamppb.Now(),
			Level:     "WARN",
			Message:   "High memory usage detected",
			Source:    "system.log",
			Service:   "monitoring",
			Host:      "server-01",
			Tags:      map[string]string{"env": "production", "metric": "memory"},
		},
	}

	for _, logReq := range logs {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		resp, err := client.SendLog(ctx, logReq)
		cancel()

		if err != nil {
			log.Printf("Failed to send log: %v", err)
			continue
		}

		if resp.Success {
			log.Printf("Log sent successfully: %s", resp.Message)
		} else {
			log.Printf("Failed to send log: %s", resp.Message)
		}
	}
}