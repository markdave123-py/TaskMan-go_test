package db

import (
	"database/sql"
	"fmt"

	"go-backend/storage"

	"github.com/sirupsen/logrus"
)

func (s *SQLiteStore) GetUsers() []storage.User {
	rows, err := s.db.Query("SELECT id, name, email, role FROM users")
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"component": "storage",
			"driver":    "sqlite",
			"operation": "GetUsers",
		}).Error("query users failed")
		return []storage.User{}
	}
	defer rows.Close()

	var users []storage.User
	for rows.Next() {
		var u storage.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role); err == nil {
			users = append(users, u)
		}
	}

	return users
}

func (s *SQLiteStore) GetUserByID(id int) (*storage.User, bool) {
	row := s.db.QueryRow(
		"SELECT id, name, email, role FROM users WHERE id = ?",
		id,
	)

	var u storage.User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Role)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"component": "storage",
			"driver":    "sqlite",
			"operation": "GetUserByID",
			"user_id":   id,
		}).Error("query user failed")
		return nil, false
	}

	return &u, true
}

func (s *SQLiteStore) GetUserByEmail(email string) (*storage.User, bool) {
	row := s.db.QueryRow(
		"SELECT id, name, email, role FROM users WHERE email = ?",
		email,
	)

	var user storage.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Role)
	if err != nil {
		if err != sql.ErrNoRows {
			logrus.WithError(err).WithFields(logrus.Fields{
				"component": "storage",
				"driver":    "sqlite",
				"operation": "GetUserByEmail",
				"email":     email,
			}).Error("query user by email failed")
		}
		return nil, false
	}

	return &user, true
}

func (s *SQLiteStore) CreateUser(name, email, role string) (storage.User, error) {
	result, err := s.db.Exec(
		"INSERT INTO users (name, email, role) VALUES (?, ?, ?)",
		name, email, role,
	)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"component": "storage",
			"driver":    "sqlite",
			"operation": "CreateUser",
			"email":     email,
		}).Error("insert user failed")
		return storage.User{}, fmt.Errorf("insert user: %w", err)
	}

	id, _ := result.LastInsertId()

	return storage.User{
		ID:    int(id),
		Name:  name,
		Email: email,
		Role:  role,
	}, nil
}
