package env

import (
	"fmt"
	"os"
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

// Get returns env variable or the provided default value when variable not found
func Get(name string, defaults ...string) string {
	value, ok := os.LookupEnv(name)
	if !ok {
		return defaultValue(defaults...)
	}

	return value
}

// Int returns an integer from env variable or the provided default value when variable not found
func Int(name string, defaults ...int) int {
	value, ok := os.LookupEnv(name)
	if !ok {
		return defaultValue(defaults...)
	}

	val, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue(defaults...)
	}

	return val
}

// Bool returns a boolean from env variable or the provided default value when variable not found
func Bool(name string, defaults ...bool) bool {
	value, ok := os.LookupEnv(name)
	if !ok {
		return defaultValue(defaults...)
	}

	val, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue(defaults...)
	}

	return val
}

// Float returns a float from env variable or the provided default value when variable not found
func Float(name string, defaults ...float64) float64 {
	value, ok := os.LookupEnv(name)
	if !ok {
		return defaultValue(defaults...)
	}

	val, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return defaultValue(defaults...)
	}

	return val
}

// MustGet returns env variable or panics when variable not found
func MustGet(name string) string {
	value, ok := os.LookupEnv(name)
	if !ok {
		panic(fmt.Errorf("MustGet Env: %s not found", name))
	}

	return value
}

// MustInt returns an integer from env variable or panics when variable not found
func MustInt(name string) int {
	value, ok := os.LookupEnv(name)
	if !ok {
		panic(fmt.Errorf("MustInt Env: %s not found", name))
	}

	val, err := strconv.Atoi(value)
	if err != nil {
		panic(fmt.Errorf("MustInt Env: %s not an integer", name))
	}

	return val
}

// MustBool returns a boolean from env variable or panics when variable not found
func MustBool(name string) bool {
	value, ok := os.LookupEnv(name)
	if !ok {
		panic(fmt.Errorf("MustBool Env: %s not found", name))
	}

	val, err := strconv.ParseBool(value)
	if err != nil {
		panic(fmt.Errorf("MustBool Env: %s not a boolean", name))
	}

	return val
}

// MustFloat returns a float from env variable or panics when variable not found
func MustFloat(name string) float64 {
	value, ok := os.LookupEnv(name)
	if !ok {
		panic(fmt.Errorf("MustFloat Env: %s not found", name))
	}

	val, err := strconv.ParseFloat(value, 64)
	if err != nil {
		panic(fmt.Errorf("MustFloat Env: %s not a float", name))
	}

	return val
}
