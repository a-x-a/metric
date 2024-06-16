package metricservice

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/a-x-a/go-metric/internal/models/metric"
	"github.com/a-x-a/go-metric/internal/storage"
)

func Test_GetWithoutErr(t *testing.T) {
	require := require.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	record, err := storage.NewRecord("PollCount")
	require.NoError(err)
	record.SetValue(metric.Counter(123))

	m := storage.NewMockDataBase(ctrl)
	m.EXPECT().Get(context.Background(), "PollCount").Return(&record, nil)

	s := New(m, zap.L())
	got, err := s.Get(context.Background(), "PollCount", string(metric.KindCounter))
	require.NoError(err)
	require.Equal(record, *got)
}

func Test_GetWithInvalidKind(t *testing.T) {
	require := require.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := storage.NewMockDataBase(ctrl)
	s := New(m, zap.L())
	got, err := s.Get(context.Background(), "PollCount", "invalid kind")
	require.Error(err)
	require.Nil(got)
}

func Test_GetAllWithoutErr(t *testing.T) {
	require := require.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	records := make([]storage.Record, 0)
	record, err := storage.NewRecord("Alloc")
	require.NoError(err)
	record.SetValue(metric.Gauge(12.3456))
	records = append(records, record)

	record, err = storage.NewRecord("PollCount")
	require.NoError(err)
	record.SetValue(metric.Counter(123))
	records = append(records, record)

	record, err = storage.NewRecord("Random")
	require.NoError(err)
	record.SetValue(metric.Gauge(1313.1313))
	records = append(records, record)

	m := storage.NewMockDataBase(ctrl)
	m.EXPECT().GetAll(context.Background()).Return(records, nil)

	s := New(m, zap.L())
	got := s.GetAll(context.Background())
	require.Equal(len(records), len(got))
	require.ElementsMatch(records, got)
}

func Test_GetAllWithErr(t *testing.T) {
	require := require.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := storage.NewMockDataBase(ctrl)
	m.EXPECT().GetAll(context.Background()).Return(nil, errors.New("get all with error"))

	s := New(m, zap.L())
	got := s.GetAll(context.Background())
	require.Nil(got)
}

func Test_PushWithoutErr(t *testing.T) {
	require := require.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	recordPollCount, err := storage.NewRecord("PollCount")
	require.NoError(err)
	recordPollCount.SetValue(metric.Counter(123))
	recordRandom, err := storage.NewRecord("Random")
	require.NoError(err)
	recordRandom.SetValue(metric.Gauge(1313.1313))

	m := storage.NewMockDataBase(ctrl)
	m.EXPECT().Push(context.Background(), "PollCount", recordPollCount).Return(nil)
	recordPollCount.SetValue(metric.Counter(0))
	m.EXPECT().Get(context.Background(), "PollCount").Return(&recordPollCount, nil)
	m.EXPECT().Push(context.Background(), "Random", recordRandom).Return(nil)

	s := New(m, zap.L())
	err = s.Push(context.Background(), "PollCount", string(metric.KindCounter), "123")
	require.NoError(err)

	err = s.Push(context.Background(), "Random", string(metric.KindGauge), "1313.1313")
	require.NoError(err)
}

func Test_PushWithInvalidKind(t *testing.T) {
	require := require.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := storage.NewMockDataBase(ctrl)
	s := New(m, zap.L())
	err := s.Push(context.Background(), "PollCount", "invalidKind", "123")
	require.Error(err)
}

func Test_PushWithInvalidName(t *testing.T) {
	require := require.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := storage.NewMockDataBase(ctrl)
	s := New(m, zap.L())
	err := s.Push(context.Background(), "", string(metric.KindCounter), "123")
	require.Error(err)
}

func Test_PushWithInvalidValue(t *testing.T) {
	require := require.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := storage.NewMockDataBase(ctrl)
	s := New(m, zap.L())
	err := s.Push(context.Background(), "PollCount", string(metric.KindCounter), "$#@!")
	require.Error(err)
	err = s.Push(context.Background(), "Random", string(metric.KindGauge), "$#@!")
	require.Error(err)
}

func Test_PushGaugeWithoutErr(t *testing.T) {
	require := require.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	recordRandom, err := storage.NewRecord("Random")
	require.NoError(err)
	recordRandom.SetValue(metric.Gauge(1313.1313))

	m := storage.NewMockDataBase(ctrl)
	m.EXPECT().Push(context.Background(), "Random", recordRandom).Return(nil)

	s := New(m, zap.L())
	got, err := s.PushGauge(context.Background(), "Random", metric.Gauge(1313.1313))
	require.NoError(err)
	require.Equal(metric.Gauge(1313.1313), got)
}

