package config

import (
	"log/slog"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Env      string
	Port     string
	DBPath   string
	LogLevel slog.Level
	SeedDemo bool

	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ShutdownTimeout   time.Duration
}

func Load() Config {
	return Config{
		Env:               envOr("NORI_ENV", "development"),
		Port:              envOr("NORI_PORT", "8080"),
		DBPath:            envOr("NORI_DB_PATH", "nori.db"),
		LogLevel:          parseLevel(envOr("NORI_LOG_LEVEL", "info")),
		SeedDemo:          envBool("NORI_SEED_DEMO", true),
		ReadHeaderTimeout: durationOr("NORI_READ_HEADER_TIMEOUT", 5*time.Second),
		ReadTimeout:       durationOr("NORI_READ_TIMEOUT", 10*time.Second),
		WriteTimeout:      durationOr("NORI_WRITE_TIMEOUT", 30*time.Second),
		IdleTimeout:       durationOr("NORI_IDLE_TIMEOUT", 120*time.Second),
		ShutdownTimeout:   durationOr("NORI_SHUTDOWN_TIMEOUT", 5*time.Second),
	}
}

func envOr(key string, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func parseLevel(s string) slog.Level {
	switch s {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func envBool(key string, fallback bool) bool {
	if v := os.Getenv(key); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return fallback
}

func durationOr(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}
