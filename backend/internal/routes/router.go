package routes

import (
	"net/http"

	"github.com/AliceOrlandini/Auto-Light-Pi/internal/controllers"
	"github.com/AliceOrlandini/Auto-Light-Pi/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(authController *controllers.AuthController) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	api := router.Group("/api")
	{
		api.POST("/register", authController.Register)
		api.POST("/login/email", authController.LoginByEmail)
		api.POST("/login/username", authController.LoginByUsername)

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