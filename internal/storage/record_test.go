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

func TestRecord_UnmarshalJSON(t *testing.T) {
	t.Run("normal unmarshal JSON", func(t *testing.T) {
		require := require.New(t)

		r := &Record{name: "counter"}
		data := []byte(`{"name":"counter"}`)
		err := r.UnmarshalJSON(data)

		require.NoError(err)
	})

	t.Run("error unmarshal JSON", func(t *testing.T) {
		require := require.New(t)

		r := &Record{name: "counter"}
		data := []byte(`invalid`)
		err := r.UnmarshalJSON(data)

		require.Error(err)

	})

}

func TestRecord_MarshalJSON(t *testing.T) {
	t.Run("normal marshal JSON", func(t *testing.T) {
		require := require.New(t)

		r, err := NewRecord("counter")

		require.NoError(err)

		r.SetValue(metric.Counter(10))

		_, err = r.MarshalJSON()

		require.NoError(err)
	})
}
