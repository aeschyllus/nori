package category

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

func (r *CategoryRepository) List(ctx context.Context) ([]Category, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT category_id, name
		FROM categories
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query categories: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("failed to close rows", "error", err)
		}
	}()

	categories := make([]Category, 0)
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate categories: %w", err)
	}

	return categories, nil
}

func (r *CategoryRepository) GetByID(ctx context.Context, id int64) (Category, error) {
	var c Category

	err := r.db.QueryRowContext(ctx, `
		SELECT category_id, name
		FROM categories
		WHERE category_id = ?
	`, id).Scan(&c.ID, &c.Name)

	if errors.Is(err, sql.ErrNoRows) {
		return Category{}, fmt.Errorf("get category %d: %w", id, ErrCategoryNotFound)
	}
	if err != nil {
		return Category{}, fmt.Errorf("failed to query category: %w", err)
	}

	return c, nil
}

func (r *CategoryRepository) Insert(ctx context.Context, category Category) (Category, error) {
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO categories (name)
		VALUES (?)
		RETURNING category_id
	`, category.Name).Scan(&category.ID)
	if err != nil {
		return Category{}, fmt.Errorf("failed to insert category: %w", err)
	}

	return category, nil
}

func (r *CategoryRepository) Update(ctx context.Context, category Category) (Category, error) {
	err := r.db.QueryRowContext(ctx, `
		UPDATE categories
		SET name = ?
		WHERE category_id = ?
		RETURNING category_id
	`, category.Name, category.ID).Scan(&category.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return Category{}, fmt.Errorf("update category %d: %w", category.ID, ErrCategoryNotFound)
	}
	if err != nil {
		return Category{}, fmt.Errorf("failed to update category: %w", err)
	}

	return category, nil
}
