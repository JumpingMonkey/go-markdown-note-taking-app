package config

import (
	"os"
)

// Config holds the application configuration
type Config struct {
	Port     string
	NotesDir string
	LogLevel string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Port:     getEnv("PORT", "8080"),
		NotesDir: getEnv("NOTES_DIR", "./notes"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
