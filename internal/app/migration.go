package app

import (
	"errors"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

func migrationRun(dsn string, log *zap.Logger) error {
	var (
		retries  uint = 10
		migrator *migrate.Migrate
		err      error
	)

	for retries > 0 {
		migrator, err = migrate.New("file://sql", string(dsn))
		if err == nil {
			break
		}

		retries--

		log.Info("migrations: trying connect to DB", zap.String("url", dsn))
		time.Sleep(time.Second)
	}

	if err != nil {
		return err
	}

	err = migrator.Up()
	defer migrator.Close()

	if err == nil {
		log.Info("migrations: success")
		return nil
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Info("migrations: no change")
		return nil
	}

	return err
}
