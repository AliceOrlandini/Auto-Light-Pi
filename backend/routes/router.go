package routes

import (
	"net/http"

	"github.com/AliceOrlandini/Auto-Light-Pi/config"
	"github.com/AliceOrlandini/Auto-Light-Pi/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(init *config.Initialization) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	api := router.Group("/api")
	{
		api.POST("/register", init.UserController.Register)
		api.POST("/login/email", init.UserController.LoginByEmail)
		api.POST("/login/username", init.UserController.LoginByUsername)

		auth := api.Group("/")
		auth.Use(middleware.AuthMiddleware())
		{
			auth.GET("/ping", func(c *gin.Context) {
        userID := c.GetString("userID")
        c.JSON(http.StatusOK, gin.H {
            "message": "Hello, user " + userID,
        })
			})
		}
	}

	return router
}