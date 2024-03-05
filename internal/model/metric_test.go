package model

import (
	"errors"
	"reflect"
	"testing"
)

func TestNewMetric(t *testing.T) {
	type args struct {
		name       string
		metricType int
	}
	tests := []struct {
		name    string
		args    args
		want    *metric
		wantErr bool
		err     error
	}{
		{
			name:    "normal create gauge",
			args:    args{name: "gauge#1", metricType: 1},
			want:    &metric{name: "gauge#1", metricType: 1, value: metricValue{0, 0}},
			wantErr: false,
			err:     nil,
		},
		{
			name:    "normal create counter",
			args:    args{name: "counter#1", metricType: 2},
			want:    &metric{name: "counter#1", metricType: 2, value: metricValue{0, 0}},
			wantErr: false,
			err:     nil,
		},
		{
			name:    "create counter without name",
			args:    args{name: "", metricType: 2},
			want:    nil,
			wantErr: true,
			err:     ErroMetricNameIsNull,
		},
		{
			name:    "create counter with incorrect type",
			args:    args{name: "counter2", metricType: -2},
			want:    nil,
			wantErr: true,
			err:     ErroInvalidMetricType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMetric(tt.args.name, tt.args.metricType)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NewMetric() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !errors.Is(err, tt.err) {
					t.Errorf("NewMetric() error = %v, want %v", err, tt.err)
					return
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMetric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetricType_String(t *testing.T) {
	tests := []struct {
		name string
		mt   MetricType
		want string
	}{
		{
			name: "metric type -1",
			mt:   -1,
			want: "",
		},
		{
			name: "metric type 0",
			mt:   0,
			want: "",
		},
		{
			name: "metric type 1 (gauge)",
			mt:   1,
			want: "gauge",
		},
		{
			name: "metric type 2 (counter)",
			mt:   2,
			want: "counter",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mt.String(); got != tt.want {
				t.Errorf("MetricType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
