package account

import (
	"database/sql"
	"fmt"
)

const createAccountsTable = `
	CREATE TABLE IF NOT EXISTS accounts (
		account_id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		amount INTEGER NOT NULL
	);
`

func Migrate(db *sql.DB) error {
	if _, err := db.Exec(createAccountsTable); err != nil {
		return fmt.Errorf("create acounts table: %w", err)
	}
	return nil
}

func SeedDemoData(db *sql.DB) error {
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM accounts").Scan(&count); err != nil {
		return fmt.Errorf("failed to count accounts: %w", err)
	}
	if count > 0 {
		return nil
	}

	if _, err := db.Exec(
		"INSERT INTO accounts (name, amount) VALUES (?,?), (?,?)",
		"Savings", 10_000,
		"Emergency Funds", 1_000_000,
	); err != nil {
		return fmt.Errorf("failed to seed data: %w", err)
	}

	return nil
}
