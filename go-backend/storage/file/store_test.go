package file

import (
	"os"
	"testing"

	"go-backend/storage"
)

func newTestStore(t *testing.T) *FileStore {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "store-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmpFile.Close()

	t.Cleanup(func() {
		os.Remove(tmpFile.Name())
	})

	store, err := NewFileStore(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to init store: %v", err)
	}

	// Seed users (needed for tasks)
	store.data.Users = []storage.User{
		{ID: 1, Name: "John", Email: "john@test.com", Role: "dev"},
	}

	return store
}
