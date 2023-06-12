package sql

import "gorm.io/gorm/logger"

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
