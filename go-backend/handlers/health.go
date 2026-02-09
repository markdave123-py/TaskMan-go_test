package handlers

import (
	"encoding/json"
	"net/http"

	"go-backend/api"
	"go-backend/storage"
)

// HealthHandler handles HTTP requests related to health checks.
type HealthHandler struct {
	Store storage.DataStore
}

func NewHealthHandler(store storage.DataStore) *HealthHandler {
	return &HealthHandler{Store: store}
}

// Health handles GET requests to /health and returns the health status of the application.
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	resp := api.HealthResponse{
		Status:  "ok",
		Message: "Go backend is running",
	}

	// Readiness check (?ready=true)
	if r.URL.Query().Get("ready") == "true" {
		if h.Store != nil && !h.Store.IsHealthy() {
			w.WriteHeader(http.StatusServiceUnavailable)
			resp.Status = "unhealthy"
			resp.Message = "datastore unavailable"
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
