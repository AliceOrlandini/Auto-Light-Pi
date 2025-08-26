package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/AliceOrlandini/Auto-Light-Pi/internal/models"
	"github.com/gin-gonic/gin"
)

type authService interface {
	Register(ctx context.Context, username string, email string, password []byte, name string, surname string) error
	LoginByUsername(ctx context.Context, username string, password []byte) (*models.User, error)
	LoginByEmail(ctx context.Context, email string, password []byte) (*models.User, error)
	GenerateJWT(userID string) (string, error)
	GenerateRefreshToken(ctx context.Context, userID string) (*models.RefreshToken, error)
	ValidateRefreshToken(ctx context.Context, token string) (string, error)
	RotateRefreshToken(ctx context.Context, userID string) (*models.RefreshToken, error)
}

type Controller struct {
	service authService
}

func NewAuthController(service authService) *Controller {
	return &Controller{service: service}
}

type registerRequest struct {
	Username string `json:"username" binding:"required,max=50"`
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name string `json:"name" binding:"required,max=50"`
	Surname string `json:"surname" binding:"required,max=50"`
}

type loginByUsernameRequest struct {
	Username string `json:"username" binding:"required,max=50"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginByEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func (uc *Controller) Register(c *gin.Context) {
	ctx := c.Request.Context()
	var request registerRequest
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
	err = uc.service.Register(ctx, 
		request.Username, 
		request.Email, 
		[]byte(request.Password),
		request.Name,
		request.Surname,
	)
	if err != nil {
		if errors.Is(err, ErrUserAlreadyExists) {
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
		// in case of a db error, log the error and return internal server error
		fmt.Printf("Register error: %v\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}

func (uc *Controller) LoginByUsername(c *gin.Context) {
	ctx := c.Request.Context()
	var request loginByUsernameRequest
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

	user, err := uc.service.LoginByUsername(ctx, request.Username, []byte(request.Password))
	if err != nil {
		if errors.Is(err, ErrUserNotExists) || errors.Is(err, ErrInvalidPassword) {
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

	accessToken, err := uc.service.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.SetCookie(
		"jwt",				// name
		accessToken,	// token
		3600,					// validity
		"/",					// path
		"",						// domain
		true,					// secure (HTTPS)
		true,					// httpOnly
	)

	refreshToken, err := uc.service.GenerateRefreshToken(ctx, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.SetCookie(
		"__Host-refresh_token",												// name
		refreshToken.RefreshToken,  									// value
		int(time.Until(refreshToken.TTL).Seconds()),	// validity
		"/",																					// path
		"",    																				// Empty Domain
		true,  																				// Secure (HTTPS)
		true,  																				// HttpOnly
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"user": user,
	})
}

func (uc *Controller) LoginByEmail(c *gin.Context) {
	ctx := c.Request.Context()
	var request loginByEmailRequest
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

	user, err := uc.service.LoginByEmail(ctx, request.Email, []byte(request.Password))
	if err != nil {
		if errors.Is(err, ErrUserNotExists) || errors.Is(err, ErrInvalidPassword) {
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

	tokenString, err := uc.service.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// REFACTORING: send the access token in the body and the refresh in the cookie
	// this is to avoid CSRF attacks
	c.SetCookie(
		"jwt",				// name
		tokenString,	// token
		3600,					// validity
		"/",					// path
		"",						// domain
		false,				// secure (HTTPS)
		true,					// httpOnly
	)

	refreshToken, err := uc.service.GenerateRefreshToken(ctx, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.SetCookie(
		"__Host-refresh_token",												// name
		refreshToken.RefreshToken,  									// value
		int(time.Until(refreshToken.TTL).Seconds()),	// validity
		"/",																					// path
		"",    																				// Empty Domain
		true,  																				// Secure (HTTPS)
		true,  																				// HttpOnly
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"user": user,
	})
}

func (uc *Controller) RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()
	const refreshCookie = "__Host-refresh_token"

	// retrieve the refresh token from the cookie
	rt, err := c.Cookie(refreshCookie)
	if err != nil || rt == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing refresh token"})
		return
	}

	// validate the refresh token
	userID, err := uc.service.ValidateRefreshToken(ctx, rt)
	if err != nil {
		if errors.Is(err, ErrInvalidToken) {
			// this two deletes the refresh token cookie in a secure way
			// with this we ensure that the browser makes another request 
			// with the same invalid token
			c.SetCookie(refreshCookie, "", -1, "/", "", true, true)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// generate a new refresh token by rotate
	newRefreshToken, err := uc.service.RotateRefreshToken(ctx, userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.SetCookie(
		"__Host-refresh_token",													// name
		newRefreshToken.RefreshToken,  									// value
		int(time.Until(newRefreshToken.TTL).Seconds()),	// validity
		"/",																						// path
		"",    																					// Empty Domain
		true,  																					// Secure (HTTPS)
		true,  																					// HttpOnly
	)

	// generate a new JWT
	tokenString, err := uc.service.GenerateJWT(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.SetCookie(
		"jwt",				// name
		tokenString,	// token
		3600,					// validity
		"/",					// path
		"",						// domain
		false,				// secure (HTTPS)
		true,					// httpOnly
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "refresh token generated",
	})
}

func (uc *Controller) validatePassword(password string) error {
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