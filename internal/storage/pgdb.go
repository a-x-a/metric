package storage

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type (
	dbStortage struct {
		dbConn *pgxpool.Pool
		logger *zap.Logger
	}
)

var _ Storage = &dbStortage{}

func NewDBStorage(dbConn *pgxpool.Pool, log *zap.Logger) *dbStortage {
	return &dbStortage{
		dbConn: dbConn,
		logger: log,
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
	return d.dbConn.Ping(ctx)
}

func (d *dbStortage) Close() error {
	d.dbConn.Close()
	return nil
}
