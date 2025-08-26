package user

import (
	"context"
	"database/sql"

	"errors"

	"github.com/AliceOrlandini/Auto-Light-Pi/internal/models"
)

type repository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *repository {
	return &repository{db: db}
}

func (r *repository) CreateOne(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO user_account(id, username, email, password, name, surname) 
		VALUES($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Username, user.Email, user.Password, user.Name, user.Surname)
	return err
}

func (r *repository) GetOneByEmail(ctx context.Context, email string) (*models.User, error) {
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

func (r *repository) GetOneByUsername(ctx context.Context, username string) (*models.User, error) {
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