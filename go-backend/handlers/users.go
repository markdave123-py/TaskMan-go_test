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

	"github.com/sirupsen/logrus"
)

// UserHandler handles HTTP requests related to users.
type UserHandler struct {
	Store storage.DataStore
	Cache *cache.RedisCache
}

func NewUserHandler(store storage.DataStore, cache *cache.RedisCache) *UserHandler {
	return &UserHandler{Store: store, Cache: cache}

}

// Users handles requests to /api/users.
// GET returns all users.
// POST creates a new user.
func (h *UserHandler) Users(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
	cacheKey := "users:all"

	switch r.Method {

	case http.MethodGet:
		if cached, ok, _ := h.Cache.Get(ctx, cacheKey); ok {
			w.Write(cached)
			return
		}

		users := h.Store.GetUsers()
		resp := api.UsersResponse{
			Users: users,
			Count: len(users),
		}

		data, _ := json.Marshal(resp)
		_ = h.Cache.Set(ctx, cacheKey, data)
		w.Write(data)

	case http.MethodPost:
		h.createUser(w, r)
		_ = h.Cache.Clear(ctx)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// UserByID handles requests to /api/users/{id}.
func (h *UserHandler) UserByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	user, ok := h.Store.GetUserByID(id)
	if !ok {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// createUser handles the logic for creating a new user.
func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	var req api.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Email == "" || req.Role == "" {
		http.Error(w, "all fields are required", http.StatusBadRequest)
		return
	}

	if !util.IsValidEmail(req.Email) {
		http.Error(w, "invalid email format", http.StatusBadRequest)
		return
	}

	if _, exists := h.Store.GetUserByEmail(req.Email); exists {
		http.Error(w, "email already in use", http.StatusBadRequest)
		return
	}

	user, err := h.Store.CreateUser(req.Name, req.Email, req.Role)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"component": "handler",
			"handler":   "users",
			"operation": "create_user",
			"email":     req.Email,
		}).Error("create user failed")
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
