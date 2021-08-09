package database

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func Migrate(db *sql.DB, to int) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:///some-http-server/internal/database/migrations",
		viper.GetString("svc_name"), driver)
	if err != nil {
		return err
	}

	err = m.Steps(to)
	if err != nil {
		return err
	}

	return nil
}
