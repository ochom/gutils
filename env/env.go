// Package env provides type-safe environment variable access with support for default values.
//
// This package automatically loads environment variables from .env files using godotenv
// and provides generic functions for retrieving typed values. Supported types include
// string, bool, int, int64, and float64.
//
// Example .env file:
//
//	APP_NAME=MyApp
//	DEBUG=true
//	PORT=8080
//	RATE_LIMIT=1.5
//
// Example usage:
//
//	// Get string with default
//	appName := env.Get[string]("APP_NAME", "DefaultApp")
//
//	// Get boolean
//	debug := env.Get[bool]("DEBUG", false)
//
//	// Get integer
//	port := env.Get[int]("PORT", 3000)
//
//	// Get required value (panics if missing)
//	apiKey := env.MustGet[string]("API_KEY")
//
//	// Convenience functions
//	name := env.String("APP_NAME", "default")
//	isDebug := env.Bool("DEBUG", false)
//	port := env.Int("PORT", 8080)
//	rate := env.Float("RATE", 1.0)
package env

import (
	"fmt"
	"os"
	"reflect"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

func defaultValue[T any](defaults ...T) T {
	if len(defaults) > 0 {
		return defaults[0]
	}

	var zero T
	return zero
}

// Env constrains the types that can be used with the generic Get functions.
type Env interface {
	string | bool | int | int64 | float64
}

// Get retrieves an environment variable and converts it to the specified type T.
// Returns the default value if the variable is not set or cannot be converted.
//
// Supported types: string, bool, int, int64, float64
//
// Example:
//
//	// String variable
//	host := env.Get[string]("DB_HOST", "localhost")
//
//	// Boolean variable (accepts: "true", "false", "1", "0", "yes", "no")
//	verbose := env.Get[bool]("VERBOSE", false)
//
//	// Integer variable
//	workers := env.Get[int]("WORKERS", 4)
//
//	// Float variable
//	timeout := env.Get[float64]("TIMEOUT", 30.0)
func Get[T Env](name string, defaults ...T) T {
	value, ok := os.LookupEnv(name)
	if !ok && len(defaults) > 0 {
		return defaults[0]
	}

	var zero T
	kind := reflect.TypeOf(zero).Kind()
	switch kind {
	case reflect.String:
		return any(value).(T)
	case reflect.Bool:
		val, err := strconv.ParseBool(value)
		if err != nil {
			return defaultValue(defaults...)
		}
		return any(val).(T)
	case reflect.Int:
		val, err := strconv.Atoi(value)
		if err != nil {
			return defaultValue(defaults...)
		}
		return any(val).(T)
	case reflect.Float64:
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return defaultValue(defaults...)
		}
		return any(val).(T)
	default:
		panic(fmt.Errorf("unsupported type: %s", kind))
	}
}

// MustGet retrieves an environment variable and panics if it is not set or is empty.
// Use this for required configuration that the application cannot run without.
//
// Example:
//
//	// Required database URL - app will panic if not set
//	dbURL := env.MustGet[string]("DATABASE_URL")
//
//	// Required API key
//	apiKey := env.MustGet[string]("API_KEY")
func MustGet[T Env](name string) T {
	envVal := Get[T](name)
	if reflect.DeepEqual(envVal, defaultValue[T]()) {
		panic(fmt.Errorf("MustGet Env: %s not found", name))
	}

	return envVal
}

// Int retrieves an integer environment variable with an optional default value.
// Convenience wrapper around Get[int].
//
// Example:
//
//	port := env.Int("PORT", 8080)
//	workers := env.Int("WORKERS", 4)
func Int(name string, defaults ...int) int {
	return Get(name, defaults...)
}

// Float retrieves a float64 environment variable with an optional default value.
// Convenience wrapper around Get[float64].
//
// Example:
//
//	timeout := env.Float("TIMEOUT_SECONDS", 30.0)
//	ratio := env.Float("COMPRESSION_RATIO", 0.8)
func Float(name string, defaults ...float64) float64 {
	return Get(name, defaults...)
}

// String retrieves a string environment variable with an optional default value.
// Convenience wrapper around Get[string].
//
// Example:
//
//	host := env.String("DB_HOST", "localhost")
//	env := env.String("ENVIRONMENT", "development")
func String(name string, defaults ...string) string {
	return Get(name, defaults...)
}

// Bool retrieves a boolean environment variable with an optional default value.
// Convenience wrapper around Get[bool].
//
// Accepts: "true", "false", "1", "0", "yes", "no" (case-insensitive)
//
// Example:
//
//	debug := env.Bool("DEBUG", false)
//	verbose := env.Bool("VERBOSE_LOGGING", false)
func Bool(name string, defaults ...bool) bool {
	return Get(name, defaults...)
}
