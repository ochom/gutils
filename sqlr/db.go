package sqlr

import (
	"database/sql"
	"log"
	"os"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// dbInstance holds the internal database connections.
type dbInstance struct {
	gormDB *gorm.DB
	sqlDB  *sql.DB
}

var instance = &dbInstance{}

// GORM returns the underlying GORM database instance for advanced operations.
func GORM() *gorm.DB {
	return instance.gormDB
}

// SQL returns the underlying standard library sql.DB instance.
func SQL() *sql.DB {
	return instance.sqlDB
}

// Init initializes the global database connection.
func Init(configs ...*Config) (err error) {
	config := parseConfig(configs...)
	gormDB, sqlDB, err := createPool(config)
	if err != nil {
		return err
	}

	instance.gormDB = gormDB
	instance.sqlDB = sqlDB
	return nil
}

// New creates and returns a new database connection without affecting the global instance.
func New(dialer gorm.Dialector, cfg ...*Config) (*gorm.DB, *sql.DB, error) {
	config := parseConfig(cfg...)
	gormDB, sqlDB, err := createPool(config)
	if err != nil {
		return nil, nil, err
	}

	return gormDB, sqlDB, nil
}

func parseConfig(configs ...*Config) *Config {
	config := &defaultConfig
	for _, cfg := range configs {
		if cfg.Conn != nil {
			config.Conn = cfg.Conn
		}

		if cfg.LogLevel != 0 {
			config.LogLevel = cfg.LogLevel
		}

		if cfg.MaxIdleConns != 0 {
			config.MaxIdleConns = cfg.MaxIdleConns
		}

		if cfg.MaxOpenConns != 0 {
			config.MaxOpenConns = cfg.MaxOpenConns
		}

		if cfg.MaxConnIdleTime != 0 {
			config.MaxConnIdleTime = cfg.MaxConnIdleTime
		}

		if cfg.MaxConnLifeTime != 0 {
			config.MaxConnLifeTime = cfg.MaxConnLifeTime
		}

		if cfg.SkipDefaultTransaction {
			config.SkipDefaultTransaction = cfg.SkipDefaultTransaction
		}
	}

	return config
}

func getGormConfig(config *Config) *gorm.Config {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),

		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  config.LogLevel,
			IgnoreRecordNotFoundError: config.IgnoreRecordNotFoundError,
			Colorful:                  true,
		})

	return &gorm.Config{
		Logger:                 newLogger,
		SkipDefaultTransaction: config.SkipDefaultTransaction,
		PrepareStmt:            config.PreparedStatements,
	}
}

func createPool(config *Config) (gormDB *gorm.DB, sqlDB *sql.DB, err error) {
	gormDB, err = gorm.Open(config.Conn, getGormConfig(config))
	if err != nil {
		return
	}

	sqlDB, err = gormDB.DB()
	if err != nil {
		return
	}

	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(config.MaxConnLifeTime)
	sqlDB.SetConnMaxIdleTime(config.MaxConnIdleTime)

	return gormDB, sqlDB, nil
}