func Test_PushGaugeWithErr(t *testing.T) {
	require := require.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	recordRandom, err := storage.NewRecord("Random")
	require.NoError(err)
	recordRandom.SetValue(metric.Gauge(1313.1313))

	m := storage.NewMockDataBase(ctrl)
	m.EXPECT().Push(context.Background(), "Random", recordRandom).Return(errors.New("push with error"))

	s := New(m, zap.L())
	got, err := s.PushGauge(context.Background(), "Random", metric.Gauge(1313.1313))
	require.Error(err)
	require.Equal(metric.Gauge(0), got)
}

func Test_PushGaugeWithInvalidName(t *testing.T) {
	require := require.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := storage.NewMockDataBase(ctrl)
	s := New(m, zap.L())
	got, err := s.PushGauge(context.Background(), "", metric.Gauge(1313.1313))
	require.Error(err)
	require.Equal(metric.Gauge(0), got)
}

func Test_PushCounterWithoutErr(t *testing.T) {
	require := require.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	recordPollCount, err := storage.NewRecord("PollCount")
	require.NoError(err)
	recordPollCount.SetValue(metric.Counter(123))

	m := storage.NewMockDataBase(ctrl)
	m.EXPECT().Push(context.Background(), "PollCount", recordPollCount).Return(nil)
	recordPollCount.SetValue(metric.Counter(0))
	m.EXPECT().Get(context.Background(), "PollCount").Return(&recordPollCount, nil)

	s := New(m, zap.L())
	got, err := s.PushCounter(context.Background(), "PollCount", metric.Counter(123))
	require.NoError(err)
	require.Equal(metric.Counter(123), got)
}

func Test_PushCounterWithInvalidName(t *testing.T) {
	require := require.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := storage.NewMockDataBase(ctrl)
	s := New(m, zap.L())
	got, err := s.PushCounter(context.Background(), "", metric.Counter(123))
	require.Error(err)
	require.Equal(metric.Counter(0), got)
}

func Test_PushBatchWithoutErr(t *testing.T) {
	require := require.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	records := make([]storage.Record, 0)
	record, err := storage.NewRecord("Alloc")
	require.NoError(err)
	record.SetValue(metric.Gauge(12.3456))
	records = append(records, record)

	record, err = storage.NewRecord("PollCount")
	require.NoError(err)
	record.SetValue(metric.Counter(123))
	records = append(records, record)

	record, err = storage.NewRecord("Random")
	require.NoError(err)
	record.SetValue(metric.Gauge(1313.1313))
	records = append(records, record)

	record, err = storage.NewRecord("PollCount1")
	require.NoError(err)
	record.SetValue(metric.Counter(123))
	records = append(records, record)

	m := storage.NewMockDataBase(ctrl)
	m.EXPECT().PushBatch(context.Background(), records).Return(nil)
	record, err = storage.NewRecord("PollCount")
	require.NoError(err)
	record.SetValue(metric.Counter(0))
	m.EXPECT().Get(context.Background(), "PollCount").Return(&record, nil)
	m.EXPECT().Get(context.Background(), "PollCount1").Return(nil, errors.New("get counter error"))

	record, err = storage.NewRecord("PollCount")
	require.NoError(err)
	record.SetValue(metric.Counter(0))
	records = append(records, record)

	s := New(m, zap.L())
	err = s.PushBatch(context.Background(), records)
	require.NoError(err)
}

func Test_Ping(t *testing.T) {
	require := require.New(t)
	tt := []struct {
		name   string
		result error
	}{
		{
			name:   "return no error DB is online",
			result: nil,
		},
		{
			name:   "return error DB is ofline",
			result: errors.New("DB is offline"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := storage.NewMockDataBase(ctrl)
			m.EXPECT().Ping(context.Background()).Return(tc.result)

			s := New(m, zap.L())
			err := s.Ping(context.Background())
			require.Equal(err, tc.result)
		})
	}
}

func Test_PingErrNotSupportedMethod(t *testing.T) {
	require := require.New(t)
	tt := []struct {
		name   string
		result error
	}{
		{
			name:   "return ErrNotSupportedMethod",
			result: ErrNotSupportedMethod,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := storage.NewMockStorage(ctrl)
			s := New(m, zap.L())
			err := s.Ping(context.Background())
			require.Equal(err, tc.result)
		})
	}
}
