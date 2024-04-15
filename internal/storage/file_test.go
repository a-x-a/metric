package storage

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/a-x-a/go-metric/internal/models/metric"
)

func Test_FileStorage(t *testing.T) {
	var err error
	require := require.New(t)
	log := zap.NewNop()
	fileName := os.TempDir() + string(os.PathSeparator) + "test_123456789.json"
	m := NewWithFileStorage(fileName, false, log)
	records := [...]Record{
		{name: "Alloc", value: metric.Gauge(12.345)},
		{name: "PollCount", value: metric.Counter(123)},
		{name: "Random", value: metric.Gauge(1313.131)},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, v := range records {
		m.Push(ctx, v.name, v)
	}

	err = m.Save()
	require.NoError(err)
	require.FileExists(fileName)

	m2 := NewWithFileStorage(fileName, false, log)
	err = m2.Load()
	require.NoError(err)

	r, err := m.GetAll(ctx)
	require.NoError(err)
	r2, err := m2.GetAll(ctx)
	require.NoError(err)

	require.ElementsMatch(r, r2)
	require.Equal(len(r), len(r2))

	err = os.Remove(fileName)
	require.NoError(err)

	err = m2.Load()
	require.Error(err)

	// dirName := os.TempDir() + string(os.PathSeparator)
	m2 = NewWithFileStorage("", false, log)
	err = m2.Save()
	require.Error(err)

	m2 = NewWithFileStorage(fileName, true, log)
	for _, v := range records {
		m2.Push(ctx, v.name, v)
	}

	err = m2.Close()
	require.NoError(err)
	require.FileExists(fileName)

	r2, err = m.GetAll(ctx)
	require.NoError(err)
	require.ElementsMatch(r, r2)
	require.Equal(len(r), len(r2))

	err = os.Remove(fileName)
	require.NoError(err)
}
