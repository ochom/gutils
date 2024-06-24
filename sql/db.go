package sql

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// internal db initialized by the New function
var db *gorm.DB

// Conn returns the database db
func Conn() *gorm.DB { return db }

// defaultConfig ...
var config = &Config{
	Driver:       Sqlite,
	Url:          "gorm.db",
	LogLevel:     logger.Silent,
	MaxIdleConns: 10,
	MaxOpenConns: 100,
	ConnLifeTime: time.Hour,
}

// New initializes the database db with GORM
func New(configs ...*Config) (err error) {
	for _, cfg := range configs {
		if cfg.Driver != Sqlite {
			config.Driver = cfg.Driver
		}

		if cfg.Url != "" {
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

	return createInstance()
}

func createInstance() error {
	switch config.Driver {
	case Postgres:
		return createPgInstance()
	case MySQL:
		return createMysqlInstance()
	default:
		return createSqliteInstance()
	}
}

func createPgInstance() (err error) {
	db, err = gorm.Open(postgres.Open(config.Url), &gorm.Config{
		Logger: logger.Default.LogMode(config.LogLevel),
	})

	if err != nil {
		return
	}

	return createPool()
}

func createMysqlInstance() (err error) {
	db, err = gorm.Open(mysql.Open(config.Url), &gorm.Config{
		Logger: logger.Default.LogMode(config.LogLevel),
	})

	if err != nil {
		return
	}

	return createPool()
}

func createSqliteInstance() (err error) {
	// - Set WAL mode (not strictly necessary each time because it's persisted in the database, but good for first run)
	// - Set busy timeout, so concurrent writers wait on each other instead of erroring immediately
	// - Enable foreign key checks
	// -  see https://www.golang.dk/articles/go-and-sqlite-in-the-cloud

	url := config.Url + "?_journal=WAL&_timeout=5000&_fk=true"
	db, err = gorm.Open(sqlite.Open(url), &gorm.Config{
		Logger: logger.Default.LogMode(config.LogLevel),
	})

	if err != nil {
		return
	}

	return createPool()
}

func createPool() error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(config.ConnLifeTime)

	return nil
}
