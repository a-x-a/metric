package storage

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/a-x-a/go-metric/internal/models/metric"
)

func Test_FileStorage(t *testing.T) {
	var err error

	log := zap.NewNop()
	fileName := os.TempDir() + string(os.PathSeparator) + "test_123456789.json"
	m := NewWithFileStorage(fileName, false, log)
	records := [...]Record{
		{name: "Alloc", value: metric.Gauge(12.345)},
		{name: "PollCount", value: metric.Counter(123)},
		{name: "Random", value: metric.Gauge(1313.131)},
	}

	for _, v := range records {
		m.Push(v.name, v)
	}

	err = m.Save()
	require.NoError(t, err)
	require.FileExists(t, fileName)

	m2 := NewWithFileStorage(fileName, false, log)
	err = m2.Load()
	require.NoError(t, err)

	r := m.GetAll()
	r2 := m2.GetAll()

	require.ElementsMatch(t, r, r2)
	require.Equal(t, len(r), len(r2))

	err = os.Remove(fileName)
	require.NoError(t, err)

	err = m2.Load()
	require.Error(t, err)

	// dirName := os.TempDir() + string(os.PathSeparator)
	m2 = NewWithFileStorage("", false, log)
	err = m2.Save()
	require.Error(t, err)

	m2 = NewWithFileStorage(fileName, true, log)
	for _, v := range records {
		m2.Push(v.name, v)
	}

	err = m2.Close()
	require.NoError(t, err)
	require.FileExists(t, fileName)

	r2 = m.GetAll()
	require.ElementsMatch(t, r, r2)
	require.Equal(t, len(r), len(r2))

	err = os.Remove(fileName)
	require.NoError(t, err)
}
