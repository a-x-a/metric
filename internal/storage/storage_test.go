package storage

import (
	"reflect"
	"testing"

	"github.com/a-x-a/go-metric/internal/models/metric"
)

func TestNewRecord(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want Record
	}{
		{
			name: "gauge1 record",
			args: args{"gauge1"},
			want: Record{name: "gauge1"},
		},
		{
			name: "counter1 record",
			args: args{"counter1"},
			want: Record{name: "counter1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRecord(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRecord() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecord_SetValue(t *testing.T) {
	recGauge := NewRecord("gauge1")
	recGauge.SetValue(metric.Gauge(123.456))

	recCounter := NewRecord("counter1")
	recCounter.SetValue(metric.Counter(123))

	type args struct {
		value metric.Metric
	}
	tests := []struct {
		name string
		rec  *Record
		args args
	}{
		{
			name: "guage1 value",
			rec:  &Record{name: "guage1"},
			args: args{value: recGauge.value},
		},
		{
			name: "counter1 value",
			rec:  &Record{name: "counter1"},
			args: args{value: recCounter.value},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.rec.SetValue(tt.args.value)
		})
	}
}
