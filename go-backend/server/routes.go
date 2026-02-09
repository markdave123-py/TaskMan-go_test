package server

import (
	"net/http"

	"go-backend/handlers"
)

// registerRoutes sets up the HTTP routes and their corresponding handlers.
func (s *Server) registerRoutes() {

	userHandler := handlers.NewUserHandler(s.dataStore, s.cache)
	taskHandler := handlers.NewTaskHandler(s.dataStore, s.cache)
	healthHandler := handlers.NewHealthHandler(s.dataStore)
	statsHandler := handlers.NewStatsHandler(s.dataStore)

	http.HandleFunc("/health", healthHandler.Health)

	http.HandleFunc("/api/users", userHandler.Users)
	http.HandleFunc("/api/users/", userHandler.UserByID)

	http.HandleFunc("/api/tasks", taskHandler.Tasks)
	http.HandleFunc("/api/tasks/", taskHandler.TaskByID)

	http.HandleFunc("/api/stats", statsHandler.Stats)
}
