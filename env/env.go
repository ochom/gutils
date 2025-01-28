package env

import (
	"fmt"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

// Get returns env variable or the provided default value when variable not found
func Get(name string, defaults ...string) string {
	defaultValue := ""
	if len(defaults) > 1 {
		defaultValue = defaults[0]
	}

	value, ok := os.LookupEnv(name)
	if !ok {
		return defaultValue
	}

	return value
}

// Int returns an integer from env variable or the provided default value when variable not found
func Int(name string, defaultValue int) int {
	value, ok := os.LookupEnv(name)
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
func Bool(name string, defaultValue bool) bool {
	value, ok := os.LookupEnv(name)
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
func Float(name string, defaultValue float64) float64 {
	value, ok := os.LookupEnv(name)
	if !ok {
		return defaultValue
	}

	val, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return defaultValue
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
