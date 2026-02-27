package testutils

import (
	"context"
	"fmt"
	"sync"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var (
	pgContainer *postgres.PostgresContainer
	pgDBURL     string
	pgOnce      sync.Once
)

func SetupPostgres() string {
	pgOnce.Do(func() {
		ctx := context.Background()

		container, err := postgres.Run(ctx,
		"postgres:18-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		postgres.WithInitScripts("../testdata/schema.sql"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		panic(fmt.Errorf("failed to start postgres container: %w", err))
	}

		url, err := container.ConnectionString(ctx, "sslmode=disable")
		if err != nil {
			panic(err)
		}

		pgContainer = container
		pgDBURL = url
	})
	return pgDBURL
}