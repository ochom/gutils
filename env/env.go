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

type Env interface {
	string | bool | int | float64
}

// Get returns env variable or the provided default value when variable not found
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

// MustGet returns env variable or panics when variable not found
func MustGet[T Env](name string) T {
	envVal := Get[T](name)
	if reflect.DeepEqual(envVal, defaultValue[T]()) {
		panic(fmt.Errorf("MustGet Env: %s not found", name))
	}

	return envVal
}
