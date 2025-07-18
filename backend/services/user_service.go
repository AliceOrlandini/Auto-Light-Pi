package services

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/AliceOrlandini/Auto-Light-Pi/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var	(
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotExists 		 = errors.New("user not exists")
	ErrInvalidPassword   = errors.New("passowrd not valid")
)

type UserRepository interface {
	CreateOne(ctx context.Context, user *models.User) error
	GetOneByEmail(ctx context.Context, email string) (*models.User, error)
	GetOneByUsername(ctx context.Context, username string) (*models.User, error)
}

type userService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *userService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) Register(ctx context.Context, username string, email string, password []byte) error {
	// first check if the email already exists
	emailAlreadyExists, err := s.repo.GetOneByEmail(ctx, email)
	if err != nil {
		// this will be an internal server error since it is a db error
		return err
	}
	if emailAlreadyExists != nil {
		return fmt.Errorf("%w: email already registered", ErrUserAlreadyExists)
	}

	// now do the same thing for the username
	usernameAlreadyExists, err := s.repo.GetOneByUsername(ctx, username)
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
	}

	return s.repo.CreateOne(ctx, user)
}

func (s *userService) LoginByUsername(ctx context.Context, username string, password []byte) (*models.User, error) {
	user, err := s.repo.GetOneByUsername(ctx, username)
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

func (s *userService) LoginByEmail(ctx context.Context, email string, password []byte) (*models.User, error) {
	user, err := s.repo.GetOneByEmail(ctx, email)
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

func (s *userService) GenerateJWT(user *models.User) (string, error) {
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