package config

import (
	"fmt" 
	"os"
)

type Config struct {
	DatabaseURL string
	ServerPort  string
}

func Load () (*Config , error) {
	DatabaseURL := os.Getenv("DATABASE_URL") 
	if DatabaseURL == "" {
		return nil , fmt.Errorf("DATABASE_URL enviroment variable is required")
	}
	ServerPort := os.Getenv("SERVER_PORT")
	if ServerPort == "" {
		ServerPort = ":8080"
	}
	return &Config{
		DatabaseURL: DatabaseURL,
		ServerPort:  ServerPort,
	 }, nil
}