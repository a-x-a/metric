package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/a-x-a/go-metric/internal/models/metric"
)

func Test_Push(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
			err := m.Push(ctx, tt.args.name, tt.args.record)
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
	require := require.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	m := NewMemStorage()
	records := [...]Record{
		{name: "Alloc", value: metric.Gauge(12.3456)},
		{name: "PollCount", value: metric.Counter(123)},
		{name: "Random", value: metric.Gauge(1313.1313)},
	}

	for _, v := range records {
		m.Push(ctx, v.name, v)
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
		{
			name:    "record unknown",
			args:    args{name: ")unknown(", record: Record{name: ")unknown(", value: metric.Metric(nil)}},
			want:    Record{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			record, err := m.Get(ctx, tt.args.name)
			if tt.wantErr {
				require.Error(err)
				return
			}

			require.NoError(err)
			require.Equal(tt.want, *record)
		})
	}
}

func Test_GetAll(t *testing.T) {
	require := require.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	m := NewMemStorage()
	records := [...]Record{
		{name: "Alloc", value: metric.Gauge(12.345)},
		{name: "PollCount", value: metric.Counter(123)},
		{name: "Random", value: metric.Gauge(1313.131)},
	}
	for _, v := range records {
		m.Push(ctx, v.name, v)
	}

	got, err := m.GetAll(ctx)
	require.NoError(err)
	require.ElementsMatch(records, got)
	require.Equal(len(records), len(got))
}

func Test_GetSnapShot(t *testing.T) {
	require := require.New(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	m := NewMemStorage()
	records := [...]Record{
		{name: "Alloc", value: metric.Gauge(12.345)},
		{name: "PollCount", value: metric.Counter(123)},
		{name: "Random", value: metric.Gauge(1313.131)},
	}
	for _, v := range records {
		m.Push(ctx, v.name, v)
	}

	snap := m.GetSnapShot()

	r, err := m.GetAll(ctx)
	require.NoError(err)

	rs, err := snap.GetAll(ctx)
	require.NoError(err)

	require.ElementsMatch(r, rs)
	require.Equal(len(r), len(rs))
}

func Test_Close(t *testing.T) {
	require := require.New(t)

	m := NewMemStorage()

	err := m.Close()
	require.NoError(err)
}
