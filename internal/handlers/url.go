package handlers

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"url-shortener/internal/auth"
	"url-shortener/internal/services"
	"url-shortener/pkg/models"
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
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer")
	if len(splitToken) != 2 {
		h.Logger.Error("Invalid token format")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid token format"))
		return
	}
	token := strings.TrimSpace(splitToken[1])

	id, err := auth.VerifyToken(token)

	if err != nil {
		h.Logger.Error("Invalid token")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid token"))
		return
	}

	var URL models.Url
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
