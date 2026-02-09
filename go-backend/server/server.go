package server

import (
	"go-backend/middleware"
	"go-backend/storage"
	"go-backend/storage/cache"
	"log"
	"net/http"
)

// Server represents the HTTP server and its dependencies.
type Server struct {
	dataStore storage.DataStore
	cache     *cache.RedisCache
}

func NewServer(dataStore storage.DataStore, cache *cache.RedisCache) *Server {
	return &Server{
		dataStore: dataStore,
		cache:     cache,
	}
}

// Start registers routes and starts the HTTP server on the given port.
// It blocks until the server exits or fails to start.
func (s *Server) Start(port string) {
	s.registerRoutes()

	// Add simple rate limiter middleware
	limiter := middleware.NewRateLimiterStore(1, 10) // 1 req/sec, burst 10

	handler := limiter.Middleware(
		// Add logging middleware
		middleware.LoggingMiddleware(http.DefaultServeMux),
	)

	log.Printf("Go backend server starting on http://localhost:%s", port)
	log.Printf("Serving data directly from Go backend")

	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
