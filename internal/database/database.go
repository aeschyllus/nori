package database

import (
	"database/sql"
	"embed"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

const (
	driverName   = "sqlite3"
	dsnSuffix    = "?_busy_timeout=5000&_foreign_keys=on"
	maxOpenConns = 1
	maxIdleConns = 1
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

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

func Migrate(db *sql.DB) error {
	return RunGoose(db, "up")
}

func RunGoose(db *sql.DB, command string) error {
	goose.SetBaseFS(migrationsFS)
	goose.SetDialect("sqlite3")

	switch command {
	case "up":
		return goose.Up(db, "migrations")
	case "down":
		return goose.Down(db, "migrations")
	case "reset":
		return goose.Reset(db, "migrations")
	case "status":
		return goose.Status(db, "migrations")
	case "version":
		return goose.Version(db, "migrations")
	default:
		return fmt.Errorf("unknown migration command %q", command)
	}
}
