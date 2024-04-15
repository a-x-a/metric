package storage

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type (
	Storage interface {
		Push(ctx context.Context, name string, record Record) error
		PushBatch(ctx context.Context, records []Record) error
		Get(ctx context.Context, name string) (*Record, error)
		GetAll(ctx context.Context) ([]Record, error)
		Close() error
	}

	DBConnPool interface {
		Acquire(ctx context.Context) (*pgxpool.Conn, error)
		Ping(ctx context.Context) error
		Close()
	}
)

func NewDataStorage(dbConn *pgxpool.Pool, path string, storeInterval time.Duration, log *zap.Logger) Storage {
	if dbConn != nil {
		log.Info("attached database storage")
		return NewDBStorage(dbConn, log)
	}

	if len(path) == 0 {
		log.Info("attached in-memory storage")
		return NewMemStorage()
	}

	log.Info("attached in-memory storage with file")
	return NewWithFileStorage(path, storeInterval == 0, log)
}
