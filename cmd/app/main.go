package main

import (
	"url-shortener/config"
	"url-shortener/internal/db"
	"url-shortener/internal/server"
	"url-shortener/pkg/logging"

	_ "github.com/lib/pq"

	_ "github.com/jackc/pgx/v5"
)

func main() {

	logger := logging.Init()
	cfg := config.Load(logger)

	dbConn, err := db.ConnectPostgres(cfg.DatabaseURL, logger)
	if err != nil {
		logger.Fatal(err)
	}

	server.StartSever(cfg.ServerPort, dbConn, logger)

}
