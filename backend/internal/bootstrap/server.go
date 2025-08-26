package bootstrap

import (
	"context"

	"github.com/AliceOrlandini/Auto-Light-Pi/internal/auth"
	"github.com/AliceOrlandini/Auto-Light-Pi/internal/refresh_token"
	"github.com/AliceOrlandini/Auto-Light-Pi/internal/routes"
	"github.com/AliceOrlandini/Auto-Light-Pi/internal/user"
	"github.com/gin-gonic/gin"
)

func InitializeServer(ctx context.Context) (*gin.Engine, error) {
    // Initialize DB connections
    if err := InitPosgresDB(ctx); err != nil {
        return nil, err
    }
    if err := InitRedisDB(ctx); err != nil {
        return nil, err
    }

    // Repositories
    userRepo := user.NewUserRepository(PostgresDB)
    refreshTokenRepo := refresh_token.NewRefreshTokenRepository(RedisDB)

    // Services
    authService := auth.NewAuthService(userRepo, refreshTokenRepo)

    // Controllers
    authController := auth.NewAuthController(authService)

    // Routes
    engine := routes.SetupRoutes(authController)
    return engine, nil
}
