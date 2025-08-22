package utils

import (
	"os"
	"strconv"
)

// GetEnv gets an environment variable with a default value
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvAsInt gets an environment variable as integer with a default value
func GetEnvAsInt(key string, defaultValue int) int {
	valueStr := GetEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// GetEnvAsBool gets an environment variable as boolean with a default value
func GetEnvAsBool(key string, defaultValue bool) bool {
	valueStr := GetEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}
