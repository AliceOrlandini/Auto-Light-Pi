package bootstrap

import (
	"context"
	"database/sql"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

var PostgresDB *sql.DB
var RedisDB *redis.Client

func InitRedisDB(ctx context.Context) error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	dbDsn := os.Getenv("REDIS_DB_DSN")

	opt, err := redis.ParseURL(dbDsn)
	if err != nil {
		return err
	}

	rdb := redis.NewClient(opt)

	err = rdb.Ping(ctx).Err()
	if err != nil {
		return err
	}

	RedisDB = rdb

	return nil
}

func InitPosgresDB(ctx context.Context) error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	dbDsn := os.Getenv("POSTGRES_DB_DSN")

	psqlDB, err := sql.Open("postgres", dbDsn)
	if err != nil {
		return err
	}

	err = psqlDB.PingContext(ctx)
	if err != nil {
		return err
	}

	PostgresDB = psqlDB

	// read from the filesystem the database schema and execute it
	content, err := os.ReadFile("schema.sql")
  if err != nil {
    return err
  }
	_, err = PostgresDB.Exec(string(content))
	if err != nil {
		return err
	}

	return nil
}