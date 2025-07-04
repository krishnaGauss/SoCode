package main

import (
	"log"
	"net"
	"strconv"

	"github.com/krishnaGauss/SoCode/internal/config"
	"github.com/krishnaGauss/SoCode/internal/server"
	"github.com/krishnaGauss/SoCode/internal/storage"
	"github.com/krishnaGauss/SoCode/proto/SoCode/proto"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()

	postgres, err:=storage.NewPostgresStorage(&cfg.Database)
	if err != nil {
        log.Fatalf("Failed to initialize PostgreSQL: %v", err)
    }
    defer postgres.Close()

	 redis, err := storage.NewRedisQueue(&cfg.Redis)
    if err != nil {
        log.Fatalf("Failed to initialize Redis: %v", err)
    }
    defer redis.Close()

	// Start log processor
    processor := server.NewLogProcessor(redis, postgres)
    go processor.Start()
    defer processor.Stop()

    // Start gRPC server
    lis, err := net.Listen("tcp", ":"+strconv.Itoa(cfg.Server.GRPCPort))
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    logServer := server.NewLogServer(redis, postgres)
    proto.RegisterLogServiceServer(grpcServer, logServer)

    log.Printf("gRPC server listening on port %d", cfg.Server.GRPCPort)
    
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve gRPC: %v", err)
    }
}