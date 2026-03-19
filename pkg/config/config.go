package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from a file.
//
// Priority (highest to lowest):
// 1) ENV_FILE environment variable
// 2) SERVICE_ENV_FILE environment variable
// 3) .env in current working directory
// 4) .env in the executable directory
func LoadEnv() {
	// Allow overriding the env file path explicitly.
	if envPath := os.Getenv("ENV_FILE"); envPath != "" {
		if err := godotenv.Overload(envPath); err != nil {
			log.Printf("Failed to load env file %q: %v", envPath, err)
		}
		return
	}

	if envPath := os.Getenv("SERVICE_ENV_FILE"); envPath != "" {
		if err := godotenv.Overload(envPath); err != nil {
			log.Printf("Failed to load env file %q: %v", envPath, err)
		}
		return
	}

	// Load global config (e.g., dist/.env) first, then overlay service-specific config
	_ = godotenv.Load()

	// Fallback: load a .env located next to the executable (service directory) and allow it to override.
	if exe, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exe)
		_ = godotenv.Overload(filepath.Join(exeDir, ".env"))
	}
}

func GetEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
