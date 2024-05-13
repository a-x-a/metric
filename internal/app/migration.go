package app

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
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

	startDir := ""
	if pwd, err := os.Getwd(); err == nil {
		startDir = strings.ReplaceAll(pwd, "\\", "/")
	}

	migrationDir := "migrations"
	sourceURL := fmt.Sprintf("file://%s/%s", startDir, migrationDir)

	m, err := migrate.NewWithDatabaseInstance(sourceURL, "go_metric", driver)
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
