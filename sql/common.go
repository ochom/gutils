package sql

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
	data := []*T{}
	err := conn.Find(&data, query).Error
	return data, err
}

// FindWithLimit ...
func FindWithLimit[T any](query *T, limit int) ([]*T, error) {
	data := []*T{}
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
