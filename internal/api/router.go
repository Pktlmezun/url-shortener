package api

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"url-shortener/internal/handlers"
)

func SetupRouter(userHandler *handlers.UserHandler, logger *logrus.Logger) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

		logger.WithFields(logrus.Fields{
			"method":     r.Method,
			"endpoint":   "/signup",
			"request_id": fmt.Sprintf("%d", os.Getpid()),
			"data":       r.Body,
		}).Info("Signup endpoint hit")

		userHandler.CreateUser(w, r)
	})

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		userHandler.LoginUser(w, r)
	})

	return mux
}
