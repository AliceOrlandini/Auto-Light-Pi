package config

import (
	"github.com/AliceOrlandini/Auto-Light-Pi/controllers"
	"github.com/AliceOrlandini/Auto-Light-Pi/services"
)

type Initialization struct {
	UserRepository services.UserRepository
	UserService controllers.UserService
	UserController *controllers.UserController
}

func NewInitialization(
	UserRepository services.UserRepository,
	UserService controllers.UserService,
	UserController *controllers.UserController,
) *Initialization {
	return &Initialization{
		UserRepository: UserRepository,
		UserService: UserService,
		UserController: UserController,
	}
}