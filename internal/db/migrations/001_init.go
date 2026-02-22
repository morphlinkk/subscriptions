package migrations

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func Init001(tx pgx.Tx) error {
	query := `CREATE TABLE IF NOT EXISTS subscriptions(
    id BIGSERIAL PRIMARY KEY,
    service_name VARCHAR NOT NULL,
    price INTEGER NOT NULL,
    user_id UUID NOT NULL,
    start_date timestamp NOT NULL,
    end_date timestamp
  );`

	if _, err := tx.Exec(context.Background(), query); err != nil {
		return err
	}

	return nil
}
