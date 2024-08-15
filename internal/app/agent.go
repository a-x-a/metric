package app

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/a-x-a/go-metric/internal/config"
	"github.com/a-x-a/go-metric/internal/logger"
	"github.com/a-x-a/go-metric/internal/models/metric"
	"github.com/a-x-a/go-metric/internal/security"
	"github.com/a-x-a/go-metric/internal/sender"
)

type (
	Agent struct {
		config    config.AgentConfig
		metrics   metric.Metrics
		logger    *zap.Logger
		transport sender.Sender
		key       security.PublicKey
	}
)

func NewAgent(logLevel string) *Agent {
	var err error
	log := logger.InitLogger(logLevel)
	defer log.Sync()

	cfg := config.NewAgentConfig()
	err = cfg.Parse()
	if err != nil {
		log.Warn("agent failed to parse config", zap.Error(err))
	}

	var publicKey security.PublicKey
	if len(cfg.CryptoKey) != 0 {
		publicKey, err = security.NewPublicKey(cfg.CryptoKey)
		if err != nil {
			log.Panic("agent failed to get public key", zap.Error(err))
		}
	}
	transport := sender.NewHTTPSender(cfg.ServerAddress, cfg.PollInterval, cfg.Key, publicKey)
	return &Agent{
		config:    cfg,
		metrics:   metric.Metrics{},
		logger:    log,
		transport: transport,
		key:       publicKey,
	}
}

func (app *Agent) poll(ctx context.Context, metrics *metric.Metrics) {
	ticker := time.NewTicker(app.config.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			func() {
				app.logger.Info("metrics gathering")

				err := metrics.Poll(ctx)
				if err != nil {
					app.logger.Error("failed to metrics gathering", zap.Error(err))

					return
				}

				app.logger.Info("metrics gathered")
			}()
		case <-ctx.Done():
			app.logger.Info("metrics gathering shutdown")

			return
		}
	}
}

func (app *Agent) report(ctx context.Context, metrics *metric.Metrics) {
	ticker := time.NewTicker(app.config.ReportInterval)
	defer ticker.Stop()

	sndr := sender.New(app.config.Transport, app.config.ServerAddress, app.config.PollInterval, app.config.Key, app.key)
	for {
		select {
		case <-ticker.C:
			func() {
				defer sndr.Reset()

				app.logger.Info("metrics sending")

				err := sendMetrics(ctx, sndr, *metrics, app.config.RateLimit)
				if err != nil {
					app.logger.Error("failed to send metrics", zap.Error(err))
					return
				}

				app.logger.Info("metrics have been sent")
			}()
		case <-ctx.Done():
			app.logger.Info("metrics sending shutdown")

			if err := sndr.Close(); err != nil {
				app.logger.Error("failed metrics sender connection close", zap.Error(err))
			}

			return
		}
	}
}

func (app *Agent) Run(ctx context.Context) {
	go app.poll(ctx, &app.metrics)
	go app.report(ctx, &app.metrics)
}

// sendMetrics отправляет метрики на сервер.
//
// Parameters:
// - ctx: контекст.
// - serverAddress: адрес сервера сбора метрик.
// - timeout: частота отправки метрик на сервер.
// - key: ключ подписи данных.
// - rateLimit: количество одновременно исходящих запросов на сервер.
// - stats: коллекция мсетрик для отправки.
// - сryptoKey публичныq ключ.
func sendMetrics(ctx context.Context, sndr sender.Sender, stats metric.Metrics, rateLimit int) error {
	// отправляем метрики пакета runtime
	sndr.
		Add("Alloc", stats.Runtime.Alloc).
		Add("BuckHashSys", stats.Runtime.BuckHashSys).
		Add("Frees", stats.Runtime.Frees).
		Add("GCCPUFraction", stats.Runtime.GCCPUFraction).
		Add("GCSys", stats.Runtime.GCSys).
		Add("HeapAlloc", stats.Runtime.HeapAlloc).
		Add("HeapIdle", stats.Runtime.HeapIdle).
		Add("HeapInuse", stats.Runtime.HeapInuse).
		Add("HeapObjects", stats.Runtime.HeapObjects).
		Add("HeapReleased", stats.Runtime.HeapReleased).
		Add("HeapSys", stats.Runtime.HeapSys).
		Add("LastGC", stats.Runtime.LastGC).
		Add("Lookups", stats.Runtime.Lookups).
		Add("MCacheInuse", stats.Runtime.MCacheInuse).
		Add("MCacheSys", stats.Runtime.MCacheSys).
		Add("MSpanInuse", stats.Runtime.MSpanInuse).
		Add("MSpanSys", stats.Runtime.MSpanSys).
		Add("Mallocs", stats.Runtime.Mallocs).
		Add("NextGC", stats.Runtime.NextGC).
		Add("NumForcedGC", stats.Runtime.NumForcedGC).
		Add("NumGC", stats.Runtime.NumGC).
		Add("OtherSys", stats.Runtime.OtherSys).
		Add("PauseTotalNs", stats.Runtime.PauseTotalNs).
		Add("StackInuse", stats.Runtime.StackInuse).
		Add("StackSys", stats.Runtime.StackSys).
		Add("Sys", stats.Runtime.Sys).
		Add("TotalAlloc", stats.Runtime.TotalAlloc)
	// отправляем метрики пакета gopsutil
	sndr.
		Add("TotalMemory", stats.PS.TotalMemory).
		Add("FreeMemory", stats.PS.FreeMemory).
		Add("CPUutilization1", stats.PS.CPUutilization1)
	// отправляем обновляемое произвольное значение
	sndr.
		Add("RandomValue", stats.RandomValue)
	// отправляем счётчик обновления метрик пакета runtime
	sndr.
		Add("PollCount", stats.PollCount)

	return sndr.Send(ctx, rateLimit)
}
