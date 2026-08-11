package account

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

func (r *AccountRepository) List(ctx context.Context) ([]Account, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT account_id, name, amount
		FROM accounts
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query accounts: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("failed to close rows", "error", err)
		}
	}()

	accounts := make([]Account, 0)
	for rows.Next() {
		var a Account
		if err := rows.Scan(&a.ID, &a.Name, &a.Amount); err != nil {
			return nil, fmt.Errorf("failed to scan account: %w", err)
		}
		accounts = append(accounts, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate accounts: %w", err)
	}

	return accounts, nil
}

func (r *AccountRepository) GetByID(ctx context.Context, id int64) (Account, error) {
	var a Account

	err := r.db.QueryRowContext(ctx, `
		SELECT account_id, name, amount
		FROM accounts
		WHERE account_id = ?
	`, id).Scan(&a.ID, &a.Name, &a.Amount)

	if errors.Is(err, sql.ErrNoRows) {
		return Account{}, fmt.Errorf("get account %d: %w", id, ErrAccountNotFound)
	}
	if err != nil {
		return Account{}, fmt.Errorf("failed to query account: %w", err)
	}

	return a, nil
}

func (r *AccountRepository) Insert(ctx context.Context, account Account) (Account, error) {
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO accounts (name, amount)
		VALUES (?, ?)
		RETURNING account_id
	`, account.Name, account.Amount).Scan(&account.ID)
	if err != nil {
		return Account{}, fmt.Errorf("failed to insert account: %w", err)
	}

	return account, nil
}

func (r *AccountRepository) Update(ctx context.Context, account Account) (Account, error) {
	err := r.db.QueryRowContext(ctx, `
		UPDATE accounts
		SET name = ?, amount = ?
		WHERE account_id = ?
		RETURNING account_id
	`, account.Name, account.Amount, account.ID).Scan(&account.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return Account{}, fmt.Errorf("update account %d: %w", account.ID, ErrAccountNotFound)
	}
	if err != nil {
		return Account{}, fmt.Errorf("failed to update account: %w", err)
	}

	return account, nil
}
