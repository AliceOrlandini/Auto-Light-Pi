package services

import (
	"errors"
	"os"
	"time"

	"github.com/AliceOrlandini/Auto-Light-Pi/models"
	"github.com/AliceOrlandini/Auto-Light-Pi/repositories"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

/*const (
	ErrUserNotFound = errors.New("Error not found")
)*/


type UserService interface {
	Register(username string, email string, password []byte) error
	LoginByUsername(username string, password []byte) (*models.User, error) 
	LoginByEmail(email string, password []byte) (*models.User, error) 
	GenerateJWT(user *models.User) (string, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) Register(username string, email string, password []byte) error {
	existing, _ := s.repo.GetByEmail(email)
	if existing != nil {
		// return fmt.Errof("%w: user id %s", ErrUserNotFound, email)
		return errors.New("email already existing")
	}
	existing, _ = s.repo.GetByUsername(username)
	if existing != nil {
		return errors.New("username already existing")
	}

	id := uuid.New()

	password, err := bcrypt.GenerateFromPassword(password, 14)
	if err != nil {
		return err
	}

	user := &models.User{
		Id: id,
		Username: username,
		Email: email,
		Password: password,
	}

	return s.repo.Create(user)
}

func (s *userService) LoginByUsername(username string, password []byte) (*models.User, error) {
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return nil, errors.New("user not exists: " + err.Error())
	}

	err = bcrypt.CompareHashAndPassword(user.Password, password)
	if err != nil {
		return nil, errors.New("invalid password: " + err.Error())
	}

	return user, nil
}

func (s *userService) LoginByEmail(email string, password []byte) (*models.User, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, errors.New("user not exists: " + err.Error())
	}

	err = bcrypt.CompareHashAndPassword(user.Password, password)
	if err != nil {
		return nil, errors.New("invalid password: " + err.Error())
	}

	return user, nil
}

func (s *userService) GenerateJWT(user *models.User) (string, error) {
	var secret = []byte(os.Getenv("JWT_SECRET"))
	var appName = os.Getenv("APPLICATION_NAME")

  token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Id.String(),
		"iss": appName,
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	})

	return token.SignedString(secret)
}