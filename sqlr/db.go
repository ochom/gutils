package sqlr

import (
	"database/sql"
	"log"
	"os"
	"strings"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// internal db initialized by the New function
type dbInstance struct {
	gormDB *gorm.DB
	sqlDB  *sql.DB
}

var instance = &dbInstance{}

// GORM returns the gorm db instance
func GORM() *gorm.DB {
	return instance.gormDB
}

// SQL returns the sql db instance
func SQL() *sql.DB {
	return instance.sqlDB
}

// Init initializes the database db with GORM
func Init(configs ...*Config) (err error) {
	config := parseConfig(configs...)
	gormDB, sqlDB, err := createInstance(config)
	if err != nil {
		return err
	}

	instance.gormDB = gormDB
	instance.sqlDB = sqlDB
	return nil
}

// New Create connection create and returns a new connection
func New(cfg ...*Config) (*gorm.DB, *sql.DB, error) {
	config := parseConfig(cfg...)
	gormDB, sqlDB, err := createInstance(config)
	if err != nil {
		return nil, nil, err
	}

	return gormDB, sqlDB, nil
}

func parseConfig(configs ...*Config) *Config {
	config := &defaultConfig
	for _, cfg := range configs {
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

func createInstance(config *Config) (gormDB *gorm.DB, sqlDB *sql.DB, err error) {
	if strings.HasPrefix(config.Url, "postgres") {
		return createPool(postgres.Open(config.Url), config)
	}

	if strings.HasPrefix(config.Url, "mysql") {
		return createPool(mysql.Open(config.Url), config)
	}

	url := config.Url + "?_journal=WAL&_timeout=5000&_fk=true"
	return createPool(sqlite.Open(url), config)
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
		PrepareStmt:            true,
	}
}

func createPool(conn gorm.Dialector, config *Config) (gormDB *gorm.DB, sqlDB *sql.DB, err error) {
	gormDB, err = gorm.Open(conn, getGormConfig(config))
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
