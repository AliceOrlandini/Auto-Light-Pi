package testutils

import (
	"context"
	"sync"

	redisContainer "github.com/testcontainers/testcontainers-go/modules/redis"
)

var (
	rContainer *redisContainer.RedisContainer
	redisDBURL       string
	redisOnce        sync.Once
)

func SetupRedis() string {
	redisOnce.Do(func() {
		ctx := context.Background()
		container, err := redisContainer.Run(ctx,
			"redis:alpine",
			redisContainer.WithSnapshotting(0, 0),
			redisContainer.WithLogLevel(redisContainer.LogLevelVerbose),
		)
		if err != nil {
			panic(err)
		}
		
		url, err := container.ConnectionString(ctx)
		if err != nil {
			panic(err)
		}

		rContainer = container
		redisDBURL = url
	})
	return redisDBURL
}