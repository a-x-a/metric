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
		config  config.AgentConfig
		metrics metric.Metrics
		logger  *zap.Logger
		key     security.PublicKey
	}
)

func NewAgent(logLevel string) *Agent {
	var err error
	log := logger.InitLogger(logLevel)
	defer log.Sync()

	cfg := config.NewAgentConfig()
	err = cfg.Parse()
	if err != nil {
		log.Panic("agent failed to parse config", zap.Error(err))
	}

	var publicKey security.PublicKey
	if len(cfg.CryptoKey) != 0 {
		publicKey, err = security.NewPublicKey(cfg.CryptoKey)
		if err != nil {
			log.Panic("agent failed to get public key", zap.Error(err))
		}
	}

	return &Agent{
		config:  cfg,
		metrics: metric.Metrics{},
		logger:  log,
		key:     publicKey,
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

	for {
		select {
		case <-ticker.C:
			func() {
				app.logger.Info("metrics sending")

				err := sender.SendMetrics(ctx, app.config.ServerAddress, app.config.PollInterval, app.config.Key, app.config.RateLimit, *metrics, app.key)
				if err != nil {
					app.logger.Error("failed to send metrics", zap.Error(err))
					return
				}

				app.logger.Info("metrics have been sent")
			}()
		case <-ctx.Done():
			app.logger.Info("metrics sending shutdown")
			return
		}
	}
}

func (app *Agent) Run(ctx context.Context) {
	go app.poll(ctx, &app.metrics)
	go app.report(ctx, &app.metrics)
}
