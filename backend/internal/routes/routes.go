package routes

import (
	"net/http"

	"github.com/AliceOrlandini/Auto-Light-Pi/internal/auth"
	"github.com/AliceOrlandini/Auto-Light-Pi/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(authController *auth.Controller) *gin.Engine {
	// create a new gin router
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// the main group is /api
	api := router.Group("/api")
	{
		api.POST("/register", authController.Register)
		api.POST("/login/email", authController.LoginByEmail)
		api.POST("/login/username", authController.LoginByUsername)
		api.POST("/refresh", authController.RefreshToken)

		// the auth group is for authenticated users only
		auth := api.Group("/")
		auth.Use(middleware.AuthMiddleware())
		{
			// endpoint to check if the user is authenticated
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