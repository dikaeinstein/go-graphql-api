package config

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config is the configuration of the server.
type Config struct {
	AppEnv           string
	DBName           string
	DBUser           string
	DBConnectTimeout int
	LogLevel         int
	Port             int
}

// New creates an instance of config.
func New() Config {
	err := godotenv.Load()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Println(".env not found. Loading env variables from process env")
		} else {
			log.Fatal(err)
		}
	}

	return Config{
		AppEnv:           getEnv("APP_ENV", "development"),
		DBName:           getEnv("DB_NAME", ""),
		DBUser:           getEnv("DB_USER", ""),
		DBConnectTimeout: getEnvAsInt("DB_CONNECT_TIMEOUT", 0),
		Port:             getEnvAsInt("PORT", 10000),
		LogLevel:         getEnvAsInt("LOG_LEVEL", 0),
	}
}

// Simple helper function to read an environment or return a default value.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Helper to read an environment variable into a bool or return a default value.
func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}
	return defaultVal
}

// Simple helper function to read an environment variable into an integer
// or return a default value.
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}
