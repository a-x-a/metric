package storage

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewDataStore(t *testing.T) {
	assert := assert.New(t)
	tt := []struct {
		name     string
		db       *pgxpool.Pool
		path     string
		interval time.Duration
		expected Storage
	}{
		{
			name:     "create database storage",
			path:     "some/path",
			db:       &pgxpool.Pool{},
			interval: 10 * time.Second,
			expected: &dbStortage{},
		},
		{
			name:     "create storage with file (path set, with interval)",
			path:     "some/path",
			interval: 10 * time.Second,
			expected: &withFileStorage{},
		},
		{
			name:     "create storage with file (path set, without interval)",
			path:     "some/path",
			interval: 0,
			expected: &withFileStorage{},
		},
		{
			name:     "create memory storage (with interval)",
			path:     "",
			interval: 10,
			expected: &memStorage{},
		},
		{
			name:     "create memory storage  (without interval)",
			path:     "",
			interval: 0,
			expected: &memStorage{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			store := NewDataStorage(tc.db, tc.path, tc.interval, zap.L())
			assert.IsType(tc.expected, store)
		})
	}
}
