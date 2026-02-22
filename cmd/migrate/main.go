package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/morphlinkk/subscriptions/internal/config"
	"github.com/morphlinkk/subscriptions/internal/db"
	"github.com/morphlinkk/subscriptions/internal/logger"
)

func main() {
	ctx := context.Background()
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load configuration", "error", err)
	}

	dbconf, err := pgxpool.ParseConfig(cfg.DatabaseURI)
	if err != nil {
		logger.Fatal("Failed to parse config", "error", err)
	}

	var pool *pgxpool.Pool

	const maxRetries = 30
	for i := 1; i <= maxRetries; i++ {
		pool, err = pgxpool.NewWithConfig(ctx, dbconf)
		if err == nil {
			err = pool.Ping(ctx)
		}

		if err == nil {
			break
		}

		slog.Warn("Waiting for database to be ready...", "attempt", i, "error", err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		logger.Fatal("Could not connect to database after retries", "error", err)
	}

	migrator := db.NewMigrator(pool)
	if err := migrator.Run(ctx); err != nil {
		logger.Fatal("Failed to run migrations", "error", err)
	}

	slog.Info("Migrations ran successfully")
}
