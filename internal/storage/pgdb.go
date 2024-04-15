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
INSERT INTO metric(id, name, kind, value)
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

func (d *dbStortage) PushBatch(ctx context.Context, records []Record) error {

	conn, err := d.dbPool.Acquire(ctx)
	if err != nil {
		return err
	}

	defer conn.Release()

	queryText := `
INSERT INTO metric(id, name, kind, value)
values ($1, $1, $2, $3)
ON CONFLICT (id) DO UPDATE
SET value = $3;
`

	batch := &pgx.Batch{}
	for _, v := range records {
		batch.Queue(queryText,
			v.GetName(),
			v.GetValue().Kind(),
			v.GetValue(),
		)
	}

	err = conn.SendBatch(ctx, batch).Close()
	if err != nil {
		return err
	}

	return nil
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
SELECT name, kind, value FROM metric WHERE id=$1
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
SELECT name, kind, value FROM metric WHERE
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
