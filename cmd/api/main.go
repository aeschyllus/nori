package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/aeschyllus/nori/internal/account"
	"github.com/aeschyllus/nori/internal/config"
	"github.com/aeschyllus/nori/internal/database"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	// Logger
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: cfg.LogLevel,
		}),
	)
	slog.SetDefault(logger)

	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Database
	db, err := database.Open(cfg.DBPath)
	if err != nil {
		slog.Error("could not initialize database", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("failed to close database", "error", err)
		}
	}()

	if err := database.Migrate(db); err != nil {
		slog.Error("could not migrate database", "error", err)
		os.Exit(1)
	}

	if cfg.SeedDemo {
		if err := account.SeedDemoData(db); err != nil {
			slog.Error("could not seed demo data", "error", err)
			os.Exit(1)
		}
	}

	// Router
	router := gin.Default()
	account.RegisterRoutes(router, db)

	// HTTP server
	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}

	// Start server
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go func() {
		slog.Info("starting server", "address", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server failed", "error", err)
			stop()
		}
	}()
	<-ctx.Done()

	// Allow existing requests to finish
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("forced shutdown", "error", err)
	}
	slog.Info("server stopped")
}
