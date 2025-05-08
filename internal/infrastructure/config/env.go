// internal/infrastructure/config/env.go
package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

func init() {
	// Get the path to the current file
	_, filename, _, _ := runtime.Caller(0)
	// Navigate to project root (adjust the number of "../" as needed)
	projectRoot := filepath.Join(filepath.Dir(filename), "../../..")

	// Load .env file from project root
	envPath := filepath.Join(projectRoot, ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("Warning: Error loading .env file from %s: %v\n", envPath, err)
		// Continue execution - in Docker the env vars will be set directly
	} else {
		log.Printf("Loaded .env file from: %s\n", envPath)
	}
}

// Helper function to get an environment variable or a default value
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
