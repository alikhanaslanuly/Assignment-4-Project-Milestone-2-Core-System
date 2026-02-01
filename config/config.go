package config

import (
	"os"
)

type Config struct {
	DBConnStr string
	Port      string
}

func LoadConfig() *Config {
	return &Config{
		DBConnStr: getEnv("DB_CONN_STR", "host=localhost port=5432 user=postgres password=09022107 dbname=eventify sslmode=disable"),
		Port:      getEnv("PORT", ":8080"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
