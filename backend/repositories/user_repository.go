package repositories

import (
	"context"
	"database/sql"

	"errors"

	"github.com/AliceOrlandini/Auto-Light-Pi/models"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserInvalidCredentials = errors.New("invalid user credentials")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateOne(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO user_account(id, username, email, password, name, surname) 
		VALUES($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Username, user.Email, user.Password, user.Name, user.Surname)
	return err
}

func (r *UserRepository) GetOneByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, username, email, password, name, surname
		FROM user_account 
		WHERE email = $1;
	`
	row := r.db.QueryRowContext(ctx, query, email)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Name, &user.Surname)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetOneByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
		SELECT id, username, email, password, name, surname
		FROM user_account 
		WHERE username = $1
	`
	row := r.db.QueryRowContext(ctx, query, username)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Name, &user.Surname)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}