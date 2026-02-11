package server

import (
	"net/http"

	"go-backend/middleware"
	"go-backend/storage"
	"go-backend/storage/cache"

	"github.com/sirupsen/logrus"
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

	logrus.WithFields(logrus.Fields{
		"component": "server",
		"event":     "startup",
		"port":      port,
		"addr":      "http://localhost:" + port,
	}).Info("HTTP server listening")

	if err := http.ListenAndServe(":"+port, handler); err != nil {
		logrus.WithError(err).WithField("component", "server").Fatal("server failed to start")
	}
}
