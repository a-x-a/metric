package storage

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestPing(t *testing.T) {
	tt := []struct {
		name   string
		result error
	}{
		{
			name: "return no error DB is online",
		},
		{
			name:   "return error DB is ofline",
			result: errors.New("DB is offline"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			m := NewDBConnPoolMock()
			m.On("Ping", mock.Anything).Return(tc.result)

			s := NewDBStorage(m, zap.L())
			assert.ErrorIs(t, tc.result, s.Ping(context.Background()))
		})
	}
}

func TestCloseNeverFails(t *testing.T) {
	m := NewDBConnPoolMock()
	m.On("Close").Return()

	s := NewDBStorage(m, zap.L())
	assert.NoError(t, s.Close())
}
