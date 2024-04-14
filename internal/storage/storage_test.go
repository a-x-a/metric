package storage

import (
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewDataStorage(t *testing.T) {
	require := require.New(t)
	log := zap.NewNop()
	fileName := os.TempDir() + string(os.PathSeparator) + "test_123456789.json"

	t.Run("storage without file", func(t *testing.T) {
		ds := NewDataStorage(nil, "", 0, log)
		require.NotNil(ds)
	})

	t.Run("storage with file", func(t *testing.T) {
		ds := NewDataStorage(nil, fileName, 0, log)
		require.NotNil(ds)
	})

	t.Run("database storage", func(t *testing.T) {
		dbConn := &pgxpool.Pool{}
		ds := NewDataStorage(dbConn, fileName, 0, log)
		require.NotNil(ds)
	})
}
