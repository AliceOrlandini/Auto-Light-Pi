package repositories

import (
	"database/sql"

	"github.com/AliceOrlandini/Auto-Light-Pi/models"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByEmail(email string) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	query := "INSERT INTO user_account(id, username, email, password) VALUES($1, $2, $3, $4)"
	_, err := r.db.Exec(query, user.Id, user.Username, user.Email, user.Password)
	return err
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, username, email, password 
		FROM user_account 
		WHERE email = $1;
	`
	row := r.db.QueryRow(query, email)

	var user models.User
	err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Password)
	if err != nil {
		// gestisci gli erroriiiiiiii Kevin <3
		return nil, nil
	}
	return &user, nil
}

func (r *userRepository) GetByUsername(username string) (*models.User, error) {
	query := "SELECT id, username, email, password FROM user_account WHERE username = $1"
	row := r.db.QueryRow(query, username)

	var user models.User
	err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}