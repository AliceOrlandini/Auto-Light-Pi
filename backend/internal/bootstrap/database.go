package bootstrap

import (
	"context"
	"database/sql"
	"log/slog"
	"os"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

var PostgresDB *sql.DB
var RedisDB *redis.Client

func InitRedisDB(ctx context.Context) error {
	// load the environment to create the redis dsn
	redisDBHost := os.Getenv("REDIS_HOST")
	redisDBPort := os.Getenv("REDIS_PORT")
	dbDsn := "redis://" + redisDBHost + ":" + redisDBPort
	opt, err := redis.ParseURL(dbDsn)
	if err != nil {
		slog.Error("failed to parse RedisDB URL", "error", err)
		return err
	}

	// create the redis client
	rdb := redis.NewClient(opt)

	// check if the redis server is reachable
	err = rdb.Ping(ctx).Err()
	if err != nil {
		slog.Error("failed to connect to RedisDB", "error", err)
		return err
	}

	// set the redis client to the global variable
	RedisDB = rdb

	return nil
}

func InitPosgresDB(ctx context.Context) error {
	// load the environment to create the postgres dsn
	postgresDBHost := os.Getenv("POSTGRES_HOST")
	postgresDBUser := os.Getenv("POSTGRES_USER")
	postgresDBPass := os.Getenv("POSTGRES_PASSWORD")
	postgresDBName := os.Getenv("POSTGRES_DB")
	postgresDBPort := os.Getenv("POSTGRES_PORT")
	dbDsn := "host=" + postgresDBHost + " user=" + postgresDBUser + " password=" + postgresDBPass + " dbname=" + postgresDBName + " port=" + postgresDBPort + " sslmode=disable"
	psqlDB, err := sql.Open("postgres", dbDsn)
	if err != nil {
		slog.Error("failed to open connection with PostgresDB", "error", err)
		return err
	}

	// check if the postgres server is reachable
	err = psqlDB.PingContext(ctx)
	if err != nil {
		slog.Error("failed to connect to PostgresDB", "error", err)
		return err
	}

	// set the postgres client to the global variable
	PostgresDB = psqlDB

	return nil
}