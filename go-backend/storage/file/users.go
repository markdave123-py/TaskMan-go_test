package file

import "go-backend/storage"

func (fs *FileStore) GetUsers() []storage.User {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	return append([]storage.User(nil), fs.data.Users...)
}

func (fs *FileStore) GetUserByID(id int) (*storage.User, bool) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	for _, u := range fs.data.Users {
		if u.ID == id {
			return &u, true
		}
	}
	return nil, false
}

func (fs *FileStore) CreateUser(name, email, role string) (storage.User, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	maxID := 0
	for _, u := range fs.data.Users {
		if u.ID > maxID {
			maxID = u.ID
		}
	}

	user := storage.User{
		ID:    maxID + 1,
		Name:  name,
		Email: email,
		Role:  role,
	}

	fs.data.Users = append(fs.data.Users, user)
	return user, fs.save()
}
