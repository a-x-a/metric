package storage

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/a-x-a/go-metric/internal/models/metric"
)

func Test_Push(t *testing.T) {
	m := NewMemStorage()
	records := [...]Record{
		{name: "Alloc", value: metric.Gauge(12.3456)},
		{name: "PollCount", value: metric.Counter(123)},
		{name: "Random", value: metric.Gauge(1313.1313)},
	}

	type args struct {
		name   string
		record Record
	}
	tests := []struct {
		name    string
		args    args
		want    Record
		wantErr bool
	}{
		{
			name:    "record " + records[0].name,
			args:    args{name: records[0].name, record: records[0]},
			want:    records[0],
			wantErr: false,
		},
		{
			name:    "record " + records[1].name,
			args:    args{name: records[1].name, record: records[1]},
			want:    records[1],
			wantErr: false,
		},
		{
			name:    "record " + records[2].name,
			args:    args{name: records[2].name, record: records[2]},
			want:    records[2],
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := m.Push(tt.args.name, tt.args.record)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, m.data[tt.args.name])
		})
	}
}

func Test_Get(t *testing.T) {
	m := NewMemStorage()
	records := [...]Record{
		{name: "Alloc", value: metric.Gauge(12.3456)},
		{name: "PollCount", value: metric.Counter(123)},
		{name: "Random", value: metric.Gauge(1313.1313)},
	}

	for _, v := range records {
		m.Push(v.name, v)
	}

	type args struct {
		name   string
		record Record
	}
	tests := []struct {
		name string
		args args
		want Record
		ok   bool
	}{
		{
			name: "record " + records[0].name,
			args: args{name: records[0].name, record: records[0]},
			want: records[0],
			ok:   true,
		},
		{
			name: "record " + records[1].name,
			args: args{name: records[1].name, record: records[1]},
			want: records[1],
			ok:   true,
		},
		{
			name: "record " + records[2].name,
			args: args{name: records[2].name, record: records[2]},
			want: records[2],
			ok:   true,
		},
		{
			name: "record unknown",
			args: args{name: ")unknown(", record: Record{name: ")unknown(", value: metric.Metric(nil)}},
			want: Record{},
			ok:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			record, ok := m.Get(tt.args.name)
			if !ok {
				require.Equal(t, tt.ok, ok)
				return
			}
			require.True(t, ok)
			require.Equal(t, tt.want, record)
		})
	}
}

func Test_GetAll(t *testing.T) {
	m := NewMemStorage()
	records := [...]Record{
		{name: "Alloc", value: metric.Gauge(12.345)},
		{name: "PollCount", value: metric.Counter(123)},
		{name: "Random", value: metric.Gauge(1313.131)},
	}
	for _, v := range records {
		m.Push(v.name, v)
	}

	got := m.GetAll()
	require.ElementsMatch(t, records, got)
	require.Equal(t, len(records), len(got))
}

func Test_GetSnapShot(t *testing.T) {
	m := NewMemStorage()
	records := [...]Record{
		{name: "Alloc", value: metric.Gauge(12.345)},
		{name: "PollCount", value: metric.Counter(123)},
		{name: "Random", value: metric.Gauge(1313.131)},
	}
	for _, v := range records {
		m.Push(v.name, v)
	}

	snap := m.GetSnapShot()

	r := m.GetAll()
	rs := snap.GetAll()

	require.ElementsMatch(t, r, rs)
	require.Equal(t, len(r), len(rs))
}

func Test_Close(t *testing.T) {
	require := require.New(t)

	m := NewMemStorage()

	err := m.Close()
	require.NoError(err)
}
