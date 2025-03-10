package services

import (
	"database/sql"
	"errors"
	"strconv"
	"url-shortener/internal/repositories"
	"url-shortener/pkg/models"

	"github.com/deatil/go-encoding/base62"
	"github.com/sirupsen/logrus"
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

func (s *URLService) AddURL(url *models.AddUrl) (string, error) {
	if url.LongUrl == "" {
		s.Logger.Error("empty url")
		return "", errors.New("url is empty")
	}
	counter, err := s.URLRepo.GetCounter()
	if err != nil {
		return "", err
	}
	if url.CustomAllias != "" && s.URLRepo.IsShortURL(url.CustomAllias) {
		url.ShortUrl = url.CustomAllias
	} else {
		url.ShortUrl = generateShortUrl(counter)
	}
	err = s.URLRepo.AddUrl(url)
	if err != nil {
		return "", err
	}
	return url.ShortUrl, nil
}

func (s *URLService) GetURL(url models.Url) (string, error) {
	if url.ShortUrl == "" {
		s.Logger.Error("empty url")
		return "", errors.New("url is empty")
	}
	return s.URLRepo.GetURL(url)
}

func (s *URLService) DeleteURL(url models.Url) error {
	if url.ShortUrl == "" {
		s.Logger.Error("empty short url")
		return errors.New("short url is empty")
	}
	return s.URLRepo.DeleteURL(url)
}

func (s *URLService) GetMyURLs(userID int64) ([]models.Url, error) {
	if userID == 0 {
		s.Logger.Error("0 userID")
		return nil, errors.New("userID is empty")
	}
	return s.URLRepo.GetMyURLs(userID)
}

func generateShortUrl(counter int) string {
	base62Chars := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	encoder := base62.NewEncoding(base62Chars)

	numStr := strconv.Itoa(counter)
	encodedStr := string(encoder.Encode([]byte(numStr)))

	return encodedStr
}
