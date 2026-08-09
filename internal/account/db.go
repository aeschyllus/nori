package account

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath+"?_busy_timeout=5000&_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if _, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS accounts (
			account_id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			amount INTEGER NOT NULL
		);
	`); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return db, nil
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
