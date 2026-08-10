# nori

Small Gin web service for managing money accounts.

## Run

```
go run ./cmd/api
```

Set `NORI_PORT`, `NORI_DB_PATH`, `NORI_LOG_LEVEL`, `NORI_SEED_DEMO` to override defaults (see `internal/config`).

## Environment variables

| Variable                   | Default     |
| -------------------------- | ----------- |
| `NORI_ENV`                 | development |
| `NORI_PORT`                | 8080        |
| `NORI_DB_PATH`             | nori.db     |
| `NORI_LOG_LEVEL`           | info        |
| `NORI_SEED_DEMO`           | true        |
| `NORI_READ_HEADER_TIMEOUT` | 5s          |
| `NORI_READ_TIMEOUT`        | 10s         |
| `NORI_WRITE_TIMEOUT`       | 30s         |
| `NORI_IDLE_TIMEOUT`        | 120s        |
| `NORI_SHUTDOWN_TIMEOUT`    | 5s          |

## Migrations

Schema changes live in `internal/database/migrations` as SQL files with `-- +goose Up` / `-- +goose Down` blocks (see [pressly/goose](https://github.com/pressly/goose)).

Migrations are embedded in the binary and applied automatically at startup when the API starts.

To manage them manually, use the `migrate` command:

| CLI command                    | What it does                                 |
| ------------------------------ | -------------------------------------------- |
| `go run ./cmd/migrate up`      | Apply pending migrations                     |
| `go run ./cmd/migrate down`    | Rollback the most recently applied migration |
| `go run ./cmd/migrate reset`   | Rollback all, then re-apply                  |
| `go run ./cmd/migrate status`  | List applied/pending migrations              |
| `go run ./cmd/migrate version` | Show current version                         |
