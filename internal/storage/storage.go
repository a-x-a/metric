package storage

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type (
	Storage interface {
		Push(name string, record Record) error
		Get(name string) (Record, bool)
		GetAll() []Record
		Close() error
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
