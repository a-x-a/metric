package storage

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewDataStorage(t *testing.T) {
	fileName := os.TempDir() + string(os.PathSeparator) + "test_123456789.json"

	t.Run("storage without file", func(t *testing.T) {
		ds := NewDataStorage("", 0)
		require.NotNil(t, ds)
	})

	t.Run("storage with file", func(t *testing.T) {
		ds := NewDataStorage(fileName, 0)
		require.NotNil(t, ds)
	})

}
