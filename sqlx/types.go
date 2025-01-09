package sqlx

import (
	"time"

	"gorm.io/gorm/logger"
)

// Drivers ...
const (
	Sqlite = iota
	Postgres
	MySQL
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
	Driver                 int
	Url                    string
	LogLevel               logger.LogLevel
	MaxIdleConns           int
	MaxOpenConns           int
	ConnLifeTime           time.Duration
	SkipDefaultTransaction bool
}
