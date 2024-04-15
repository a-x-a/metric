package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"github.com/a-x-a/go-metric/internal/models/metric"
)

type (
	dbStortage struct {
		dbPool DBConnPool
		logger *zap.Logger
	}
)

var _ Storage = &dbStortage{}

func NewDBStorage(dbConn DBConnPool, log *zap.Logger) *dbStortage {
	d := dbStortage{
		dbPool: dbConn,
		logger: log,
	}

	return &d
}

func (d *dbStortage) Push(ctx context.Context, key string, record Record) error {
	conn, err := d.dbPool.Acquire(ctx)
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
INSERT INTO metrics(id, name, kind, value)
values ($1, $2, $3, $4)
ON CONFLICT (id) DO UPDATE
SET value = $4;
`

	if _, err = tx.Exec(ctx, queryText,
		key,
		record.GetName(),
		record.GetValue().Kind(),
		record.GetValue(),
	); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (d *dbStortage) Get(ctx context.Context, key string) (*Record, error) {
	conn, err := d.dbPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Release()

	var (
		name     string
		kindRaw  string
		valueRaw float64
	)

	queryText := `
SELECT name, kind, value FROM metrics WHERE id=$1
	`
	err = conn.QueryRow(ctx, queryText, key).Scan(&name, &kindRaw, &valueRaw)
	if err != nil {
		return nil, err
	}

	kind, err := metric.GetKind(kindRaw)
	if err != nil {
		return nil, err
	}

	record, err := NewRecord(name)
	if err != nil {
		return nil, err
	}

	switch kind {
	case metric.KindCounter:
		value := metric.Counter(valueRaw)
		record.SetValue(value)

		return &record, nil
	case metric.KindGauge:
		value := metric.Gauge(valueRaw)
		record.SetValue(value)

		return &record, nil
	default:
		return nil, metric.ErrorInvalidMetricKind
	}
}

func (d *dbStortage) GetAll(ctx context.Context) ([]Record, error) {
	conn, err := d.dbPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Release()

	queryText := `
SELECT name, kind, value FROM metrics WHERE
	`
	rows, err := conn.Query(ctx, queryText)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var (
		name     string
		kindRaw  string
		valueRaw float64
	)

	records := make([]Record, 0)
	_, err = pgx.ForEachRow(rows, []any{&name, &kindRaw, &valueRaw}, func() error {
		kind, err := metric.GetKind(kindRaw)
		if err != nil {
			return err
		}

		record, err := NewRecord(name)
		if err != nil {
			return err
		}

		switch kind {
		case metric.KindCounter:
			value := metric.Counter(valueRaw)
			record.SetValue(value)
		case metric.KindGauge:
			value := metric.Gauge(valueRaw)
			record.SetValue(value)
		default:
			return metric.ErrorInvalidMetricKind
		}

		records = append(records, record)

		return nil
	})

	return records, err
}

func (d *dbStortage) Ping(ctx context.Context) error {
	return d.dbPool.Ping(ctx)
}

func (d *dbStortage) Close() error {
	d.dbPool.Close()
	return nil
}

func (d *dbStortage) Bootstrap(ctx context.Context) error {
	conn, err := d.dbPool.Acquire(ctx)
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
