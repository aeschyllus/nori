package category

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/aeschyllus/nori/internal/database"
)

func newTestRepo(t *testing.T) *CategoryRepository {
	t.Helper()
	db, err := database.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("database.Open: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := database.Migrate(db); err != nil {
		t.Fatalf("Migrate: %v", err)
	}
	return NewCategoryRepository(db)
}

func TestListReturnsEmptySlice(t *testing.T) {
	repo := newTestRepo(t)

	categories, err := repo.List(context.Background())
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if categories == nil {
		t.Fatalf("List returned nil slice, want non-nil empty slice")
	}
	if len(categories) != 0 {
		t.Fatalf("List returned %d categories, want 0", len(categories))
	}
}

func TestInsertThenGetByID(t *testing.T) {
	repo := newTestRepo(t)

	created, err := repo.Insert(context.Background(), Category{Name: "New"})
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
	if got.Name != "New" {
		t.Fatalf("GetByID: %+v, want name=New", got)
	}
}

func TestGetByIDNotFound(t *testing.T) {
	repo := newTestRepo(t)

	_, err := repo.GetByID(context.Background(), 999)
	if !errors.Is(err, ErrCategoryNotFound) {
		t.Fatalf("GetByID error = %v, want ErrCategoryNotFound", err)
	}
}

func TestUpdatePersistsChanges(t *testing.T) {
	repo := newTestRepo(t)

	created, err := repo.Insert(context.Background(), Category{Name: "Old"})
	if err != nil {
		t.Fatalf("Insert: %v", err)
	}

	created.Name = "New"
	updated, err := repo.Update(context.Background(), created)
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if updated.Name != "New" {
		t.Fatalf("Update = %+v, want name=New", updated)
	}

	got, err := repo.GetByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("GetByID after update: %v", err)
	}
	if got.Name != "New" {
		t.Fatalf("stored category = %+v, want name=New", got)
	}
}

func TestUpdateNotFound(t *testing.T) {
	repo := newTestRepo(t)

	_, err := repo.Update(context.Background(), Category{ID: 999, Name: "X"})
	if !errors.Is(err, ErrCategoryNotFound) {
		t.Fatalf("Update error = %v, want ErrCategoryNotFound", err)
	}
}
