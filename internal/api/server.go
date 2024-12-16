package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"url-shortener/internal/handlers"
	"url-shortener/internal/repositories"
	"url-shortener/internal/services"
)

func StartSever(port string, db *sql.DB) {

	userRepo := repositories.NewUserRepository(db)

	userService := services.NewUserService(userRepo)

	userHandler := handlers.NewUserHandler(userService)

	router := SetupRouter(userHandler)

	log.Println("Server is running on port ", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))

}
