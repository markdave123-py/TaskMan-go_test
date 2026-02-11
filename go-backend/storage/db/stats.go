package db

import (
	"go-backend/storage"

	"github.com/sirupsen/logrus"
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
		logrus.WithError(err).WithFields(logrus.Fields{
			"component": "storage",
			"driver":    "sqlite",
			"operation": "GetStats",
		}).Error("query task status breakdown failed")
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
