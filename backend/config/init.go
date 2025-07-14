package config

import (
	"github.com/AliceOrlandini/Auto-Light-Pi/controllers"
	"github.com/AliceOrlandini/Auto-Light-Pi/repositories"
	"github.com/AliceOrlandini/Auto-Light-Pi/services"
)

type Initialization struct {
	UserRepository repositories.UserRepository
	UserService services.UserService
	UserController *controllers.UserController
}

func NewInitialization(
	UserRepository repositories.UserRepository,
	UserService services.UserService,
	UserController *controllers.UserController,
) *Initialization {
	return &Initialization{
		UserRepository: UserRepository,
		UserService: UserService,
		UserController: UserController,
	}
}