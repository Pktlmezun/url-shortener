package config

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
)

type Config struct {
	DatabaseURL string
	ServerPort  string
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
	}
}
