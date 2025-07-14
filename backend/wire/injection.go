//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitializeServer() (*gin.Engine, error) {
    wire.Build(ProviderSet)
    return nil, nil
}