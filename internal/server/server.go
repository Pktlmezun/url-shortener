package server

import (
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"url-shortener/internal/api"
	"url-shortener/internal/handlers"
	"url-shortener/internal/repositories"
	"url-shortener/internal/services"
)

func StartSever(port string, db *sql.DB, logger *logrus.Logger) {

	userRepo := repositories.NewUserRepository(db, logger)

	userService := services.NewUserService(userRepo, logger)

	userHandler := handlers.NewUserHandler(userService, logger)

	router := api.SetupRouter(userHandler, logger)

	logger.Info("Server is running on port ", port)

	logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))

}
