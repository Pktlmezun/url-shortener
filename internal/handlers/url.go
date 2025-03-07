package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"url-shortener/internal/auth"
	"url-shortener/internal/services"
	"url-shortener/pkg/models"

	"github.com/sirupsen/logrus"
)

type URLHandler struct {
	Service *services.URLService
	Logger  *logrus.Logger
}

func NewURLHandler(service *services.URLService, logger *logrus.Logger) *URLHandler {
	return &URLHandler{
		Service: service,
		Logger:  logger,
	}
}

func (h *URLHandler) AddURL(w http.ResponseWriter, r *http.Request) {
	token, err := splitToken(r, h.Logger)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid token"))
		return
	}

	id, err := auth.VerifyToken(token)

	if err != nil {
		h.Logger.Error("Invalid token")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid token"))
		return
	}

	var URL models.AddUrl
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&URL)
	if err != nil {
		h.Logger.Error("Invalid data")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	URL.UserId = id
	shortUrl, err := h.Service.AddURL(&URL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	response := map[string]string{"message": "URL added", "short_url": shortUrl}
	_ = json.NewEncoder(w).Encode(response)
}

func (h *URLHandler) GetURL(w http.ResponseWriter, r *http.Request, shortUrl string) {
	token, err := splitToken(r, h.Logger)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid token"))
		return
	}

	id, err := auth.VerifyToken(token)

	if err != nil {
		h.Logger.Error("Invalid token")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid token"))
		return
	}

	var URL models.Url
	URL.UserId = id
	URL.ShortUrl = shortUrl
	longURL, err := h.Service.GetURL(URL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	response := map[string]string{"message": "URL retrieved", "long_url": longURL}
	_ = json.NewEncoder(w).Encode(response)
	//http.Redirect(w, r, longURL, http.StatusMovedPermanently)
}

func (h *URLHandler) DeleteURL(w http.ResponseWriter, r *http.Request, shortUrl string) {
	token, err := splitToken(r, h.Logger)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid token"))
		return
	}

	id, err := auth.VerifyToken(token)

	if err != nil {
		h.Logger.Error("Invalid token")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid token"))
		return
	}
	var URL models.Url
	URL.UserId = id
	URL.ShortUrl = shortUrl
	err = h.Service.DeleteURL(URL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	response := map[string]string{"message": "URL deleted"}
	_ = json.NewEncoder(w).Encode(response)
}

func (h *URLHandler) GetMyURLs(w http.ResponseWriter, r *http.Request) {
	token, err := splitToken(r, h.Logger)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid token"))
		return
	}

	id, err := auth.VerifyToken(token)

	if err != nil {
		h.Logger.Error("Invalid token")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid token"))
		return
	}

	urls, err := h.Service.GetMyURLs(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	_ = json.NewEncoder(w).Encode(urls)

}

func splitToken(r *http.Request, logger *logrus.Logger) (string, error) {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer")
	if len(splitToken) != 2 {
		logger.Error("invalid token format")
		return "", errors.New("invalid token format")
	}
	token := strings.TrimSpace(splitToken[1])
	return token, nil
}
