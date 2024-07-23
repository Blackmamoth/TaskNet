package db

import (
	"context"
	"fmt"

	"github.com/blackmamoth/tasknet/pkg/config"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func init() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.GlobalConfig.RedisDBConfig.REDIS_DB_HOST, config.GlobalConfig.RedisDBConfig.REDIS_DB_PORT),
		Password: config.GlobalConfig.RedisDBConfig.REDIS_DB_PASS,
		DB:       0,
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			config.Logger.INFO("Application connected to Redis Server")
			return nil
		},
	})

	if status := RedisClient.Ping(context.Background()); status.Err() != nil {
		config.Logger.CRITICAL("Application disconnected from Redis Server: %v", status.Err())
	}
}
