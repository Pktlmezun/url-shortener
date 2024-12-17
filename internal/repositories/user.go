package repositories

import (
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"url-shortener/pkg/models"
)

type UserRepository struct {
	DB     *sql.DB
	Logger *logrus.Logger
}

func NewUserRepository(db *sql.DB, logger *logrus.Logger) *UserRepository {
	return &UserRepository{DB: db, Logger: logger}
}

func (r *UserRepository) InsertUser(user models.User) (int64, error) {
	query := `
        INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id
    `
	fmt.Println("Repo", user)
	var id int64 = 0
	err := r.DB.QueryRow(query, user.Username, user.Email, user.Password).Scan(&id)
	if err != nil {
		r.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("error inserting user to db")
		return 0, err
	}
	r.Logger.Infof("Inserted new user with id: %d", id)
	return id, nil
}

func (r *UserRepository) LoginUser(user *models.User) (models.User, error) {
	query := `
        SELECT id, username, password, email, urlCounter FROM users WHERE email = $1;
    `

	//var dbPassword string
	err := r.DB.QueryRow(query, user.Email).Scan(&user.Id, &user.Username, &user.Password, &user.Email, &user.UrlCounter)
	fmt.Println("DB RES:", user)
	if err != nil {
		r.Logger.Error(err)
		return models.User{}, err
	}

	return *user, nil
}
