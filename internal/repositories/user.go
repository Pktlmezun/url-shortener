package repositories

import (
	"database/sql"
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

func (r *UserRepository) LoginUser(email string) (models.User, error) {
	query := `
        SELECT id, username, password, email, urlCounter FROM users WHERE email = $1;
    `
	var dbUser models.User
	err := r.DB.QueryRow(query, email).Scan(&dbUser.Id, &dbUser.Username, &dbUser.Password, &dbUser.Email, &dbUser.UrlCounter)
	if err != nil {
		r.Logger.Error(err)
		return models.User{}, err
	}

	return dbUser, nil
}
