package config

import "os"

type Config struct {
	HTTPPort    string
	DatabaseURL string
	RedisURL    string
	JWTSecret   string
	JWTIssuer   string
}

func Load() (Config, error) {
	return Config{
		HTTPPort:    getEnv("HTTP_PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/strutfy?sslmode=disable"),
		RedisURL:    getEnv("REDIS_URL", "localhost:6379"),
		JWTSecret:   getEnv("JWT_SECRET", "dev-secret"),
		JWTIssuer:   getEnv("JWT_ISSUER", "strutfy-api"),
	}, nil
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}