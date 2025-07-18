package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/AliceOrlandini/Auto-Light-Pi/models"
	"github.com/AliceOrlandini/Auto-Light-Pi/services"
	"github.com/gin-gonic/gin"
)

type UserService interface {
	Register(ctx context.Context, username string, email string, password []byte) error
	LoginByUsername(ctx context.Context, username string, password []byte) (*models.User, error) 
	LoginByEmail(ctx context.Context, email string, password []byte) (*models.User, error) 
	GenerateJWT(user *models.User) (string, error)
}

type UserController struct {
	Service UserService
}

func NewUserController(service UserService) *UserController {
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
	ctx := c.Request.Context()
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
	// I want to check:
	// 1. if the user already exists (from the service layer)
	// 2. db errors 
	// 3. context cancelled or timeouted
	err = uc.Service.Register(ctx, request.Username, request.Email, []byte(request.Password))
	if err != nil {
		if errors.Is(err, services.ErrUserAlreadyExists) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, context.Canceled) {
			// in this case the client disconnected so we do not return nothing
			return
		}
		if errors.Is(err, context.DeadlineExceeded) {
			c.JSON(http.StatusRequestTimeout, gin.H{"error": "request timeout"})
			return
		}
		// in case of a db error, I log the error and i return internal server error
		fmt.Printf("Register error: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}

func (uc *UserController) LoginByUsername(c *gin.Context) {
	ctx := c.Request.Context()
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

	user, err := uc.Service.LoginByUsername(ctx, request.Username, []byte(request.Password))
	if err != nil {
		if errors.Is(err, services.ErrUserNotExists) || errors.Is(err, services.ErrInvalidPassword) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		if errors.Is(err, context.Canceled) {
			return
		}
		if errors.Is(err, context.DeadlineExceeded) {
			c.JSON(http.StatusRequestTimeout, gin.H{"error": "request timeout"})
			return
		}
		fmt.Printf("Login by username error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	tokenString, err := uc.Service.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
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
	ctx := c.Request.Context()
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

	user, err := uc.Service.LoginByEmail(ctx, request.Email, []byte(request.Password))
	if err != nil {
		if errors.Is(err, services.ErrUserNotExists) || errors.Is(err, services.ErrInvalidPassword) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		if errors.Is(err, context.Canceled) {
			return
		}
		if errors.Is(err, context.DeadlineExceeded) {
			c.JSON(http.StatusRequestTimeout, gin.H{"error": "request timeout"})
			return
		}
		fmt.Printf("Login by email error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	tokenString, err := uc.Service.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
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