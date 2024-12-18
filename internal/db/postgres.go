package db

import (
	"database/sql"
	"github.com/sirupsen/logrus"
)

func ConnectPostgres(dbURL string, logger *logrus.Logger) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		logger.Fatalf("Error connecting to postgres: %v", err)
		return nil, err
	}
	logger.Info("Successfully connected to postgres")
	return db, nil
}

func migratePostgres(db *sql.DB, logger *logrus.Logger) {

}
