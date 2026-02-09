package main

import (
	"go-backend/server"
	"go-backend/storage/file"
	"log"
	"os"
	"time"

	"go-backend/storage/cache"
)

const (
	defaultPort = "8080"
)

func main() {
	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Initialize file-based data store
	store, err := file.NewFileStore("data/data.json")
	if err != nil {
		log.Fatalf("failed to init store: %v", err)
	}

	// Initialize Redis cache with 2 minute TTL
	cache := cache.NewRedisCache("localhost:6380", 2*time.Minute)

	server := server.NewServer(store, cache)
	server.Start(port)
}
