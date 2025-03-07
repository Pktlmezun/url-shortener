package repositories

import (
	"database/sql"
	"errors"
	"url-shortener/pkg/models"

	"github.com/sirupsen/logrus"
)

type URLRepository struct {
	Logger *logrus.Logger
	DB     *sql.DB
}

func NewURLRepository(logger *logrus.Logger, db *sql.DB) *URLRepository {
	return &URLRepository{
		Logger: logger,
		DB:     db,
	}
}

func (r *URLRepository) AddUrl(url *models.AddUrl) error {
	query := `INSERT INTO urls (user_id, short_url, long_url) VALUES ($1, $2, $3)`
	_, err := r.DB.Exec(query, url.UserId, url.ShortUrl, url.LongUrl)
	if err != nil {
		r.Logger.Error(err)
		return err
	}
	r.Logger.Info("Successfully inserted URL into PostgreSQL")
	return nil
}

func (r *URLRepository) GetURL(url models.Url) (string, error) {
	query := `SELECT long_url FROM urls WHERE user_id = $1 AND short_url = $2`
	var longUrl string
	err := r.DB.QueryRow(query, url.UserId, url.ShortUrl).Scan(&longUrl)
	if err != nil {
		if err == sql.ErrNoRows {
			r.Logger.Error("Could not find URL")
			return "", errors.New("could not find URL")
		}
		r.Logger.Error(err)
		return "", err
	}

	r.Logger.Info("Successfully queried URL from PostgreSQL")
	return longUrl, nil
}

func (r *URLRepository) DeleteURL(url models.Url) error {
	query := `DELETE FROM urls WHERE user_id = $1 AND short_url = $2`
	_, err := r.DB.Exec(query, url.UserId, url.ShortUrl)
	if err != nil {
		r.Logger.Error("Error deleting the URL: ", err)
		return err
	}
	r.Logger.Info("Successfully deleted URL from PostgreSQL")
	return nil
}

func (r *URLRepository) GetMyURLs(userID int64) ([]models.Url, error) {
	query := `SELECT user_id, short_url, long_url FROM urls WHERE user_id = $1`
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		r.Logger.Error(err)
		return nil, err
	}
	defer rows.Close()

	var URLs []models.Url
	for rows.Next() {
		var url models.Url
		if err := rows.Scan(&url.UserId, &url.ShortUrl, &url.LongUrl); err != nil {
			r.Logger.Error(err)
			return nil, err
		}
		URLs = append(URLs, url)
	}

	if err := rows.Err(); err != nil {
		r.Logger.Error(err)
		return nil, err
	}

	r.Logger.Info("Successfully retrieved user URLs from PostgreSQL")
	return URLs, nil
}

func (r *URLRepository) GetCounter() (int, error) {
	query := `UPDATE counter SET counter = counter + 1 WHERE id = 1 RETURNING counter - 1`
	var counter int
	err := r.DB.QueryRow(query).Scan(&counter)
	if err != nil {
		r.Logger.Error(err)
		return 0, err
	}
	return counter, nil
}

func (r *URLRepository) IsShortURL(short_url string) bool {
	query := `SELECT COUNT(*) FROM urls WHERE short_url = $1`
	var count int
	err := r.DB.QueryRow(query, short_url).Scan(&count)

	if err != nil && err != sql.ErrNoRows {
		r.Logger.Error(err)
		return false
	}
	return count == 0
}
