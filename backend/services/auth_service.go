package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/AliceOrlandini/Auto-Light-Pi/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const refreshTokenVersion = "rt1"

var	(
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotExists 		 = errors.New("user not exists")
	ErrInvalidPassword   = errors.New("password not valid")
	ErrInvalidToken      = errors.New("invalid refresh token")
	ErrExpired           = errors.New("expired")
)

type UserRepository interface {
	CreateOne(ctx context.Context, user *models.User) error
	GetOneByEmail(ctx context.Context, email string) (*models.User, error)
	GetOneByUsername(ctx context.Context, username string) (*models.User, error)
}

type RefreshTokenRepository interface {
	CreateOne(ctx context.Context, refreshToken *models.RefreshToken) error
	GetOneUserIDByTokenHash(ctx context.Context, tokenHash string) (string, error)
	GetOneTokenHashByUserID(ctx context.Context, userID string) (string, error)
	RevokeOneByUserID(ctx context.Context, userID string) error
}

type AuthService struct {
	userRepo UserRepository
	refreshTokenRepo RefreshTokenRepository
}

func NewAuthService(userRepo UserRepository, refreshTokenRepo RefreshTokenRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		refreshTokenRepo: refreshTokenRepo,
	}
}

func (s *AuthService) Register(ctx context.Context, username string, email string, password []byte, name string, surname string) error {
	// first check if the email already exists
	emailAlreadyExists, err := s.userRepo.GetOneByEmail(ctx, email)
	if err != nil {
		// this will be an internal server error since it is a db error
		return err
	}
	if emailAlreadyExists != nil {
		return fmt.Errorf("%w: email already registered", ErrUserAlreadyExists)
	}

	// now do the same thing for the username
	usernameAlreadyExists, err := s.userRepo.GetOneByUsername(ctx, username)
	if err != nil {
		return err
	}
	if usernameAlreadyExists != nil {
		return fmt.Errorf("%w: username already registered", ErrUserAlreadyExists)
	}

	id := uuid.New()

	password, err = bcrypt.GenerateFromPassword(password, 14)
	if err != nil {
		// this is also an internal server error
		return err
	}

	user := &models.User{
		ID: id.String(),
		Username: username,
		Email: email,
		Password: password,
		Name: name,
		Surname: surname,
	}

	return s.userRepo.CreateOne(ctx, user)
}

func (s *AuthService) LoginByUsername(ctx context.Context, username string, password []byte) (*models.User, error) {
	user, err := s.userRepo.GetOneByUsername(ctx, username)
	// if there is an error than it is an internal server error
	// since is db releted 
	if err != nil {
		return nil, err
	}
	// if the user field is empty than the username does not exists in the db
	if user == nil {
		return nil, ErrUserNotExists
	}

	err = bcrypt.CompareHashAndPassword(user.Password, password)
	// if there is an error comparing the passwords, check if the error is 
	// a mismatched one or is an internal server error
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, ErrInvalidPassword
		}
		return nil, err
	}

	return user, nil
}

func (s *AuthService) LoginByEmail(ctx context.Context, email string, password []byte) (*models.User, error) {
	user, err := s.userRepo.GetOneByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotExists
	}

	err = bcrypt.CompareHashAndPassword(user.Password, password)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, ErrInvalidPassword
		}
		return nil, err
	}

	return user, nil
}

func (s *AuthService) GenerateJWT(user *models.User) (string, error) {
	var secret = []byte(os.Getenv("JWT_SECRET"))
	var appName = os.Getenv("APPLICATION_NAME")

  token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"iss": appName,
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	})

	return token.SignedString(secret)
}

func (s *AuthService) GenerateRefreshToken(ctx context.Context, userID string) (*models.RefreshToken, error) {
	// generate 32 random bytes that will be used as the refresh token
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}

	// the TTL of the token is 7 days
	TTL := time.Now().Add(7 * 24 * time.Hour)

	// the opaque is the base64 encoded version of the random bytes
	opaque := base64.RawURLEncoding.EncodeToString(bytes)
	
	// the version is used in case in the future we need to change the format of the token
	// in this way we can easily identify the version of the token
	token := fmt.Sprintf("%s.%s", refreshTokenVersion, opaque)
	
	// in the db we will not save the opaque value, but only the hash
	sum := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(sum[:])

	refreshToken := &models.RefreshToken{
		RefreshToken:     token,
		RefreshTokenHash: tokenHash,
		UserID:           userID,
		TTL:              TTL,
		CreatedAt:        time.Now(),
	}

	err = s.refreshTokenRepo.CreateOne(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	return refreshToken, nil
}

func (s *AuthService) hashRefreshToken(token string) (string, error) {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:]), nil
}

func isRefreshTokenSupported(token string) bool {
	return strings.HasPrefix(token, refreshTokenVersion+".")
}
