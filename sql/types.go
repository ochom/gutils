package sql

import (
	"time"

	"gorm.io/gorm/logger"
)

type Driver int

// Drivers ...
const (
	Sqlite Driver = iota
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
	Driver       Driver
	Url          string
	LogLevel     logger.LogLevel
	MaxIdleConns int
	MaxOpenConns int
	ConnLifeTime time.Duration
}
