package sqlr

import (
	"context"

	"github.com/ochom/gutils/logs"
	"gorm.io/gorm"
)

// Create inserts a new record into the database.
// The table name is automatically derived from the type.
//
// Example:
//
//	user := &User{Name: \"Alice\", Email: \"alice@example.com\"}
//	err := sqlr.Create(user)
//	if err != nil {
//		log.Error(\"Failed to create user: %v\", err)
//	}
//	fmt.Println(\"Created user with ID:\", user.ID)
func Create[T any](data *T) error {
	return instance.gormDB.Create(data).Error
}

// CreateWithCtx inserts a new record with context support for cancellation/timeout.
//
// Example:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	err := sqlr.CreateWithCtx(ctx, user)
func CreateWithCtx[T any](ctx context.Context, data *T) error {
	return instance.gormDB.WithContext(ctx).Create(data).Error
}

// Update saves all fields of an existing record.
// The record must have a primary key set.
//
// Example:
//
//	user, _ := sqlr.FindOneById[User](1)
//	user.Name = \"New Name\"
//	err := sqlr.Update(user)
func Update[T any](data *T) error {
	return instance.gormDB.Save(data).Error
}

// UpdateWithCtx saves a record with context support.
//
// Example:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	err := sqlr.UpdateWithCtx(ctx, user)
func UpdateWithCtx[T any](ctx context.Context, data *T) error {
	return instance.gormDB.WithContext(ctx).Save(data).Error
}

// UpdateOne updates specific fields on records matching the scope.
// Only the fields in the updates map are modified.
//
// Example:
//
//	// Update user's status
//	err := sqlr.UpdateOne[User](
//		func(db *gorm.DB) *gorm.DB {
//			return db.Where(\"id = ?\", userID)
//		},
//		map[string]any{\"status\": \"active\", \"updated_at\": time.Now()},
//	)
func UpdateOne[T any](scope func(db *gorm.DB) *gorm.DB, updates map[string]any) error {
	var model T
	return instance.gormDB.Model(&model).Scopes(scope).Updates(updates).Error
}

// Delete removes records matching the provided scopes.
// WARNING: Without scopes, this will delete all records!
//
// Example:
//
//	// Delete inactive users
//	err := sqlr.Delete[User](func(db *gorm.DB) *gorm.DB {
//		return db.Where(\"active = ?\", false)
//	})
func Delete[T any](scopes ...func(db *gorm.DB) *gorm.DB) error {
	return instance.gormDB.Scopes(scopes...).Delete(new(T)).Error
}

// DeleteWithCtx removes records with context support.
//
// Example:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	err := sqlr.DeleteWithCtx[User](ctx, func(db *gorm.DB) *gorm.DB {
//		return db.Where(\"expired_at < ?\", time.Now())
//	})
func DeleteWithCtx[T any](ctx context.Context, scopes ...func(db *gorm.DB) *gorm.DB) error {
	return instance.gormDB.WithContext(ctx).Scopes(scopes...).Delete(new(T)).Error
}

// DeleteById removes a record by its primary key ID.
// Additional scopes can be provided for extra conditions.
//
// Example:
//
//	err := sqlr.DeleteById[User](123)
//
//	// With soft-delete scope
//	err := sqlr.DeleteById[User](123, func(db *gorm.DB) *gorm.DB {
//		return db.Unscoped() // Hard delete
//	})
func DeleteById[T any](id any, scopes ...func(db *gorm.DB) *gorm.DB) error {
	return Delete[T](append(scopes, func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	})...)
}

// FindOne retrieves the first record matching the provided scopes.
// Returns nil and an error if no record is found.
//
// Example:
//
//	// Find by email
//	user, err := sqlr.FindOne[User](func(db *gorm.DB) *gorm.DB {
//		return db.Where("email = ?", email)
//	})
//	if err != nil {
//		// User not found or error
//	}
//
//	// With multiple conditions
//	product, err := sqlr.FindOne[Product](func(db *gorm.DB) *gorm.DB {
//		return db.Where("sku = ? AND active = ?", sku, true)
//	})
func FindOne[T any](scopes ...func(db *gorm.DB) *gorm.DB) (*T, error) {
	var data T
	if err := instance.gormDB.Scopes(scopes...).First(&data).Error; err != nil {
		return nil, err
	}

	return &data, nil
}

// FindOneById retrieves a record by its primary key.
// Additional scopes can be provided for extra conditions (e.g., preloading).
//
// Example:
//
//	user, err := sqlr.FindOneById[User](123)
//
//	// With preloading
//	order, err := sqlr.FindOneById[Order](orderID, func(db *gorm.DB) *gorm.DB {
//		return db.Preload("Items")
//	})
func FindOneById[T any](id any, scopes ...func(db *gorm.DB) *gorm.DB) (*T, error) {
	return FindOne[T](append(scopes, func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	})...)
}

// FindAll retrieves all records matching the provided scopes.
// Returns an empty slice if no records are found (never returns nil).
//
// Example:
//
//	// Find all active users
//	users := sqlr.FindAll[User](func(db *gorm.DB) *gorm.DB {
//		return db.Where("active = ?", true).Order("name ASC")
//	})
//
//	// Find all records
//	products := sqlr.FindAll[Product]()
func FindAll[T any](scopes ...func(db *gorm.DB) *gorm.DB) []*T {
	data := []*T{}
	if err := instance.gormDB.Scopes(scopes...).Find(&data).Error; err != nil {
		logs.Info("FindAll: %s", err.Error())
		return []*T{}
	}

	return data
}

