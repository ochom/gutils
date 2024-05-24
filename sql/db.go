package sql

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// conn initialized in Init() and used in all other functions
var conn *gorm.DB

// Conn returns the database connection
func Conn() *gorm.DB { return conn }

// defaultConfig ...
func defaultConfig() *Config {
	return &Config{
		DatabaseType: Sqlite,
		Url:          "gorm.db",
		LogLevel:     logger.Silent,
		MaxIdleConns: 10,
		MaxOpenConns: 100,
		ConnLifeTime: time.Hour,
	}
}

// Init initializes the database connection with GORM
func Init(configs ...Config) error {
	var err error

	config := defaultConfig()
	for _, cfg := range configs {
		if cfg.DatabaseType != "" {
			config.Url = cfg.Url
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

		if cfg.ConnLifeTime != 0 {
			config.ConnLifeTime = cfg.ConnLifeTime
		}
	}

	switch config.DatabaseType {
	case Postgres:
		conn, err = gorm.Open(postgres.Open(config.Url), &gorm.Config{
			Logger: logger.Default.LogMode(config.LogLevel),
		})
	case MySQL:
		conn, err = gorm.Open(mysql.Open(config.Url), &gorm.Config{
			Logger: logger.Default.LogMode(config.LogLevel),
		})
	default:
		// - Set WAL mode (not strictly necessary each time because it's persisted in the database, but good for first run)
		// - Set busy timeout, so concurrent writers wait on each other instead of erroring immediately
		// - Enable foreign key checks
		// -  see https://www.golang.dk/articles/go-and-sqlite-in-the-cloud
		url := config.Url + "?_journal=WAL&_timeout=5000&_fk=true"
		conn, err = gorm.Open(sqlite.Open(url), &gorm.Config{
			Logger: logger.Default.LogMode(config.LogLevel),
		})
	}

	if err != nil {
		return err
	}

	sqlDB, err := conn.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(config.ConnLifeTime)

	return nil
}
