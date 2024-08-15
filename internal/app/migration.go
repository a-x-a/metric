package app

import (
	"database/sql"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"go.uber.org/zap"

	dbmigrations "github.com/a-x-a/go-metric/db"
)

func migrationRun(dsn string, log *zap.Logger) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	d, err := iofs.New(dbmigrations.FS, "migrations")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("iofs", d, "go_metric", driver)
	if err != nil {
		return err
	}

	err = m.Up()
	defer m.Close()

	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info("migrations: no change")
			return nil
		}

		return err
	}

	log.Info("migrations: success")
	return nil
}
