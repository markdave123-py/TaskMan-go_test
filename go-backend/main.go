package main

import (
	"go-backend/server"
	"go-backend/storage"
	"go-backend/storage/db"
	"go-backend/storage/file"
	"log"
	"os"
	"time"

	"go-backend/storage/cache"

	"github.com/sirupsen/logrus"
)

const (
	defaultPort = "8080"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)
}

func main() {
	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Storage driver selection
	driver := os.Getenv("STORAGE_DRIVER")
	if driver == "" {
		driver = "file"
	}

	var store storage.DataStore
	var err error

	switch driver {

	// SQLite storage
	case "sqlite":
		log.Println("using sqlite storage")
		store, err = db.NewSQLiteStore("data/app.db")
		if err != nil {
			log.Fatalf("failed to initialize sqlite store: %v", err)
		}

	// File-based storage
	case "file":
		log.Println("using file storage")
		store, err = file.NewFileStore("data/data.json")
		if err != nil {
			log.Fatalf("failed to initialize file store: %v", err)
		}

	default:
		log.Fatalf("unknown STORAGE_DRIVER: %s", driver)
	}

	// Initialize Redis cache with 2 minute TTL
	cache := cache.NewRedisCache("localhost:6380", 2*time.Minute)

	server := server.NewServer(store, cache)
	server.Start(port)
}
