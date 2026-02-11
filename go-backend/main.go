package main

import (
	"os"
	"time"

	"go-backend/server"
	"go-backend/storage"
	"go-backend/storage/cache"
	"go-backend/storage/db"
	"go-backend/storage/file"

	"github.com/sirupsen/logrus"
)

const (
	defaultPort = "8080"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyMsg: "message",
		},
	})
	logrus.SetLevel(logrus.InfoLevel)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	driver := os.Getenv("STORAGE_DRIVER")
	if driver == "" {
		driver = "file"
	}

	log := logrus.WithFields(logrus.Fields{
		"component": "main",
		"port":      port,
		"driver":    driver,
	})

	var store storage.DataStore
	var err error

	switch driver {
	case "sqlite":
		log.Info("initializing SQLite storage")
		store, err = db.NewSQLiteStore("data/app.db")
		if err != nil {
			log.WithError(err).Fatal("failed to initialize SQLite store")
		}
	case "file":
		log.Info("initializing file storage")
		store, err = file.NewFileStore("data/data.json")
		if err != nil {
			log.WithError(err).Fatal("failed to initialize file store")
		}
	default:
		log.WithField("driver", driver).Fatal("unknown STORAGE_DRIVER")
	}

	log.Info("storage initialized successfully")

	cacheSvc := cache.NewRedisCache("localhost:6380", 2*time.Minute)
	srv := server.NewServer(store, cacheSvc)
	srv.Start(port)
}
