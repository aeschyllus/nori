package account

import (
	"context"
	"errors"
	"testing"
)

type fakeRepo struct {
	accounts []Account
}

func (f *fakeRepo) List(ctx context.Context) ([]Account, error) {
	return f.accounts, nil
}

func (f *fakeRepo) GetByID(ctx context.Context, id int64) (Account, error) {
	for _, a := range f.accounts {
		if a.ID == id {
			return a, nil
		}
	}
	return Account{}, ErrAccountNotFound
}

func (f *fakeRepo) Insert(ctx context.Context, a Account) (Account, error) {
	a.ID = int64(len(f.accounts)) + 1
	f.accounts = append(f.accounts, a)
	return a, nil
}

func (f *fakeRepo) Update(ctx context.Context, a Account) (Account, error) {
	for i := range f.accounts {
		if f.accounts[i].ID == a.ID {
			f.accounts[i] = a
			return a, nil
		}
	}
	return Account{}, ErrAccountNotFound
}

func TestCreateTrimsNameAndAssignsID(t *testing.T) {
	svc := NewAccountService(&fakeRepo{})

	got, err := svc.Create(context.Background(), Account{Name: "  Savings  ", Amount: 100})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if got.Name != "Savings" {
		t.Fatalf("Create name = %q, want trimmed %q", got.Name, "Savings")
	}
	if got.ID == 0 {
		t.Fatalf("Create did not assign an ID")
	}
}

func TestCreateRejectsBlankName(t *testing.T) {
	svc := NewAccountService(&fakeRepo{})

	_, err := svc.Create(context.Background(), Account{Name: "   "})
	if !errors.Is(err, ErrInvalidAccount) {
		t.Fatalf("Create error = %v, want ErrInvalidAccount", err)
	}
}

func TestCreateRejectsOverLongName(t *testing.T) {
	svc := NewAccountService(&fakeRepo{})
	long := make([]byte, maxNameLength+1)
	for i := range long {
		long[i] = 'a'
	}

	_, err := svc.Create(context.Background(), Account{Name: string(long)})
	if !errors.Is(err, ErrInvalidAccount) {
		t.Fatalf("Create error = %v, want ErrInvalidAccount", err)
	}
}

func TestUpdateRejectsNonPositiveID(t *testing.T) {
	svc := NewAccountService(&fakeRepo{})

	_, err := svc.Update(context.Background(), 0, Account{Name: "X", Amount: 1})
	if !errors.Is(err, ErrInvalidAccount) {
		t.Fatalf("Update error = %v, want ErrInvalidAccount", err)
	}
}

func TestUpdateSetsIDBeforePersisting(t *testing.T) {
	repo := &fakeRepo{accounts: []Account{{ID: 7, Name: "Old", Amount: 1}}}
	svc := NewAccountService(repo)

	got, err := svc.Update(context.Background(), 7, Account{Name: "New", Amount: 2})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if got.ID != 7 {
		t.Fatalf("Update returned ID %d, want 7", got.ID)
	}
	if repo.accounts[0].Name != "New" {
		t.Fatalf("Update did not persist, stored = %+v", repo.accounts[0])
	}
}
