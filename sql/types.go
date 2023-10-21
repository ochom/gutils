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
	DSN          string
	LogLevel     logger.LogLevel
	MaxIdleConns int
	MaxOpenConns int
	ConnLifeTime time.Duration
}

// ConfigFunc ...
type ConfigFunc func(*Config)

// WithDatabaseType ...
func WithDatabaseType(dbType Platform) ConfigFunc {
	return func(c *Config) {
		c.DatabaseType = dbType
	}
}

// WithDSN ...
func WithDSN(dsn string) ConfigFunc {
	return func(c *Config) {
		c.DSN = dsn
	}
}

// WithLogLevel ...
func WithLogLevel(logLevel logger.LogLevel) ConfigFunc {
	return func(c *Config) {
		c.LogLevel = logLevel
	}
}

// WithMaxIdleConns ...
func WithMaxIdleConns(maxIdleConns int) ConfigFunc {
	return func(c *Config) {
		c.MaxIdleConns = maxIdleConns
	}
}

// WithMaxOpenConns ...
func WithMaxOpenConns(maxOpenConns int) ConfigFunc {
	return func(c *Config) {
		c.MaxOpenConns = maxOpenConns
	}
}

// WithConnLifeTime ...
func WithConnLifeTime(connLifeTime time.Duration) ConfigFunc {
	return func(c *Config) {
		c.ConnLifeTime = connLifeTime
	}
}
