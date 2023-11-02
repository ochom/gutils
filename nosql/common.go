package nosql

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Create ...
func Create[T any](ctx context.Context, v *T) error {
	_, err := Col(v).InsertOne(ctx, v)
	return err
}

// Update ...
func Update[T any](ctx context.Context, v *T) error {
	id, err := getIDField(v)
	if err != nil {
		return err
	}

	_, err = Col(v).ReplaceOne(ctx, bson.M{"_id": id}, v)
	return err
}

// Delete ...
func Delete[T any](ctx context.Context, filter bson.M) error {
	var v T
	_, err := Col(v).DeleteOne(ctx, filter)
	return err
}

// DeleteByID ...
func DeleteByID[T any](ctx context.Context, id string) error {
	return Delete[T](ctx, bson.M{"_id": id})
}

// FindOne ...
func FindOne[T any](ctx context.Context, filter bson.M) (*T, error) {
	var v T
	if err := Col(v).FindOne(ctx, filter).Decode(&v); err != nil {
		return nil, err
	}

	return &v, nil
}

// FindOneByID ...
func FindOneByID[T any](ctx context.Context, id string) (*T, error) {
	return FindOne[T](ctx, bson.M{"_id": id})
}

// FindAll ...
func FindAll[T any](ctx context.Context, filter bson.M) ([]*T, error) {
	var v T
	cur, err := Col(v).Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)
	var vs []*T
	if err := cur.All(ctx, &vs); err != nil {
		return nil, err
	}

	return vs, nil
}

// FindWithLimit ...
func FindWithLimit[T any](ctx context.Context, filter bson.M, limit int64) ([]*T, error) {
	var v T
	cur, err := Col(v).Find(ctx, filter, options.Find().SetLimit(limit))
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)
	var vs []*T
	if err := cur.All(ctx, &vs); err != nil {
		return nil, err
	}

	return vs, nil
}

// Count ...
func Count[T any](ctx context.Context, filter bson.M) (int64, error) {
	var v T
	return Col(v).CountDocuments(ctx, filter)
}
