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
	services.NewUserService,
	wire.Bind(new(controllers.UserService), new(*services.UserService)),
	controllers.NewUserController,
	config.NewInitialization,
	routes.SetupRoutes,
)