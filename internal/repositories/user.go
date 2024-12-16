package repositories

import (
	"database/sql"
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

	var id int64 = 0
	err := r.DB.QueryRow(query, user.Username, user.Email, user.Password).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
