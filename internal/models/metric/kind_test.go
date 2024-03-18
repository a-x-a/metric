package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetKind(t *testing.T) {
	type args struct {
		kindRaw string
	}
	tests := []struct {
		name    string
		args    args
		want    MetricKind
		wantErr bool
	}{
		{
			name:    "gauge kind",
			args:    args{"gauge"},
			want:    KindGauge,
			wantErr: false,
		},
		{
			name:    "counter kind",
			args:    args{"counter"},
			want:    KindCounter,
			wantErr: false,
		},
		{
			name:    "zero kind",
			args:    args{""},
			want:    MetricKind(""),
			wantErr: true,
		},
		{
			name:    "zero kind",
			args:    args{"124!@#$"},
			want:    MetricKind(""),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetKind(tt.args.kindRaw)
			if tt.wantErr {
				require.EqualError(t, err, ErrorInvalidMetricKind.Error())
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