// FindWithLimit retrieves records with pagination support.
// Offset is calculated as (page - 1) * limit.
//
// Example:
//
//	// Get page 1 with 10 items per page
//	users := sqlr.FindWithLimit[User](1, 10)
//
//	// Get page 3 with 20 items, sorted by created_at
//	posts := sqlr.FindWithLimit[Post](3, 20, func(db *gorm.DB) *gorm.DB {
//		return db.Order("created_at DESC")
//	})
func FindWithLimit[T any](page, limit int, scopes ...func(db *gorm.DB) *gorm.DB) []*T {
	data := []*T{}
	if err := instance.gormDB.Scopes(scopes...).Offset((page - 1) * limit).Limit(limit).Find(&data).Error; err != nil {
		logs.Info("FindWithLimit: %s", err.Error())
		return []*T{}
	}
	return data
}

// Count returns the number of records matching the provided scopes.
//
// Example:
//
//	// Count active users
//	count := sqlr.Count[User](func(db *gorm.DB) *gorm.DB {
//		return db.Where("active = ?", true)
//	})
//
//	// Count all products
//	total := sqlr.Count[Product]()
func Count[T any](scopes ...func(db *gorm.DB) *gorm.DB) int {
	var count int64
	if err := instance.gormDB.Model(new(T)).Scopes(scopes...).Count(&count).Error; err != nil {
		logs.Info("Count: %s", err.Error())
		return 0
	}

	return int(count)
}

// Exists checks if any record matches the provided scopes.
// More efficient than Count for existence checks.
//
// Example:
//
//	// Check if email exists
//	if sqlr.Exists[User](func(db *gorm.DB) *gorm.DB {
//		return db.Where("email = ?", email)
//	}) {
//		return errors.Conflict("email already registered")
//	}
func Exists[T any](scopes ...func(db *gorm.DB) *gorm.DB) bool {
	query := instance.gormDB.Model(new(T)).Scopes(scopes...).Select("1").Limit(1)
	var exists bool
	if err := query.Scan(&exists).Error; err != nil {
		logs.Info("Exists: %s", err.Error())
		return false
	}

	return exists
}

// Raw executes a raw SQL query and returns the GORM DB for further processing.
//
// Example:
//
//	var result struct {
//		Total int
//		Avg   float64
//	}
//	sqlr.Raw("SELECT COUNT(*) as total, AVG(price) as avg FROM products WHERE active = ?", true).Scan(&result)
func Raw(query string, values ...any) *gorm.DB {
	return instance.gormDB.Raw(query, values...)
}

// Exec executes a raw SQL statement that doesn't return rows.
//
// Example:
//
//	// Update multiple records
//	err := sqlr.Exec("UPDATE users SET status = ? WHERE last_login < ?", "inactive", oneYearAgo)
//
//	// Delete old records
//	err := sqlr.Exec("DELETE FROM logs WHERE created_at < ?", thirtyDaysAgo)
func Exec(query string, values ...any) error {
	return instance.gormDB.Exec(query, values...).Error
}

// Transact executes multiple operations within a database transaction.
// If any function returns an error, the transaction is rolled back.
// All functions are executed in order within the same transaction.
//
// Example:
//
//	err := sqlr.Transact(
//		func(tx *gorm.DB) error {
//			return tx.Create(&order).Error
//		},
//		func(tx *gorm.DB) error {
//			return tx.Create(&payment).Error
//		},
//		func(tx *gorm.DB) error {
//			return tx.Model(&inventory).Update("quantity", gorm.Expr("quantity - ?", order.Quantity)).Error
//		},
//	)
//	if err != nil {
//		// Transaction was rolled back
//	}
func Transact(fn ...func(tx *gorm.DB) error) error {
	err := instance.gormDB.Transaction(func(db *gorm.DB) error {
		for _, f := range fn {
			if err := f(db); err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

// TransactWithCtx executes operations within a transaction with context support.
// Useful for timeout/cancellation control on long-running transactions.
//
// Example:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//	defer cancel()
//
//	err := sqlr.TransactWithCtx(ctx,
//		func(tx *gorm.DB) error {
//			return tx.Create(&order).Error
//		},
//	)
func TransactWithCtx(ctx context.Context, fn ...func(tx *gorm.DB) error) error {
	err := instance.gormDB.WithContext(ctx).Transaction(func(db *gorm.DB) error {
		for _, f := range fn {
			if err := f(db); err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

// Migrate automatically migrates the database schema for the provided models.
// Creates tables, adds missing columns, and creates indexes.
//
// Note: Does not delete columns or change column types for safety.
//
// Example:
//
//	// Migrate single model
//	sqlr.Migrate(&User{})
//
//	// Migrate multiple models
//	sqlr.Migrate(&User{}, &Post{}, &Comment{})
func Migrate(models ...any) error {
	return instance.gormDB.AutoMigrate(models...)
}
