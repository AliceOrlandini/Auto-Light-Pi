package controllers

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/AliceOrlandini/Auto-Light-Pi/services"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	Service services.UserService
}

func NewUserController(service services.UserService) *UserController {
	return &UserController{Service: service}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,max=50"`
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginByUsernameRequest struct {
	Username string `json:"username" binding:"required,max=50"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginByEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func (uc *UserController) Register(c *gin.Context) {
	var request RegisterRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = uc.validatePassword(request.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = uc.Service.Register(request.Username, request.Email, []byte(request.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}

func (uc *UserController) LoginByUsername(c *gin.Context) {
	var request LoginByUsernameRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = uc.validatePassword(request.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := uc.Service.LoginByUsername(request.Username, []byte(request.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	tokenString, err := uc.Service.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	maxAge := 3600
	c.SetCookie(
		"jwt",				// name
		tokenString,	// token
		maxAge,				// validity
		"/",					// path
		"",						// domain
		false,				// secure (HTTPS)
		true,					// httpOnly
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"user": user,
	})
}

func (uc *UserController) LoginByEmail(c *gin.Context) {
	var request LoginByEmailRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = uc.validatePassword(request.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := uc.Service.LoginByEmail(request.Email, []byte(request.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	tokenString, err := uc.Service.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	maxAge := 3600
	c.SetCookie(
		"jwt",				// name
		tokenString,	// token
		maxAge,				// validity
		"/",					// path
		"",						// domain
		false,				// secure (HTTPS)
		true,					// httpOnly
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"user": user,
	})
}

func (uc *UserController) validatePassword(password string) error {
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	if !hasUpper {
		return errors.New("password must contain at least capital letter")
	}
	if !hasNumber {
		return errors.New("password must contain at least a number")
	}

	return nil
}