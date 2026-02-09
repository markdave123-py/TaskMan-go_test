package handlers

import (
	"encoding/json"
	"net/http"

	"go-backend/storage"
)

type StatsHandler struct {
	store storage.DataStore
}

func NewStatsHandler(store storage.DataStore) *StatsHandler {
	return &StatsHandler{store: store}
}

func (h *StatsHandler) Stats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	stats := h.store.GetStats()
	json.NewEncoder(w).Encode(stats)
}
