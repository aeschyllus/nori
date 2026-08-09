package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const (
	driverName   = "sqlite3"
	dsnSuffix    = "?_busy_timeout=5000&_foreign_keys=on"
	maxOpenConns = 1
	maxIdleConns = 1
)

func Open(dbPath string) (*sql.DB, error) {
	db, err := sql.Open(driverName, dbPath+dsnSuffix)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
