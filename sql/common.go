package sql

import "gorm.io/gorm"

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
func FindOne[T any](query *T, scopes ...func(*gorm.DB) *gorm.DB) (*T, error) {
	var data T
	err := conn.Scopes(scopes...).First(&data, query).Error
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// FindAll ...
func FindAll[T any](query *T, scopes ...func(*gorm.DB) *gorm.DB) ([]*T, error) {
	data := []*T{}
	err := conn.Scopes(scopes...).Find(&data, query).Error
	return data, err
}

// FindWithLimit ...
func FindWithLimit[T any](query *T, page, limit int, scopes ...func(*gorm.DB) *gorm.DB) ([]*T, error) {
	data := []*T{}

	err := conn.Scopes(scopes...).Offset((page-1)*limit).Limit(limit).Find(&data, query).Error
	return data, err
}

// Count ...
func Count[T any](query *T, scopes ...func(*gorm.DB) *gorm.DB) (int64, error) {
	var count int64
	var model T
	err := conn.Model(&model).Where(query).Count(&count).Error
	return count, err
}
