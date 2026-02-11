package db

import (
	"go-backend/storage"
	"log"
)

func (s *SQLiteStore) GetStats() storage.Stats {
	var stats storage.Stats

	// Total users
	row := s.db.QueryRow("SELECT COUNT(*) FROM users")
	row.Scan(&stats.Users.Total)

	// Total tasks
	row = s.db.QueryRow("SELECT COUNT(*) FROM tasks")
	row.Scan(&stats.Tasks.Total)

	// Task status breakdown
	rows, err := s.db.Query("SELECT status, COUNT(*) FROM tasks GROUP BY status")
	if err != nil {
		log.Fatalf("GetStats: query task status breakdown: %v", err)
		return stats
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		var count int
		rows.Scan(&status, &count)

		switch status {
		case "pending":
			stats.Tasks.Pending = count
		case "in-progress":
			stats.Tasks.InProgress = count
		case "completed":
			stats.Tasks.Completed = count
		}
	}

	return stats
}
