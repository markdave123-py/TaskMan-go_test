package file

import "testing"

func TestCreateUser(t *testing.T) {
	store := newTestStore(t)

	user, err := store.CreateUser("Alice", "alice@test.com", "dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if user.ID == 0 {
		t.Error("expected user ID to be set")
	}

	if user.Email != "alice@test.com" {
		t.Errorf("unexpected email: %s", user.Email)
	}
}

func TestGetUserByID(t *testing.T) {
	store := newTestStore(t)

	created, _ := store.CreateUser("Bob", "bob@test.com", "qa")

	user, ok := store.GetUserByID(created.ID)
	if !ok {
		t.Fatal("expected user to be found")
	}

	if user.Name != "Bob" {
		t.Errorf("expected name Bob, got %s", user.Name)
	}
}
