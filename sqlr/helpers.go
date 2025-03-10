package sqlr

import (
	"github.com/ochom/gutils/logs"
	"gorm.io/gorm"
)

// Create ...
func Create[T any](data *T) error {
	return instance.gormDB.Create(data).Error
}

// Update ...
func Update[T any](data *T) error {
	return instance.gormDB.Save(data).Error
}

// Delete ...
func Delete[T any](scopes ...func(db *gorm.DB) *gorm.DB) error {
	return instance.gormDB.Scopes(scopes...).Delete(new(T)).Error
}

// DeleteById ...
func DeleteById[T any](id any, scopes ...func(db *gorm.DB) *gorm.DB) error {
	return Delete[T](append(scopes, func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	})...)
}

// FindOne ...
func FindOne[T any](scopes ...func(db *gorm.DB) *gorm.DB) (*T, error) {
	var data T
	if err := instance.gormDB.Scopes(scopes...).First(&data).Error; err != nil {
		return nil, err
	}

	return &data, nil
}

// FindOneById ...
func FindOneById[T any](id any, scopes ...func(db *gorm.DB) *gorm.DB) (*T, error) {
	return FindOne[T](append(scopes, func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	})...)
}

// FindAll ...
func FindAll[T any](scopes ...func(db *gorm.DB) *gorm.DB) []*T {
	data := []*T{}
	if err := instance.gormDB.Scopes(scopes...).Find(&data).Error; err != nil {
		logs.Info("FindAll: %s", err.Error())
		return []*T{}
	}

	return data
}

// FindWithLimit ...
func FindWithLimit[T any](page, limit int, scopes ...func(db *gorm.DB) *gorm.DB) []*T {
	data := []*T{}
	if err := instance.gormDB.Scopes(scopes...).Offset((page - 1) * limit).Limit(limit).Find(&data).Error; err != nil {
		logs.Info("FindWithLimit: %s", err.Error())
		return []*T{}
	}
	return data
}

// Count ...
func Count[T any](scopes ...func(db *gorm.DB) *gorm.DB) int {
	var count int64
	var model T
	if err := instance.gormDB.Model(&model).Scopes(scopes...).Count(&count).Error; err != nil {
		logs.Info("Count: %s", err.Error())
		return 0
	}

	return int(count)
}

// Raw ...
func Raw(query string, values ...any) *gorm.DB {
	return instance.gormDB.Raw(query, values...)
}

// Exec ...
func Exec(query string, values ...any) error {
	return instance.gormDB.Exec(query, values...).Error
}
