package config

import "os"

// Config holds all application configuration loaded from environment variables.
type Config struct {
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	ServerPort string
}

// Load reads configuration from environment variables, applying defaults where needed.
func Load() *Config {
	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBName:     getEnv("DB_NAME", "golang_demo"),
		DBUser:     getEnv("DB_USER", "golang_user"),
		DBPassword: getEnv("DB_PASSWORD", "golang_pass"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
