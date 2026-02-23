package db

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/morphlinkk/subscriptions/internal/config"
)

type Store struct {
	db *pgxpool.Pool
	*Queries
}

var (
	storeInstance *Store
	storeOnce     sync.Once
	storeErr      error
)

func NewStore(ctx context.Context, cfg config.Config) (*Store, error) {
	storeOnce.Do(func() {
		conf, err := pgxpool.ParseConfig(cfg.DatabaseURI)
		if err != nil {
			storeErr = fmt.Errorf("invalid database URI: %w", err)
			return
		}

		conf.MinConns = int32(cfg.DatabaseMinConnections)
		conf.MaxConns = int32(cfg.DatabaseMaxConnections)
		conf.MaxConnLifetime = cfg.DatabaseMaxConnLifetime
		conf.MaxConnIdleTime = 10 * time.Minute
		conf.HealthCheckPeriod = time.Minute

		pool, err := pgxpool.NewWithConfig(ctx, conf)
		if err != nil {
			storeErr = fmt.Errorf("failed to create connection pool: %w", err)
			return
		}

		storeInstance = &Store{db: pool, Queries: New(pool)}
	})

	return storeInstance, storeErr
}

func (s *Store) Ping(ctx context.Context) error {
	return s.db.Ping(ctx)
}

func (s *Store) Close() {
	if s.db != nil {
		s.db.Close()
	}
}

func (s *Store) ExecTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}

	q := s.Queries.WithTx(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("error during transaction: %v, rollback error: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit(ctx)
}
