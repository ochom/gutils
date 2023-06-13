package sql

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var conn *gorm.DB

// Init initializes the database connection with GORM
func Init(dbType Platform, dsn string, logLevel logger.LogLevel) error {
	var err error

	switch dbType {
	case Postgres:
		conn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(0),
		})
	case MySQL:
		conn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logLevel)})
	default:
		conn, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logLevel),
		})
	}

	return err
}

// Conn returns the database connection
func Conn() *gorm.DB {
	return conn
}

// Create ...
func Create[T any](data *T) error {
	return conn.Create(data).Error
}

// Update ...
func Update[T any](data *T) error {
	return conn.Save(data).Error
}

// Delete ...
func Delete[T any](query *T) error {
	return conn.Delete(query).Error
}

// FindOne ...
func FindOne[T any](query *T) (*T, error) {
	var data T
	err := conn.First(&data, query).Error
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// FindAll ...
func FindAll[T any](query *T) ([]*T, error) {
	var data []*T
	err := conn.Find(&data, query).Error
	return data, err
}

// FindWithLimit ...
func FindWithLimit[T any](query *T, limit int) ([]*T, error) {
	var data []*T
	err := conn.Limit(limit).Find(&data, query).Error
	return data, err
}

// Count ...
func Count[T any](query *T) (int64, error) {
	var count int64
	var model T
	err := conn.Model(&model).Where(query).Count(&count).Error
	return count, err
}
