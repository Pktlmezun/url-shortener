package repositories

import (
	"database/sql"
	"github.com/gocql/gocql"
	"github.com/sirupsen/logrus"
	"url-shortener/pkg/models"
)

type URLRepository struct {
	Session *gocql.Session
	Logger  *logrus.Logger
	DB      *sql.DB
}

func NewURLRepository(session *gocql.Session, logger *logrus.Logger, db *sql.DB) *URLRepository {
	return &URLRepository{
		Session: session,
		Logger:  logger,
		DB:      db,
	}
}

func (r *URLRepository) AddUrl(url *models.Url) error {
	query := `
        INSERT INTO urls (id, user_id, short_url, long_url) 
        VALUES (?, ?, ?, ?)
    `
	err := r.Session.Query(query, url.Id, url.UserId, url.ShortUrl, url.LongUrl).Exec()
	if err != nil {
		r.Logger.Error(err)
		return err
	}
	r.Logger.Info("Successfully queried url to cassandra")
	return nil
}

func (r *URLRepository) GetCounter() (int, error) {
	query := `UPDATE counter
	SET counter = counter + 1
	WHERE id = 0
	RETURNING counter - 1`
	var counter int
	err := r.DB.QueryRow(query).Scan(&counter)
	if err != nil {
		r.Logger.Error(err)
		return 0, err
	}
	return counter, nil
}
