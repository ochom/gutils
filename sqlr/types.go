package sqlr

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config holds the database connection configuration options.
//
// Example:
//
//	config := &sqlr.Config{
//		Url:          "postgres://user:pass@localhost:5432/mydb",
//		LogLevel:     logger.Warn,
//		MaxOpenConns: 50,
//		MaxIdleConns: 10,
//	}
//	err := sqlr.Init(config)
type Config struct {
	// Conn is the GORM dialector for the database connection
	Conn gorm.Dialector

	// LogLevel controls GORM's logging verbosity (Silent, Error, Warn, Info)
	LogLevel logger.LogLevel

	// IgnoreRecordNotFoundError suppresses "record not found" errors in logs
	IgnoreRecordNotFoundError bool

	// MaxOpenConns sets the maximum number of open connections to the database
	MaxOpenConns int

	// MaxIdleConns sets the maximum number of idle connections in the pool
	MaxIdleConns int

	// MaxConnIdleTime sets the maximum time a connection can be idle before closing
	MaxConnIdleTime time.Duration

	// MaxConnLifeTime sets the maximum lifetime of a connection
	MaxConnLifeTime time.Duration

	// SkipDefaultTransaction disables wrapping each operation in a transaction
	SkipDefaultTransaction bool

	// PreparedStatements enables prepared statement caching
	PreparedStatements bool
}

// defaultConfig provides sensible defaults for database connections.
var defaultConfig = Config{
	LogLevel:               logger.Info,
	MaxIdleConns:           10,
	MaxOpenConns:           100,
	MaxConnLifeTime:        time.Hour,
	MaxConnIdleTime:        time.Minute,
	SkipDefaultTransaction: true,
}
