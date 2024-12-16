package repositories

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"url-shortener/pkg/models"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) InsertUser(user models.User) (int64, error) {
	query := `
        INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id
    `
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return 0, err
	}

	var id int64 = 0
	err = r.DB.QueryRow(query, user.Username, user.Email, string(hashedPassword)).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func hashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}
