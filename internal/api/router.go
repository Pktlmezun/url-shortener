package api

import (
	"fmt"
	"net/http"
	"os"
	"url-shortener/internal/handlers"

	"github.com/sirupsen/logrus"
)

func SetupRouter(userHandler *handlers.UserHandler, urlHandler *handlers.URLHandler, logger *logrus.Logger) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("./static")))

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

	mux.HandleFunc("/my_urls", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		logger.WithFields(logrus.Fields{
			"method":     r.Method,
			"endpoint":   "/my_urls",
			"request_id": fmt.Sprintf("%d", os.Getpid()),
			"data":       r.Body,
		}).Info("get my urls endpoint hit")
		urlHandler.GetMyURLs(w, r)
	})

	mux.HandleFunc("/{short_url}", func(w http.ResponseWriter, r *http.Request) {
		shortUrl := r.URL.Path[1:]
		if r.Method == http.MethodGet {
			logger.WithFields(logrus.Fields{
				"method":     r.Method,
				"endpoint":   "/{short_url}",
				"request_id": fmt.Sprintf("%d", os.Getpid()),
				"data":       r.Body,
			}).Info("short_url endpoint hit")
			urlHandler.GetURL(w, r, shortUrl)
		} else if r.Method == http.MethodDelete {
			logger.WithFields(logrus.Fields{
				"method":     r.Method,
				"endpoint":   "/url",
				"request_id": fmt.Sprintf("%d", os.Getpid()),
				"data":       r.Body,
			}).Info("url endpoint hit")
			urlHandler.DeleteURL(w, r, shortUrl)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	return mux
}
