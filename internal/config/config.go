package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	DBHost         string
	DBPort         string
	DBName         string
	DBUser         string
	DBPassword     string
	DBSSLMode      string
	RabbitMQURL    string
	KeycloakURL    string
	KeycloakRealm  string
	KeycloakClient string
}

func Load() *Config {
	_ = godotenv.Load() // Load from .env if exists

	return &Config{
		Port:           getEnv("PORT", "8080"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBName:         getEnv("DB_NAME", ""),
		DBUser:         getEnv("DB_USER", ""),
		DBPassword:     getEnv("DB_PASSWORD", ""),
		DBSSLMode:      getEnv("DB_SSLMODE", "disable"),
		RabbitMQURL:    getEnv("RABBITMQ_URL", ""),
		KeycloakURL:    getEnv("KEYCLOAK_URL", ""),
		KeycloakRealm:  getEnv("KEYCLOAK_REALM", ""),
		KeycloakClient: getEnv("KEYCLOAK_CLIENT", ""),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	if fallback == "" {
		log.Fatalf("Required env var %s not set", key)
	}
	return fallback
}
