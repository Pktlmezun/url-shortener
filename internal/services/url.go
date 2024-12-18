package services

import (
	"database/sql"
	"errors"
	"github.com/deatil/go-encoding/base62"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"strconv"
	"url-shortener/internal/repositories"
	"url-shortener/pkg/models"
)

type URLService struct {
	URLRepo *repositories.URLRepository
	Logger  *logrus.Logger
}

func NewURLService(urlRepo *repositories.URLRepository, logger *logrus.Logger, db *sql.DB) *URLService {
	return &URLService{
		URLRepo: urlRepo,
		Logger:  logger,
	}
}

func (s *URLService) AddURL(url *models.Url) (string, error) {
	if url.LongUrl == "" {
		s.Logger.Error("empty url")
		return "", errors.New("url is empty")
	}
	counter, err := s.URLRepo.GetCounter()
	if err != nil {
		return "", err
	}
	url.ShortUrl = generateShortUrl(counter)
	url.Id = generateUUID()
	err = s.URLRepo.AddUrl(url)
	if err != nil {
		return "", err
	}

	return url.ShortUrl, nil
}

func generateShortUrl(counter int) string {
	base62Chars := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	encoder := base62.NewEncoding(base62Chars)

	numStr := strconv.Itoa(counter)
	encodedStr := string(encoder.Encode([]byte(numStr)))

	return encodedStr
}

func generateUUID() gocql.UUID {
	newUUID := uuid.New()
	gocqlUUID := gocql.UUID(newUUID)
	return gocqlUUID
}
