package db

import (
	"os"
	"testing"
)

func newTestStore(t *testing.T) *SQLiteStore {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "store-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmpFile.Close()

	t.Cleanup(func() {
		os.Remove(tmpFile.Name())
	})

	store, err := NewSQLiteStore(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to init store: %v", err)
	}

	// Seed users (needed for tasks)
	store.CreateUser("John", "john@test.com", "dev")

	return store
}
