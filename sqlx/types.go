package sqlx

import (
	"time"

	"gorm.io/gorm/logger"
)

// Driver supported by this package
type Driver string

// String returns the string representation of the driver
func (d Driver) String() string {
	return string(d)
}

// Drivers ...
const (
	Sqlite   Driver = "sqlite"
	Postgres Driver = "postgres"
	MySQL    Driver = "mysql"
)

// LogLevels ...
const (
	Silent logger.LogLevel = iota
	Error
	Warn
	Info
)

// Database configuration
type Config struct {
	Driver                 string
	Url                    string
	LogLevel               logger.LogLevel
	MaxOpenConns           int
	MaxIdleConns           int
	MaxConnIdleTime        time.Duration
	MaxConnLifeTime        time.Duration
	SkipDefaultTransaction bool
}
