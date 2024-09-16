package storage

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/a-x-a/go-metric/internal/models/metric"
)

func TestNewDBStorage(t *testing.T) {
	require := require.New(t)

	dbConn := &pgxpool.Pool{}
	ds := NewDBStorage(dbConn, zap.L())
	require.NotNil(ds)
}

func Test_dbStortage_Push(t *testing.T) {
	require := require.New(t)

	d := dbStortage{}
	r := Record{
		name:  "Alloc",
		value: metric.Gauge(12.345),
	}

	err := d.Push(r.name, r)

	require.Error(err)
}

func Test_dbStortage_Get(t *testing.T) {
	require := require.New(t)

	d := dbStortage{}
	r, ok := d.Get("temp")

	require.False(ok)
	require.NotNil(r)
}

func Test_dbStortage_GetAll(t *testing.T) {
	require := require.New(t)

	d := dbStortage{}
	r := d.GetAll()

	require.NotNil(r)
}

func Test_dbStortage_Ping(t *testing.T) {
	require := require.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	dbConn, err := pgxpool.New(ctx, "host=localhost")
	require.NoError(err)

	d := dbStortage{
		dbConn: dbConn,
	}

	err = d.Ping(ctx)
	require.Error(err)
}

func Test_dbStortage_Close(t *testing.T) {
	require := require.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	dbConn, err := pgxpool.New(ctx, "host=localhost")
	require.NoError(err)

	d := dbStortage{
		dbConn: dbConn,
	}

	err = d.Close()
	require.NoError(err)
}
