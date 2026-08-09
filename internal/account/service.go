package account

import (
	"context"
	"fmt"
	"strings"
)

const maxNameLength = 255

// Deleting accounts is intentionally not supported
type Repository interface {
	List(ctx context.Context) ([]Account, error)
	GetByID(ctx context.Context, id int64) (Account, error)
	Insert(ctx context.Context, account Account) (Account, error)
	Update(ctx context.Context, account Account) (Account, error)
}

type AccountService struct {
	repo Repository
}

func NewAccountService(repo Repository) *AccountService {
	return &AccountService{
		repo: repo,
	}
}

func (s *AccountService) List(ctx context.Context) ([]Account, error) {
	return s.repo.List(ctx)
}

func (s *AccountService) GetByID(ctx context.Context, id int64) (Account, error) {
	if id <= 0 {
		return Account{}, fmt.Errorf("%w: id must be positive", ErrInvalidAccount)
	}
	return s.repo.GetByID(ctx, id)
}

func (s *AccountService) Create(ctx context.Context, account Account) (Account, error) {
	if err := validateAccount(&account); err != nil {
		return Account{}, err
	}
	return s.repo.Insert(ctx, account)
}

func (s *AccountService) Update(ctx context.Context, id int64, account Account) (Account, error) {
	if id <= 0 {
		return Account{}, fmt.Errorf("%w: id must be positive", ErrInvalidAccount)
	}
	if err := validateAccount(&account); err != nil {
		return Account{}, err
	}
	account.ID = id
	return s.repo.Update(ctx, account)
}

func validateAccount(a *Account) error {
	name := strings.TrimSpace(a.Name)
	if name == "" {
		return fmt.Errorf("%w: name is required", ErrInvalidAccount)
	}
	if len(name) > maxNameLength {
		return fmt.Errorf("%w: name must be at most %d characters", ErrInvalidAccount, maxNameLength)
	}

	a.Name = name
	return nil
}
