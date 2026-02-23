// Package sqlr provides a GORM-based SQL database abstraction layer with generic CRUD operations.
//
// This package supports PostgreSQL, MySQL, and SQLite databases with automatic driver selection
// based on the connection URL prefix. It provides both a global instance for simple use cases
// and the ability to create multiple connections.
//
// Features:
//   - Generic type-safe CRUD operations
//   - Connection pooling with configurable limits
//   - Flexible transaction support
//   - Scoped queries using GORM's scope pattern
//   - Automatic database migrations
//
// Example usage:
//
//	// Define your model
//	type User struct {
//		ID        uint   `gorm:\"primaryKey\"`
//		Name      string `gorm:\"size:255\"`
//		Email     string `gorm:\"unique\"`
//		CreatedAt time.Time
//	}
//
//	// Initialize the database
//	err := sqlr.Init(&sqlr.Config{
//		Url: \"postgres://user:pass@localhost:5432/mydb?sslmode=disable\",
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Auto-migrate
//	sqlr.Migrate(&User{})
//
//	// Create a user
//	user := &User{Name: \"Alice\", Email: \"alice@example.com\"}
//	sqlr.Create(user)
//
//	// Find by ID
//	user, err := sqlr.FindOneById[User](1)
//
//	// Find with scopes
//	users := sqlr.FindAll[User](func(db *gorm.DB) *gorm.DB {
//		return db.Where(\"name LIKE ?\", \"A%\")
//	})
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

// dbInstance holds the internal database connections.
type dbInstance struct {
	gormDB *gorm.DB
	sqlDB  *sql.DB
}

var instance = &dbInstance{}

// GORM returns the underlying GORM database instance for advanced operations.
//
// Example:
//
//	// Use for complex queries
//	db := sqlr.GORM()
//	db.Joins(\"JOIN orders ON orders.user_id = users.id\").
//		Where(\"orders.total > ?\", 100).
//		Find(&users)
//
//	// Use with raw SQL
//	db.Raw(\"SELECT * FROM users WHERE age > ?\", 18).Scan(&users)
func GORM() *gorm.DB {
	return instance.gormDB
}

// SQL returns the underlying standard library sql.DB instance.
// Useful for operations that require direct database/sql access.
//
// Example:
//
//	db := sqlr.SQL()
//	err := db.Ping()
//	stats := db.Stats()
func SQL() *sql.DB {
	return instance.sqlDB
}

// Init initializes the global database connection.
// Call this once at application startup before using other sqlr functions.
//
// The database driver is automatically selected based on the URL prefix:
//   - postgres:// or postgresql:// -> PostgreSQL
//   - mysql:// -> MySQL
//   - file path or other -> SQLite
//
// Example:
//
//	// PostgreSQL
//	err := sqlr.Init(&sqlr.Config{
//		Url: \"postgres://user:pass@localhost:5432/mydb?sslmode=disable\",
//	})
//
//	// MySQL
//	err := sqlr.Init(&sqlr.Config{
//		Url: \"mysql://user:pass@tcp(localhost:3306)/mydb\",
//	})
//
//	// SQLite
//	err := sqlr.Init(&sqlr.Config{
//		Url: \"./data.db\",
//	})
//
//	// With custom settings
//	err := sqlr.Init(&sqlr.Config{
//		Url:          \"postgres://...\",
//		MaxOpenConns: 50,
//		LogLevel:     logger.Warn,
//	})
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

// New creates and returns a new database connection without affecting the global instance.
// Use this when you need multiple database connections or want to manage connections explicitly.
//
// Example:
//
//	// Create a separate connection
//	gormDB, sqlDB, err := sqlr.New(&sqlr.Config{
//		Url: \"postgres://user:pass@localhost:5432/analytics\",
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer sqlDB.Close()
//
//	// Use the connection
//	var users []User
//	gormDB.Find(&users)
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
		return createPool(postgres.New(postgres.Config{
			DSN:                  config.Url,
			PreferSimpleProtocol: config.PreparedStatements,
		}), config)
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
		PrepareStmt:            config.PreparedStatements,
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
