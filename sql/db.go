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
		DSN:          "gorm.db",
		LogLevel:     logger.Silent,
		MaxIdleConns: 10,
		MaxOpenConns: 100,
		ConnLifeTime: time.Hour,
	}
}

// Init initializes the database connection with GORM
func Init(configs ...ConfigFunc) error {
	var err error

	config := defaultConfig()
	for _, fn := range configs {
		fn(config)
	}

	switch config.DatabaseType {
	case Postgres:
		conn, err = gorm.Open(postgres.Open(config.DSN), &gorm.Config{
			Logger: logger.Default.LogMode(config.LogLevel),
		})
	case MySQL:
		conn, err = gorm.Open(mysql.Open(config.DSN), &gorm.Config{
			Logger: logger.Default.LogMode(config.LogLevel),
		})
	default:
		conn, err = gorm.Open(sqlite.Open(config.DSN), &gorm.Config{
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
