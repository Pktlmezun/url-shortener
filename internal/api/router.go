package api

import (
	"net/http"
	"url-shortener/internal/handlers"
)

func SetupRouter(userHandler *handlers.UserHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		userHandler.CreateUser(w, r)
	})

	return mux
}
