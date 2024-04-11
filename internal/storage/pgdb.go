package storage

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	dbStortage struct {
		dbPool *pgxpool.Pool
	}
)

var _ Storage = &dbStortage{}

func NewDBStorage(dbPool *pgxpool.Pool) *dbStortage {
	return &dbStortage{
		dbPool: dbPool,
	}
}

func (d *dbStortage) Push(key string, record Record) error {
	return errors.New("not implemented")
}

func (d *dbStortage) Get(key string) (Record, bool) {
	return Record{}, false
}

func (d *dbStortage) GetAll() []Record {
	return make([]Record, 0)
}

func (d *dbStortage) Ping(ctx context.Context) error {
	return d.dbPool.Ping(ctx)
}

func (d *dbStortage) Close() error {
	d.dbPool.Close()
	return nil
}
