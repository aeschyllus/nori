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
	db, err := account.InitDB(cfg.DBPath)
	if err != nil {
		slog.Error("could not initialize database", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("failed to close database", "error", err)
		}
	}()

	if cfg.SeedDemo {
		if err := account.SeedDemoData(db); err != nil {
			slog.Error("could not seed demo data", "error", err)
			os.Exit(1)
		}
	}

	// Dependencies
	repo := account.NewAccountRepository(db)
	svc := account.NewAccountService(repo)
	h := account.NewAccountHandler(svc)

	// Router
	router := gin.Default()

	accounts := router.Group("/accounts")
	{
		accounts.GET("/", h.ListAccounts)
		accounts.GET("/:id", h.GetAccount)
		accounts.POST("/", h.CreateAccount)
		accounts.PUT("/:id", h.UpdateAccount)
	}

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
