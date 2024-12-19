package repositories

import (
	"database/sql"
	"errors"
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
        INSERT INTO urls (user_id, short_url, long_url) 
        VALUES (?, ?, ?)
    `
	err := r.Session.Query(query, url.UserId, url.ShortUrl, url.LongUrl).Exec()
	if err != nil {
		r.Logger.Error(err)
		return err
	}
	r.Logger.Info("Successfully queried url to cassandra")
	return nil
}

func (r *URLRepository) GetURL(url models.Url) (string, error) {
	query := `
	   SELECT long_url FROM urls WHERE user_id = ? AND short_url = ?
   `
	rows := r.Session.Query(query, url.UserId, url.ShortUrl).Iter()
	var longUrl string
	for rows.Scan(&longUrl) {
		break
	}
	if longUrl == "" {
		r.Logger.Error("Could not find url")
		return "", errors.New("could not find url")
	}
	r.Logger.Info("Successfully queried url to cassandra")
	return longUrl, nil
}

func (r *URLRepository) GetMyURLs(userID int64) ([]models.Url, error) {
	query := `
	   SELECT user_id, short_url, long_url FROM urls WHERE user_id = ?
   `
	rows := r.Session.Query(query, userID).Iter()
	var URLs []models.Url
	var url models.Url
	for rows.Scan(&url.UserId, &url.ShortUrl, &url.LongUrl) {
		URLs = append(URLs, url)
	}
	r.Logger.Info("Successfully queried find my urls to cassandra")
	return URLs, nil
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
