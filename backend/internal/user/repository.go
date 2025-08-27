package user

import (
	"context"
	"database/sql"

	"errors"

	"github.com/google/uuid"
)

type UserEntity struct {
	ID       uuid.UUID
	Username string
	Email    string
	Password []byte
	Name     string
	Surname  string
}

type repository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *repository {
	return &repository{db: db}
}

func (r *repository) CreateOne(ctx context.Context, user *User) error {
	UserEntity, err := toEntity(user)
	if err != nil {
		// it is an internal server error
		return err
	}
	query := `
		INSERT INTO user_account(id, username, email, password, name, surname) 
		VALUES($1, $2, $3, $4, $5, $6)
	`
	_, err = r.db.ExecContext(ctx, query, UserEntity.ID, UserEntity.Username, UserEntity.Email, UserEntity.Password, UserEntity.Name, UserEntity.Surname)
	return err
}

func (r *repository) GetOneByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, username, email, password, name, surname
		FROM user_account 
		WHERE email = $1;
	`
	row := r.db.QueryRowContext(ctx, query, email)

	var user UserEntity
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Name, &user.Surname)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user.toUser(), nil
}

func (r *repository) GetOneByUsername(ctx context.Context, username string) (*User, error) {
	query := `
		SELECT id, username, email, password, name, surname
		FROM user_account 
		WHERE username = $1
	`
	row := r.db.QueryRowContext(ctx, query, username)

	var user UserEntity
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Name, &user.Surname)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user.toUser(), nil
}

func (ue *UserEntity) toUser() *User {
	return &User{
		ID:       ue.ID.String(),
		Username: ue.Username,
		Email:    ue.Email,
		Password: string(ue.Password),
		Name:     ue.Name,
		Surname:  ue.Surname,
	}
}

func toEntity(user *User) (*UserEntity, error) {
	var id uuid.UUID
	if user.ID == "" {
		id = uuid.New()
	} else {
		var err error
		id, err = uuid.Parse(user.ID)
		if err != nil {
			return nil, err
		}
	}
	return &UserEntity{
		ID:       id,
		Username: user.Username,
		Email:    user.Email,
		Password: []byte(user.Password),
		Name:     user.Name,
		Surname:  user.Surname,
	}, nil
}