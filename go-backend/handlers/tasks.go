package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"go-backend/api"
	"go-backend/storage"
	"go-backend/storage/cache"
	"go-backend/util"
)

// TaskHandler handles HTTP requests related to tasks.
type TaskHandler struct {
	Store storage.DataStore
	Cache *cache.RedisCache
}

func NewTaskHandler(store storage.DataStore, cache *cache.RedisCache) *TaskHandler {
	return &TaskHandler{
		Store: store,
		Cache: cache,
	}
}

// Tasks handles requests to /api/tasks.
func (h *TaskHandler) Tasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
	cacheKey := "tasks:" + r.URL.RawQuery

	switch r.Method {

	case http.MethodGet:
		if cached, ok, _ := h.Cache.Get(ctx, cacheKey); ok {
			w.Write(cached)
			return
		}

		status := r.URL.Query().Get("status")

		var userID *int
		if userIDStr := r.URL.Query().Get("userId"); userIDStr != "" {
			id, err := strconv.Atoi(userIDStr)
			if err != nil {
				http.Error(w, "invalid userId", http.StatusBadRequest)
				return
			}
			userID = &id
		}

		tasks := h.Store.GetTasks(status, userID)
		resp := api.TasksResponse{
			Tasks: tasks,
			Count: len(tasks),
		}

		data, _ := json.Marshal(resp)
		_ = h.Cache.Set(ctx, cacheKey, data)
		w.Write(data)

	case http.MethodPost:
		h.createTask(w, r)
		_ = h.Cache.Clear(ctx)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// TaskByID handles requests to /api/tasks/{id}.
func (h *TaskHandler) TaskByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	idStr := strings.TrimPrefix(r.URL.Path, "/api/tasks/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	var req api.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	// Validate status if provided
	if req.Status != nil && !util.IsValidTaskStatus(*req.Status) {
		http.Error(w, "invalid task status", http.StatusBadRequest)
		return
	}

	// Validate userId if provided
	if req.UserID != nil {
		if _, ok := h.Store.GetUserByID(*req.UserID); !ok {
			http.Error(w, "user not found", http.StatusBadRequest)
			return
		}
	}

	updatedTask, found, err := h.Store.UpdateTask(
		id,
		req.Title,
		req.Status,
		req.UserID,
	)

	if err != nil {
		http.Error(w, "failed to update task", http.StatusInternalServerError)
		return
	}

	if !found {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	_ = h.Cache.Clear(ctx)
	json.NewEncoder(w).Encode(updatedTask)
}

// createTask handles the logic for creating a new task.
func (h *TaskHandler) createTask(w http.ResponseWriter, r *http.Request) {
	var req api.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if req.Title == "" || req.Status == "" || req.UserID == 0 {
		http.Error(w, "all fields are required", http.StatusBadRequest)
		return
	}

	if !util.IsValidTaskStatus(req.Status) {
		http.Error(w, "invalid task status", http.StatusBadRequest)
		return
	}

	task, err := h.Store.CreateTask(req.Title, req.Status, req.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// getTasks handles the logic for retrieving tasks with optional filters.
func (h *TaskHandler) getTasks(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")

	var userID *int
	if userIDStr := r.URL.Query().Get("userId"); userIDStr != "" {
		id, err := strconv.Atoi(userIDStr)
		if err != nil {
			http.Error(w, "invalid userId", http.StatusBadRequest)
			return
		}
		userID = &id
	}

	tasks := h.Store.GetTasks(status, userID)

	resp := api.TasksResponse{
		Tasks: tasks,
		Count: len(tasks),
	}

	json.NewEncoder(w).Encode(resp)
}
