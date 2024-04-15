package app

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func initDB(ctx context.Context, dbPool *pgxpool.Pool) error {
	conn, err := dbPool.Acquire(ctx)
	if err != nil {
		return err
	}

	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	queryText := `
--create types
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'metrickind') THEN
        CREATE TYPE metrickind AS ENUM ('counter', 'gauge');
    END IF;
END$$;

--create tables
CREATE TABLE IF NOT EXISTS metric(
    id    varchar(255) primary key,
    name  varchar(255) not null,
    kind  metrickind not null,
    value double precision
);
`
	if _, err = tx.Exec(ctx, queryText); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
