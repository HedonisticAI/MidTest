package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_DSN   string
	HttpPort string
	ADMToken string
}

func NewConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
		return nil
	}
	DSN, exist := os.LookupEnv("DB_DSN")
	if !exist {
		return nil
	}
	HttpPort, exist := os.LookupEnv("HTTP_PORT")
	if !exist {
		return nil
	}
	ADMToken, exist := os.LookupEnv("ADM_TOKEN")
	if !exist {
		return nil
	}
	return &Config{DB_DSN: DSN, HttpPort: HttpPort, ADMToken: ADMToken}
}
