package storage

type DataStore interface {
	// Users
	GetUsers() []User
	GetUserByID(id int) (*User, bool)
	CreateUser(name, email, role string) (User, error)

	// Tasks
	GetTasks(status string, userID *int) []Task
	CreateTask(title, status string, userID int) (Task, error)
	UpdateTask(id int, title *string, status *string, userID *int) (Task, bool, error)

	// Stats
	GetStats() Stats

	// Health check
	IsHealthy() bool
}
