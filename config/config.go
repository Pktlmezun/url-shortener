package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	DatabaseURL string
	ServerPort  string
}

func Load() *Config {
	_ = godotenv.Load()

	fmt.Println("in config", os.Getenv("SERVER_PORT"))
	return &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		ServerPort:  os.Getenv("SERVER_PORT"),
	}
}
