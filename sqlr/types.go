package sqlr

import (
	"time"

	"gorm.io/gorm/logger"
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
	Url                    string
	LogLevel               logger.LogLevel
	MaxOpenConns           int
	MaxIdleConns           int
	MaxConnIdleTime        time.Duration
	MaxConnLifeTime        time.Duration
	SkipDefaultTransaction bool
}

// defaultConfig ...
var defaultConfig = Config{
	Url:                    "gorm.db",
	LogLevel:               logger.Silent,
	MaxIdleConns:           10,
	MaxOpenConns:           100,
	MaxConnLifeTime:        time.Hour,
	MaxConnIdleTime:        time.Minute,
	SkipDefaultTransaction: true,
}
