package helpers

import (
	"os"
	"strconv"
)

// GetEnv ...
func GetEnv(key string, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	return value
}

// GetEnvInt ...
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

// GetEnvBool ...
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

// GetEnvFloat ...
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
