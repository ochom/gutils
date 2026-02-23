// Package nosql provides a generic MongoDB database abstraction layer with type-safe operations.
//
// This package simplifies MongoDB operations by providing generic CRUD functions that
// automatically handle collection naming based on struct types. Collection names are
// derived from struct names using snake_case with pluralization.
//
// Features:
//   - Type-safe generic CRUD operations
//   - Automatic collection name generation from struct names
//   - Aggregation pipeline support
//
// Example usage:
//
//	// Define your model
//	type User struct {
//		ID    string `bson:\"_id\"`
//		Name  string `bson:\"name\"`
//		Email string `bson:\"email\"`
//	}
//
//	// Initialize connection
//	err := nosql.Init(\"mongodb://localhost:27017\", \"mydb\")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Create a user
//	user := &User{ID: uuid.New(), Name: \"Alice\", Email: \"alice@example.com\"}
//	nosql.Create(user)
//
//	// Find by ID
//	user, err := nosql.FindOneByID[User](\"user-id\")
//
//	// Find with filter
//	users, err := nosql.FindAll[User](bson.M{\"email\": \"alice@example.com\"})\n//\n// Collection naming examples:\n//   - User -> users\n//   - UserGroup -> user_groups\n//   - Category -> categories\npackage nosql
package nosql

import (
	"context"
	"errors"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var conn *mongo.Database

// GetTableName converts a struct name to a MongoDB collection name.
// Uses snake_case with pluralization.
//
// Conversion rules:
//   - CamelCase becomes snake_case: UserGroup -> user_group
//   - Adds 's' suffix for pluralization: user_group -> user_groups
//   - Handles 'y' endings: Category -> categories
//
// Example:
//
//	name := nosql.GetTableName(\"User\")        // \"users\"
//	name := nosql.GetTableName(\"UserProfile\") // \"user_profiles\"
//	name := nosql.GetTableName(\"Category\")    // \"categories\"
func GetTableName(name string) string {
	if name == "" {
		panic("empty collection name")
	}

	var tableName string
	for i, r := range name {
		if i > 0 && 'A' <= r && r <= 'Z' {
			tableName += "_"
		}
		tableName += string(r)
	}

	// convert to lower case
	tableName = strings.ToLower(tableName)

	// pluralize
	if tableName[len(tableName)-1] == 'y' {
		tableName = tableName[:len(tableName)-1] + "ies"
	}

	if tableName[len(tableName)-1] != 's' {
		tableName += "s"
	}

	return tableName
}

// Init initializes the MongoDB connection.
// Must be called before using any other nosql functions.
//
// Example:
//
//	err := nosql.Init("mongodb://localhost:27017", "myapp")
//	if err != nil {
//		log.Fatal("Failed to connect to MongoDB:", err)
//	}
//
//	// With authentication
//	err := nosql.Init(
//		"mongodb://user:password@localhost:27017",
//		"myapp",
//	)
func Init(dsn, dbName string) error {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(dsn))
	if err != nil {
		return err
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		return err
	}

	conn = client.Database(dbName)
	return nil
}

// Col returns the MongoDB collection for the given type.
// The collection name is automatically derived from the type name.
//
// Example:
//
//	// Get collection for User type
//	col := nosql.Col(User{})
//	// col = collection "users"
//
//	// Use for custom operations
//	col := nosql.Col(User{})
//	cursor, err := col.Find(ctx, bson.M{"active": true})
func Col(v any) *mongo.Collection {
	var name string
	t := reflect.TypeOf(v)

	if t.Kind() == reflect.Ptr {
		name = t.Elem().Name()
	} else {
		name = t.Name()
	}

	return conn.Collection(GetTableName(name))
}

func getIDField(v interface{}) (id string, err error) {
	value := reflect.ValueOf(v)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	field := value.FieldByName("ID")
	if !field.IsValid() {
		return "", errors.New("object does not have ID field")
	}

	fieldValue := field.Interface()

	return fieldValue.(string), nil
}
