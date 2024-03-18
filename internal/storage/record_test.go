package storage

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/a-x-a/go-metric/internal/models/metric"
)

func TestNewRecord(t *testing.T) {
	tests := []struct {
		name       string
		recordName string
		want       Record
		wantErr    bool
	}{
		{
			name:       "gauge1 record",
			recordName: "gauge1",
			want:       Record{name: "gauge1"},
			wantErr:    false,
		},
		{
			name:       "zero name record",
			recordName: "",
			want:       Record{},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRecord(tt.recordName)
			if tt.wantErr {
				require.Error(t, err)
				require.Equal(t, tt.want, got)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestRecordMetods(t *testing.T) {
	record := Record{
		name:  "counter",
		value: metric.Counter(123),
	}

	tests := []struct {
		name   string
		record *Record
		value  metric.Metric
		want   *Record
	}{
		{
			name:   "set record value",
			record: &Record{name: "counter"},
			value:  metric.Counter(123),
			want:   &record,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.record.SetValue(tt.value)

			v := tt.record.GetValue()
			require.NotEmpty(t, v)
			require.Equal(t, tt.want.value, v)

			n := tt.record.GetName()
			require.NotEmpty(t, n)
			require.Equal(t, tt.want.name, n)
		})
	}
}
