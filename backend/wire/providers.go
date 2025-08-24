package wire

import (
	"github.com/AliceOrlandini/Auto-Light-Pi/config"
	"github.com/AliceOrlandini/Auto-Light-Pi/controllers"
	"github.com/AliceOrlandini/Auto-Light-Pi/repositories"
	"github.com/AliceOrlandini/Auto-Light-Pi/routes"
	"github.com/AliceOrlandini/Auto-Light-Pi/services"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	ProvidePostgresDB,
	ProvideRedisDB,
	repositories.NewUserRepository,
	wire.Bind(new(services.UserRepository), new(*repositories.UserRepository)),
	repositories.NewRefreshTokenRepository,
	wire.Bind(new(services.RefreshTokenRepository), new(*repositories.RefreshTokenRepository)),
	services.NewAuthService,
	wire.Bind(new(controllers.AuthService), new(*services.AuthService)),
	controllers.NewAuthController,
	config.NewInitialization,
	routes.SetupRoutes,
)