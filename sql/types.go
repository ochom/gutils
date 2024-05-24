package sql

import (
	"time"

	"gorm.io/gorm/logger"
)

// Platform ...
type Platform string

// Platforms ...
const (
	Postgres Platform = "postgres"
	MySQL    Platform = "mysql"
	Sqlite   Platform = "sqlite"
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
	DatabaseType Platform
	Url          string
	LogLevel     logger.LogLevel
	MaxIdleConns int
	MaxOpenConns int
	ConnLifeTime time.Duration
}
