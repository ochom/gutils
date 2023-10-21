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

// GetTableName returns collection name in lower case separated by underscore
// e.g User -> users, UserGroup -> user_groups
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

// Init initializes the database connection with MongoDB
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

// Col returns a collection
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
