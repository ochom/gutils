package helpers

import (
	"os"
	"strconv"

	"github.com/ochom/gutils/logs"
)

// GetEnv returns env variable or the provided default value when variable not found
func GetEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		logs.Error("Environment variable %s not found", key)
	}

	return value
}

// GetEnv returns env variable or the provided default value when variable not found
func GetEnvDefault(key string, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	return value
}

// GetEnvInt returns an integer from env variable or the provided default value when variable not found
func GetEnvInt(key string, defaultValue int) int {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	val, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return val
}

// GetEnvBool returns a boolean from env variable or the provided default value when variable not found
func GetEnvBool(key string, defaultValue bool) bool {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	val, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}

	return val
}

// GetEnvFloat returns a float from env variable or the provided default value when variable not found
func GetEnvFloat(key string, defaultValue float64) float64 {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	val, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return defaultValue
	}

	return val
}
