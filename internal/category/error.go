package category

import "errors"

var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrInvalidCategory  = errors.New("invalid category")
)
