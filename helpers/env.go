package helpers

import (
	"os"
	"strconv"

	"github.com/ochom/gutils/logger"
)

var log = logger.NewLogger()

// GetEnv ...
func GetEnv(key string, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Warn("Environment variable `%s` not found, returning default value `%s`\n", key, defaultValue)
		return defaultValue
	}

	return value
}

// GetEnvInt ...
func GetEnvInt(key string, defaultValue int) int {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Warn("Environment variable `%s` not found, returning default value `%v`\n", key, defaultValue)
		return defaultValue
	}

	val, err := strconv.Atoi(value)
	if err != nil {
		log.Warn("Environment variable `%s` error: %s, returning default value `%v`\n", key, err.Error(), defaultValue)
		return defaultValue
	}

	return val
}

// GetEnvBool ...
func GetEnvBool(key string, defaultValue bool) bool {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Warn("Environment variable `%s` not found, returning default value `%v`\n", key, defaultValue)
		return defaultValue
	}

	val, err := strconv.ParseBool(value)
	if err != nil {
		log.Warn("Environment variable `%s` error: %s, returning default value `%v`\n", key, err.Error(), defaultValue)
		return defaultValue
	}

	return val
}

// GetEnvFloat ...
func GetEnvFloat(key string, defaultValue float64) float64 {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Warn("Environment variable `%s` not found, returning default value `%v`\n", key, defaultValue)
		return defaultValue
	}

	val, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Warn("Environment variable `%s` error: %s, returning default value `%v`\n", key, err.Error(), defaultValue)
		return defaultValue
	}

	return val
}
