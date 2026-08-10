package main

import (
	"log/slog"
	"os"

	"github.com/aeschyllus/nori/internal/config"
	"github.com/aeschyllus/nori/internal/database"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	command := "up"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	cfg := config.Load()
	db, err := database.Open(cfg.DBPath)
	if err != nil {
		slog.Error("could not open database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := database.RunGoose(db, command); err != nil {
		slog.Error("migration failed", "command", command, "error", err)
		os.Exit(1)
	}
}
