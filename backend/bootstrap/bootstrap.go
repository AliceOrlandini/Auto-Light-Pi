package bootstrap

import (
	"context"

	"github.com/AliceOrlandini/Auto-Light-Pi/config"
	"github.com/AliceOrlandini/Auto-Light-Pi/controllers"
	"github.com/AliceOrlandini/Auto-Light-Pi/repositories"
	"github.com/AliceOrlandini/Auto-Light-Pi/routes"
	"github.com/AliceOrlandini/Auto-Light-Pi/services"
	"github.com/gin-gonic/gin"
)

func InitializeServer(ctx context.Context) (*gin.Engine, error) {
    // Initialize DB connections
    if err := config.InitPosgresDB(ctx); err != nil {
        return nil, err
    }
    if err := config.InitRedisDB(ctx); err != nil {
        return nil, err
    }

    // Repositories
    userRepo := repositories.NewUserRepository(config.PostgresDB)
    refreshTokenRepo := repositories.NewRefreshTokenRepository(config.RedisDB)

    // Services
    authService := services.NewAuthService(userRepo, refreshTokenRepo)

    // Controllers
    authController := controllers.NewAuthController(authService)

    // App initialization container
    initialization := config.NewInitialization(userRepo, authService, authController)

    // Routes
    engine := routes.SetupRoutes(initialization)
    return engine, nil
}
 
