package server

import (
	"database/sql"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/sirupsen/logrus"
	"net/http"
	"url-shortener/internal/api"
	"url-shortener/internal/handlers"
	"url-shortener/internal/repositories"
	"url-shortener/internal/services"
)

func StartSever(port string, db *sql.DB, session *gocql.Session, logger *logrus.Logger) {

	userRepo := repositories.NewUserRepository(db, logger)
	userService := services.NewUserService(userRepo, logger)
	userHandler := handlers.NewUserHandler(userService, logger)

	URLRepo := repositories.NewURLRepository(session, logger, db)
	URLService := services.NewURLService(URLRepo, logger, db)
	URLHandler := handlers.NewURLHandler(URLService, logger)

	router := api.SetupRouter(userHandler, URLHandler, logger)

	logger.Info("Server is running on port ", port)

	logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))

}
