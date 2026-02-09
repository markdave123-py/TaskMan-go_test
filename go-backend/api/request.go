package api

// CreateUserRequest represents the request body for creating a new user.
type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// CreateTaskRequest represents the request body for creating a new task.
type CreateTaskRequest struct {
	Title  string `json:"title"`
	Status string `json:"status"`
	UserID int    `json:"userId"`
}

// UpdateTaskRequest represents the request body for updating an existing task.
type UpdateTaskRequest struct {
	Title  *string `json:"title"`
	Status *string `json:"status"`
	UserID *int    `json:"userId"`
}
