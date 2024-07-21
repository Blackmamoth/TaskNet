package db

import (
	"context"
	"fmt"

	"github.com/blackmamoth/tasknet/pkg/config"
	"github.com/jackc/pgx/v5"
)

var Conn *pgx.Conn

func init() {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=verify-full",
		config.GlobalConfig.CockroachDBConfig.COCKROACH_DB_USER,
		config.GlobalConfig.CockroachDBConfig.COCKROACH_DB_PASS,
		config.GlobalConfig.CockroachDBConfig.COCKROACH_DB_HOST,
		config.GlobalConfig.CockroachDBConfig.COCKROACH_DB_PORT,
		config.GlobalConfig.CockroachDBConfig.COCKROACH_DB_DBNAME,
	)

	conn, err := pgx.Connect(context.Background(), dsn)

	if err != nil {
		config.Logger.CRITICAL("Application disconnected from CockroachDB Server: %v", err)
	}

	if err := conn.Ping(context.Background()); err != nil {
		config.Logger.CRITICAL("Application disconnected from CockroachDB Server: %v", err)
	}

	Conn = conn

	config.Logger.INFO("Application connected to CockroachDB Server")
}
