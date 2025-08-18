//go:build wireinject
// +build wireinject

package wire

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitializeServer(ctx context.Context) (*gin.Engine, error) {
    wire.Build(ProviderSet)
    return nil, nil
}