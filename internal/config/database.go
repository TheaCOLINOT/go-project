package config

import (
	"fmt"
	"os"
	"strconv"
)

// Database contains the settings required to connect to PostgreSQL.
type Database struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	SSLMode  string
}

// HTTPAddress returns the address used by the HTTP server.
func HTTPAddress() string {
	return envOrDefault("HTTP_ADDR", ":8080")
}

// LoadDatabase reads PostgreSQL settings from the environment.
func LoadDatabase() (Database, error) {
	port, err := strconv.Atoi(envOrDefault("POSTGRES_PORT", "5432"))
	if err != nil || port < 1 || port > 65535 {
		return Database{}, fmt.Errorf("POSTGRES_PORT must be a valid TCP port")
	}

	return Database{
		Host:     envOrDefault("POSTGRES_HOST", "localhost"),
		Port:     port,
		Name:     envOrDefault("POSTGRES_DB", "ecommerce"),
		User:     envOrDefault("POSTGRES_USER", "ecommerce"),
		Password: envOrDefault("POSTGRES_PASSWORD", "dev_password_change_me"),
		SSLMode:  envOrDefault("POSTGRES_SSLMODE", "disable"),
	}, nil
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
