package config

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
	"url-shortener/internal/db"
)

type Config struct {
	DatabaseURL string
	ServerPort  string
	Cassandra   db.Cassandra
}

func Load(logger *logrus.Logger) *Config {
	err := godotenv.Load()
	if err != nil {
		logger.Fatalf("Error loading config file (.env)")
	}
	logger.Info("Successfully loaded config file")
	return &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		ServerPort:  os.Getenv("SERVER_PORT"),
		Cassandra: db.Cassandra{
			Host:     os.Getenv("CASSANDRA_HOST"),
			Username: os.Getenv("CASSANDRA_USER"),
			Password: os.Getenv("CASSANDRA_PASSWORD"),
			Keyspace: os.Getenv("CASSANDRA_KEYSPACE"),
		},
	}
}
