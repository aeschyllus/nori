package account

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/aeschyllus/nori/internal/database"
)

func newTestRepo(t *testing.T) *AccountRepository {
	t.Helper()
	db, err := database.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("database.Open: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return NewAccountRepository(db)
}

func TestListReturnsEmptySlice(t *testing.T) {
	repo := newTestRepo(t)

	accounts, err := repo.List(context.Background())
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if accounts == nil {
		t.Fatalf("List returned nil slice, want non-nil empty slice")
	}
	if len(accounts) != 0 {
		t.Fatalf("List returned %d accounts, want 0", len(accounts))
	}
}

func TestInsertThenGetByID(t *testing.T) {
	repo := newTestRepo(t)

	created, err := repo.Insert(context.Background(), Account{Name: "Travel", Amount: 25_00})
	if err != nil {
		t.Fatalf("Insert: %v", err)
	}
	if created.ID == 0 {
		t.Fatalf("Insert did not assign an ID")
	}

	got, err := repo.GetByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if got.Name != "Travel" || got.Amount != 25_00 {
		t.Fatalf("GetByID = %+v, want name=Travel amount=2500", got)
	}
}

func TestGetByIDNotFound(t *testing.T) {
	repo := newTestRepo(t)

	_, err := repo.GetByID(context.Background(), 999)
	if !errors.Is(err, ErrAccountNotFound) {
		t.Fatalf("GetByID error = %v, want ErrAccountNotFound", err)
	}
}

func TestUpdatePersistsChanges(t *testing.T) {
	repo := newTestRepo(t)

	created, err := repo.Insert(context.Background(), Account{Name: "Old", Amount: 1})
	if err != nil {
		t.Fatalf("Insert: %v", err)
	}

	created.Name = "New"
	created.Amount = 2
	updated, err := repo.Update(context.Background(), created)
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if updated.Name != "New" || updated.Amount != 2 {
		t.Fatalf("Update = %+v, want name=New Amount=2", updated)
	}

	got, err := repo.GetByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("GetByID after update: %v", err)
	}
	if got.Name != "New" || got.Amount != 2 {
		t.Fatalf("stored account = %+v, want name=New Amount=2", got)
	}
}

func TestUpdateNotFound(t *testing.T) {
	repo := newTestRepo(t)

	_, err := repo.Update(context.Background(), Account{ID: 999, Name: "X", Amount: 1})
	if !errors.Is(err, ErrAccountNotFound) {
		t.Fatalf("Update error = %v, want ErrAccountNotFound", err)
	}
}
