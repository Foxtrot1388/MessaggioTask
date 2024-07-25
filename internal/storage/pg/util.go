package storage

import (
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
)

func (con *dbStorage) migrateDB() (err error) {

	driver, err := postgres.WithInstance(con.DB.DB, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		con.cfg.DB,
		driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			con.log.Info("No migration up")
			return nil
		} else {
			return err
		}
	}

	return nil

}
