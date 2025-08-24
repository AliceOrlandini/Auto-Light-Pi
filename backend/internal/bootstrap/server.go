package bootstrap

import (
	"context"

	"github.com/AliceOrlandini/Auto-Light-Pi/internal/controllers"
	"github.com/AliceOrlandini/Auto-Light-Pi/internal/repositories"
	"github.com/AliceOrlandini/Auto-Light-Pi/internal/routes"
	"github.com/AliceOrlandini/Auto-Light-Pi/internal/services"
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
    userRepo := repositories.NewUserRepository(PostgresDB)
    refreshTokenRepo := repositories.NewRefreshTokenRepository(RedisDB)

    // Services
    authService := services.NewAuthService(userRepo, refreshTokenRepo)

    // Controllers
    authController := controllers.NewAuthController(authService)

    // Routes
    engine := routes.SetupRoutes(authController)
    return engine, nil
}
