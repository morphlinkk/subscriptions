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

		storeInstance = &Store{db: pool}
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
