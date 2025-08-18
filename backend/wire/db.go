package wire

import (
	"context"
	"database/sql"

	"github.com/AliceOrlandini/Auto-Light-Pi/config"
	"github.com/redis/go-redis/v9"
)

func ProvidePostgresDB(ctx context.Context) (*sql.DB, error) {
	err := config.InitPosgresDB(ctx)
	if err != nil {
		return nil, err
	}

	return config.PostgresDB, nil
}

func ProvideRedisDB(ctx context.Context) (*redis.Client, error) {
	err := config.InitRedisDB(ctx)
	if err != nil {
		return nil, err
	}

	return config.RedisDB, nil
}