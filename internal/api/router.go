package api

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"url-shortener/internal/handlers"
)

func SetupRouter(userHandler *handlers.UserHandler, urlHandler *handlers.URLHandler, logger *logrus.Logger) http.Handler {
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
		logger.WithFields(logrus.Fields{
			"method":     r.Method,
			"endpoint":   "/login",
			"request_id": fmt.Sprintf("%d", os.Getpid()),
			"data":       r.Body,
		}).Info("Login endpoint hit")
		userHandler.LoginUser(w, r)
	})

	mux.HandleFunc("/url", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		logger.WithFields(logrus.Fields{
			"method":     r.Method,
			"endpoint":   "/url",
			"request_id": fmt.Sprintf("%d", os.Getpid()),
			"data":       r.Body,
		}).Info("url endpoint hit")
		urlHandler.AddURL(w, r)
	})

	mux.HandleFunc("/{short_url}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		logger.WithFields(logrus.Fields{
			"method":     r.Method,
			"endpoint":   "/{short_url}",
			"request_id": fmt.Sprintf("%d", os.Getpid()),
			"data":       r.Body,
		}).Info("short_url endpoint hit")
		short_url := r.URL.Path[1:]
		urlHandler.GetURL(w, r, short_url)
	})

	return mux
}
