package category

import (
	"context"
	"fmt"
	"strings"
)

const maxNameLength = 255

type Repository interface {
	List(ctx context.Context) ([]Category, error)
	GetByID(ctx context.Context, id int64) (Category, error)
	Insert(ctx context.Context, category Category) (Category, error)
	Update(ctx context.Context, category Category) (Category, error)
}

type CategoryService struct {
	repo Repository
}

func NewCategoryService(repo Repository) *CategoryService {
	return &CategoryService{
		repo: repo,
	}
}

func (s *CategoryService) List(ctx context.Context) ([]Category, error) {
	return s.repo.List(ctx)
}

func (s *CategoryService) GetByID(ctx context.Context, id int64) (Category, error) {
	if id <= 0 {
		return Category{}, fmt.Errorf("%w: id must be positive", ErrInvalidCategory)
	}
	return s.repo.GetByID(ctx, id)
}

func (s *CategoryService) Create(ctx context.Context, category Category) (Category, error) {
	if err := validateCategory(&category); err != nil {
		return Category{}, err
	}
	return s.repo.Insert(ctx, category)
}

func (s *CategoryService) Update(ctx context.Context, id int64, category Category) (Category, error) {
	if id <= 0 {
		return Category{}, fmt.Errorf("%w: id must be positive", ErrInvalidCategory)
	}
	if err := validateCategory(&category); err != nil {
		return Category{}, err
	}
	category.ID = id
	return s.repo.Update(ctx, category)
}

func validateCategory(c *Category) error {
	name := strings.TrimSpace(c.Name)
	if name == "" {
		return fmt.Errorf("%w: name is required", ErrInvalidCategory)
	}
	if len(name) > maxNameLength {
		return fmt.Errorf("%w: name must be at most %d characters", ErrInvalidCategory, maxNameLength)
	}

	c.Name = name
	return nil
}
