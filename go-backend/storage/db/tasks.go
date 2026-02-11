package db

import (
	"fmt"
	"go-backend/storage"
	"log"
)

func (s *SQLiteStore) GetTasks(status string, userID *int) []storage.Task {
	query := "SELECT id, title, status, user_id FROM tasks WHERE 1=1"
	var args []interface{}

	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}

	if userID != nil {
		query += " AND user_id = ?"
		args = append(args, userID)
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		log.Printf("GetTasks: query tasks: %v", err)
		return []storage.Task{}
	}
	defer rows.Close()

	var tasks []storage.Task
	for rows.Next() {
		var t storage.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Status, &t.UserID); err == nil {
			tasks = append(tasks, t)
		}
	}

	return tasks
}

func (s *SQLiteStore) GetTaskByID(id int) (storage.Task, bool) {
	row := s.db.QueryRow(
		"SELECT id, title, status, user_id FROM tasks WHERE id = ?",
		id,
	)

	var t storage.Task
	err := row.Scan(&t.ID, &t.Title, &t.Status, &t.UserID)
	if err != nil {
		log.Printf("GetTaskByID: scan task: %v", err)
		return storage.Task{}, false
	}

	return t, true
}

func (s *SQLiteStore) CreateTask(title, status string, userID int) (storage.Task, error) {
	result, err := s.db.Exec(
		"INSERT INTO tasks (title, status, user_id) VALUES (?, ?, ?)",
		title, status, userID,
	)
	if err != nil {
		log.Printf("CreateTask: insert task: %v", err)
		return storage.Task{}, fmt.Errorf("insert task: %w", err)
	}

	id, _ := result.LastInsertId()

	return storage.Task{
		ID:     int(id),
		Title:  title,
		Status: status,
		UserID: userID,
	}, nil
}

func (s *SQLiteStore) UpdateTask(
	id int,
	title *string,
	status *string,
	userID *int,
) (storage.Task, bool, error) {

	query := "UPDATE tasks SET "
	var args []interface{}
	updated := false

	if title != nil {
		query += "title = ?, "
		args = append(args, *title)
		updated = true
	}

	if status != nil {
		query += "status = ?, "
		args = append(args, *status)
		updated = true
	}

	if userID != nil {
		query += "user_id = ?, "
		args = append(args, *userID)
		updated = true
	}

	if !updated {
		// No fields provided
		return storage.Task{}, false, fmt.Errorf("no fields provided")
	}

	// Remove trailing comma
	query = query[:len(query)-2] + " WHERE id = ?"
	args = append(args, id)

	result, err := s.db.Exec(query, args...)
	if err != nil {
		log.Printf("UpdateTask: update task: %v", err)
		return storage.Task{}, false, fmt.Errorf("update task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return storage.Task{}, false, fmt.Errorf("rows affected: %w", err)
	}

	if rowsAffected == 0 {
		// Task not found
		return storage.Task{}, false, nil
	}

	// Fetch updated task
	task, found := s.GetTaskByID(id)
	if !found {
		return storage.Task{}, false, nil
	}

	return task, true, nil
}
