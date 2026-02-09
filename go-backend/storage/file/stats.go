package file

import "go-backend/storage"

func (fs *FileStore) GetStats() storage.Stats {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	var stats storage.Stats

	stats.Users.Total = len(fs.data.Users)
	stats.Tasks.Total = len(fs.data.Tasks)

	for _, task := range fs.data.Tasks {
		switch task.Status {
		case "pending":
			stats.Tasks.Pending++
		case "in-progress":
			stats.Tasks.InProgress++
		case "completed":
			stats.Tasks.Completed++
		}
	}

	return stats
}
