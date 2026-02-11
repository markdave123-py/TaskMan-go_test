package file

import "testing"

func TestCreateTask(t *testing.T) {
	store := newTestStore(t)

	task, err := store.CreateTask("Test Task", "pending", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if task.ID == 0 {
		t.Error("expected task ID to be set")
	}

	if task.Status != "pending" {
		t.Errorf("expected status 'pending', got %s", task.Status)
	}
}

func TestUpdateTask_Partial(t *testing.T) {
	store := newTestStore(t)

	task, _ := store.CreateTask("Old Title", "pending", 1)

	newTitle := "New Title"
	updated, found, err := store.UpdateTask(task.ID, &newTitle, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !found {
		t.Fatal("expected task to be found")
	}

	if updated.Title != newTitle {
		t.Errorf("expected title '%s', got '%s'", newTitle, updated.Title)
	}

	if updated.Status != "pending" {
		t.Error("status should not change")
	}
}
