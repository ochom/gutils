package env

import (
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
	"github.com/ochom/gutils/logs"
)

// Get returns env variable or the provided default value when variable not found
func Get(props ...string) string {
	if len(props) == 0 {
		logs.Error("Get Env: props cannot be empty")
		return ""
	}

	key := props[0]
	defaultValue := ""
	if len(props) > 1 {
		defaultValue = props[1]
	}

	value, ok := os.LookupEnv(key)
	if !ok {
		logs.Warn("Get Env: %s not found", key)
		return defaultValue
	}

	return value
}

// Int returns an integer from env variable or the provided default value when variable not found
func Int(key string, defaultValue int) int {
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

// Bool returns a boolean from env variable or the provided default value when variable not found
func Bool(key string, defaultValue bool) bool {
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

// Float returns a float from env variable or the provided default value when variable not found
func Float(key string, defaultValue float64) float64 {
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
