package nosql

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Create ...
func Create[T any](v *T) error {
	_, err := Col(v).InsertOne(context.Background(), v)
	return err
}

// Update ...
func Update[T any](v *T) error {
	id, err := getIDField(v)
	if err != nil {
		return err
	}

	_, err = Col(v).ReplaceOne(context.Background(), bson.M{"_id": id}, v)
	return err
}

// Delete ...
func Delete[T any](filter bson.M) error {
	var v T
	_, err := Col(v).DeleteOne(context.Background(), filter)
	return err
}

// DeleteByID ...
func DeleteByID[T any](id string) error {
	return Delete[T](bson.M{"_id": id})
}

// FindOne ...
func FindOne[T any](filter bson.M) (*T, error) {
	var v T
	if err := Col(v).FindOne(context.Background(), filter).Decode(&v); err != nil {
		return nil, err
	}

	return &v, nil
}

// FindOneByID ...
func FindOneByID[T any](id string) (*T, error) {
	return FindOne[T](bson.M{"_id": id})
}

// FindAll ...
func FindAll[T any](filter bson.M) (vs []*T, err error) {
	var v T
	vs = []*T{}
	cur, err := Col(v).Find(context.Background(), filter)
	if err != nil {
		return
	}

	defer cur.Close(context.Background())

	err = cur.All(context.Background(), &vs)
	if err != nil {
		return
	}

	return vs, nil
}

// FindWithLimit ...
func FindWithLimit[T any](filter bson.M, limit int64) (vs []*T, err error) {
	var v T
	vs = []*T{}
	cur, err := Col(v).Find(context.Background(), filter, options.Find().SetLimit(limit))
	if err != nil {
		return
	}

	defer cur.Close(context.Background())

	err = cur.All(context.Background(), &vs)
	if err != nil {
		return
	}

	return vs, nil
}

// Count ...
func Count[T any](filter bson.M) (int64, error) {
	var v T
	return Col(v).CountDocuments(context.Background(), filter)
}
