package memory

import (
	"reflect"
	"testing"

	"github.com/a-x-a/go-metric/internal/models/metric"
	"github.com/a-x-a/go-metric/internal/storage"
)

func TestNewMemStorage(t *testing.T) {
	tests := []struct {
		name string
		want *memStorage
	}{
		{
			name: "normal memstorage create",
			want: &memStorage{data: make(map[string]storage.Record)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMemStorage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMemStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_memStorage_Save(t *testing.T) {
	recGauge := storage.NewRecord("gauge1")
	recGauge.SetValue(metric.Gauge(123.456))

	recCounter := storage.NewRecord("counter1")
	recCounter.SetValue(metric.Counter(123))

	type args struct {
		name string
		rec  storage.Record
	}
	tests := []struct {
		name    string
		ms      *memStorage
		args    args
		wantErr bool
	}{
		{
			name: "gauge1 save",
			ms:   &memStorage{data: make(map[string]storage.Record)},
			args: args{
				name: "gauge1",
				rec:  recGauge,
			},
			wantErr: false,
		},
		{
			name: "counter1 save",
			ms:   &memStorage{data: make(map[string]storage.Record)},
			args: args{
				name: "counter1",
				rec:  recCounter,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.ms.Save(tt.args.name, tt.args.rec); (err != nil) != tt.wantErr {
				t.Errorf("memStorage.Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
