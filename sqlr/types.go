package sqlr

import (
	"time"

	"gorm.io/gorm/logger"
)

// Database configuration
type Config struct {
	Url                       string
	LogLevel                  logger.LogLevel
	IgnoreRecordNotFoundError bool
	MaxOpenConns              int
	MaxIdleConns              int
	MaxConnIdleTime           time.Duration
	MaxConnLifeTime           time.Duration
	SkipDefaultTransaction    bool
}

// defaultConfig ...
var defaultConfig = Config{
	Url:                    "gorm.db",
	LogLevel:               logger.Info,
	MaxIdleConns:           10,
	MaxOpenConns:           100,
	MaxConnLifeTime:        time.Hour,
	MaxConnIdleTime:        time.Minute,
	SkipDefaultTransaction: true,
}
