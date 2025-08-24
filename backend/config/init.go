package config

import (
	"github.com/AliceOrlandini/Auto-Light-Pi/controllers"
	"github.com/AliceOrlandini/Auto-Light-Pi/services"
)

type Initialization struct {
	UserRepository services.UserRepository
	AuthService controllers.AuthService
	AuthController *controllers.AuthController
}

func NewInitialization(
	UserRepository services.UserRepository,
	AuthService controllers.AuthService,
	AuthController *controllers.AuthController,
) *Initialization {
	return &Initialization{
		UserRepository: UserRepository,
		AuthService: AuthService,
		AuthController: AuthController,
	}
}