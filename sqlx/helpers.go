package sqlx

import (
	"github.com/ochom/gutils/logs"
	"gorm.io/gorm"
)

// Create ...
func Create[T any](data *T) error {
	return db.Create(data).Error
}

// Update ...
func Update[T any](data *T) error {
	return db.Save(data).Error
}

// Delete ...
func Delete[T any](scopes ...func(*gorm.DB) *gorm.DB) error {
	return db.Scopes(scopes...).Delete(new(T)).Error
}

// DeleteById ...
func DeleteById[T any](id any, scopes ...func(*gorm.DB) *gorm.DB) error {
	return Delete[T](append(scopes, func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	})...)
}

// FindOne ...
func FindOne[T any](scopes ...func(*gorm.DB) *gorm.DB) (*T, error) {
	var data T
	if err := db.Scopes(scopes...).First(&data).Error; err != nil {
		return nil, err
	}

	return &data, nil
}

// FindOneById ...
func FindOneById[T any](id any, scopes ...func(*gorm.DB) *gorm.DB) (*T, error) {
	return FindOne[T](append(scopes, func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	})...)
}

// FindAll ...
func FindAll[T any](scopes ...func(*gorm.DB) *gorm.DB) []*T {
	data := []*T{}
	if err := db.Scopes(scopes...).Find(&data).Error; err != nil {
		logs.Info("FindAll: %s", err.Error())
		return []*T{}
	}

	return data
}

// FindWithLimit ...
func FindWithLimit[T any](page, limit int, scopes ...func(*gorm.DB) *gorm.DB) []*T {
	data := []*T{}
	if err := db.Scopes(scopes...).Offset((page - 1) * limit).Limit(limit).Find(&data).Error; err != nil {
		logs.Info("FindWithLimit: %s", err.Error())
		return []*T{}
	}
	return data
}

// Count ...
func Count[T any](scopes ...func(*gorm.DB) *gorm.DB) int {
	var count int64
	var model T
	if err := db.Model(&model).Scopes(scopes...).Count(&count).Error; err != nil {
		logs.Info("Count: %s", err.Error())
		return 0
	}

	return int(count)
}
