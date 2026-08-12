package category

import (
	"context"
	"errors"
	"testing"
)

type fakeRepo struct {
	categories []Category
}

func (f *fakeRepo) List(ctx context.Context) ([]Category, error) {
	return f.categories, nil
}

func (f *fakeRepo) GetByID(ctx context.Context, id int64) (Category, error) {
	for _, c := range f.categories {
		if c.ID == id {
			return c, nil
		}
	}
	return Category{}, ErrCategoryNotFound
}

func (f *fakeRepo) Insert(ctx context.Context, c Category) (Category, error) {
	c.ID = int64(len(f.categories)) + 1
	f.categories = append(f.categories, c)
	return c, nil
}

func (f *fakeRepo) Update(ctx context.Context, c Category) (Category, error) {
	for i := range f.categories {
		if f.categories[i].ID == c.ID {
			f.categories[i] = c
			return c, nil
		}
	}
	return Category{}, ErrCategoryNotFound
}

func TestCreateTrimsNameAndAssignsID(t *testing.T) {
	svc := NewCategoryService(&fakeRepo{})

	created, err := svc.Create(context.Background(), Category{Name: "Food"})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if created.Name != "Food" {
		t.Fatalf("Create name = %q, want trimmed %q", created.Name, "Food")
	}
	if created.ID == 0 {
		t.Fatalf("Create did not assign ID")
	}
}

func TestCreateRejectsBlankName(t *testing.T) {
	svc := NewCategoryService(&fakeRepo{})

	_, err := svc.Create(context.Background(), Category{Name: "   "})
	if !errors.Is(err, ErrInvalidCategory) {
		t.Fatalf("Create error = %v, want ErrInvalidCategory", err)
	}
}

func TestCreateRejectsOverLongName(t *testing.T) {
	svc := NewCategoryService(&fakeRepo{})
	long := make([]byte, maxNameLength+1)
	for i := range long {
		long[i] = 'a'
	}

	_, err := svc.Create(context.Background(), Category{Name: string(long)})
	if !errors.Is(err, ErrInvalidCategory) {
		t.Fatalf("Create error = %v, want ErrInvalidCategory", err)
	}
}

func TestUpdateRejectsNonPositiveID(t *testing.T) {
	svc := NewCategoryService(&fakeRepo{})

	_, err := svc.Update(context.Background(), 0, Category{Name: "X"})
	if !errors.Is(err, ErrInvalidCategory) {
		t.Fatalf("Update error = %v, want ErrInvalidCategory", err)
	}
}

func TestUpdateSetsIDBeforePersisting(t *testing.T) {
	repo := &fakeRepo{categories: []Category{{ID: 7, Name: "Old"}}}
	svc := NewCategoryService(repo)

	got, err := svc.Update(context.Background(), 7, Category{Name: "New"})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if got.ID != 7 {
		t.Fatalf("Update returned ID %d, want 7", got.ID)
	}
	if repo.categories[0].Name != "New" {
		t.Fatalf("Update did not persist, stored = %+v", repo.categories[0])
	}
}
