package file

import (
	"errors"
	"go-backend/storage"
)

func (fs *FileStore) CreateTask(title, status string, userID int) (storage.Task, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	userExists := false
	for _, u := range fs.data.Users {
		if u.ID == userID {
			userExists = true
			break
		}
	}
	if !userExists {
		return storage.Task{}, errors.New("user not found")
	}

	maxID := 0
	for _, t := range fs.data.Tasks {
		if t.ID > maxID {
			maxID = t.ID
		}
	}

	task := storage.Task{
		ID:     maxID + 1,
		Title:  title,
		Status: status,
		UserID: userID,
	}

	fs.data.Tasks = append(fs.data.Tasks, task)
	return task, fs.save()
}

func (fs *FileStore) GetTasks(status string, userID *int) []storage.Task {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	var filtered []storage.Task

	for _, task := range fs.data.Tasks {
		if status != "" && task.Status != status {
			continue
		}

		if userID != nil && task.UserID != *userID {
			continue
		}

		filtered = append(filtered, task)
	}

	return filtered
}

func (fs *FileStore) UpdateTask(
	id int,
	title *string,
	status *string,
	userID *int,
) (storage.Task, bool, error) {

	fs.mu.Lock()
	defer fs.mu.Unlock()

	for i, task := range fs.data.Tasks {
		if task.ID != id {
			continue
		}

		// Apply partial updates
		if title != nil {
			task.Title = *title
		}
		if status != nil {
			task.Status = *status
		}
		if userID != nil {
			task.UserID = *userID
		}

		// Persist update
		fs.data.Tasks[i] = task

		if err := fs.save(); err != nil {
			return storage.Task{}, false, err
		}

		return task, true, nil
	}

	return storage.Task{}, false, nil
}
