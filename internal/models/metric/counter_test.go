package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCounter_Kind(t *testing.T) {
	tests := []struct {
		name string
		c    Counter
		want string
	}{
		{
			name: "counter kind",
			c:    Counter(123),
			want: "counter",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.c.Kind())
		})
	}
}

func TestCounter_String(t *testing.T) {
	tests := []struct {
		name string
		c    Counter
		want string
	}{
		{
			name: "counter to string",
			c:    Counter(123),
			want: "123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.c.String())
		})
	}
}
