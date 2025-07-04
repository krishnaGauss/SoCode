package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/krishnaGauss/SoCode/internal/api"
	"github.com/krishnaGauss/SoCode/internal/config"
	"github.com/krishnaGauss/SoCode/internal/storage"
)

func main() {
    cfg := config.Load()

    // Initialize storage
    postgres, err := storage.NewPostgresStorage(&cfg.Database)
    if err != nil {
        log.Fatalf("Failed to initialize PostgreSQL: %v", err)
    }
    defer postgres.Close()

    // Create API server
    server := api.NewServer(postgres)
    handler := server.SetupRoutes()

    addr := cfg.Server.Host + ":" + strconv.Itoa(cfg.Server.Port)
    log.Printf("API server listening on %s", addr)
    
    if err := http.ListenAndServe(addr, handler); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}