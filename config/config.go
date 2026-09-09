package config

import (
	"os"
)

type Config struct {
	DB_DSN   string
	HttpPort string
	ADMToken string
}

func NewConfig() *Config {
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
