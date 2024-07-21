package main

import (
	"fmt"
	"os"

	"github.com/blackmamoth/tasknet/pkg/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/cockroachdb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	dsn := fmt.Sprintf("cockroachdb://%s:%s@%s:%s/%s?sslmode=verify-full",
		config.GlobalConfig.CockroachDBConfig.COCKROACH_DB_USER,
		config.GlobalConfig.CockroachDBConfig.COCKROACH_DB_PASS,
		config.GlobalConfig.CockroachDBConfig.COCKROACH_DB_HOST,
		config.GlobalConfig.CockroachDBConfig.COCKROACH_DB_PORT,
		config.GlobalConfig.CockroachDBConfig.COCKROACH_DB_DBNAME,
	)
	m, err := migrate.New(
		"file://cmd/migrate/migrations",
		dsn,
	)

	if err != nil {
		config.Logger.CRITICAL(err.Error())
	}

	v, d, _ := m.Version()
	config.Logger.INFO("Version: %d, dirty: %v", v, d)

	cmd := os.Args[len(os.Args)-1]
	if cmd == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			config.Logger.CRITICAL(err.Error())
		}
	}
	if cmd == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			config.Logger.CRITICAL(err.Error())
		}
	}
}
