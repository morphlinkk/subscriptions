package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/morphlinkk/subscriptions/internal/db/migrations"
)

type Migrator struct {
	db *pgxpool.Pool
}

func NewMigrator(db *pgxpool.Pool) *Migrator {
	return &Migrator{
		db,
	}
}

type migration struct {
	version int
	apply   func(pgx.Tx) error
}

var migrationList = []migration{
	{1, migrations.Init001},
}

func (s *Migrator) Run(ctx context.Context) error {
	if _, err := s.db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY
		);
	`); err != nil {
		return err
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, m := range migrationList {
		var exists bool
		err := tx.QueryRow(ctx,
			`SELECT EXISTS(
        SELECT 1 FROM schema_migrations WHERE version = $1
			)`,
			m.version,
		).Scan(&exists)

		if err != nil {
			return err
		}

		if exists {
			continue
		}

		if err := m.apply(tx); err != nil {
			return fmt.Errorf("migration %03d failed: %w", m.version, err)
		}
		if _, err := tx.Exec(ctx,
			`INSERT INTO schema_migrations (version) VALUES ($1)`,
			m.version,
		); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
