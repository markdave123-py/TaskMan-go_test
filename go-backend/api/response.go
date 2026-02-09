package api

import "go-backend/storage"

// UsersResponse represents the response for the /api/users endpoint.
type UsersResponse struct {
	Users []storage.User `json:"users"`
	Count int            `json:"count"`
}

// TasksResponse represents the response for the /api/tasks endpoint.
type TasksResponse struct {
	Tasks []storage.Task `json:"tasks"`
	Count int            `json:"count"`
}

// StatsResponse represents the response for the /api/stats endpoint.
type StatsResponse struct {
	Users struct {
		Total int `json:"total"`
	} `json:"users"`
	Tasks struct {
		Total      int `json:"total"`
		Pending    int `json:"pending"`
		InProgress int `json:"inProgress"`
		Completed  int `json:"completed"`
	} `json:"tasks"`
}

// HealthResponse represents the response for the health check endpoint.
type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
