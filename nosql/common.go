package nosql

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Create ...
func Create[T any](v *T) error {
	ctx := context.Background()
	_, err := Col(v).InsertOne(ctx, v)
	return err
}

// CreateMany ...
func CreateMany[T any](vs []*T) error {
	ctx := context.Background()
	newData := make([]interface{}, len(vs))
	for i, v := range vs {
		newData[i] = v
	}
	_, err := Col(vs[0]).InsertMany(ctx, newData)
	return err
}

// Update ...
func Update[T any](v *T) error {
	ctx := context.Background()
	id, err := getIDField(v)
	if err != nil {
		return err
	}

	_, err = Col(v).ReplaceOne(ctx, bson.M{"_id": id}, v)
	return err
}

// UpdateMany ...
func UpdateMany[T any](vs []*T) error {
	ctx := context.Background()
	newData := make([]interface{}, len(vs))
	for i, v := range vs {
		newData[i] = v
	}

	_, err := Col(vs[0]).BulkWrite(ctx, []mongo.WriteModel{
		&mongo.ReplaceOneModel{
			Filter:      bson.M{"_id": bson.M{"$in": newData}},
			Replacement: newData,
		},
	})
	return err
}

// Delete ...
func Delete[T any](filter bson.M) error {
	ctx := context.Background()
	var v T
	_, err := Col(v).DeleteOne(ctx, filter)
	return err
}

// DeleteByID ...
func DeleteByID[T any](id string) error {
	return Delete[T](bson.M{"_id": id})
}

// FindOne ...
func FindOne[T any](filter bson.M) (*T, error) {
	ctx := context.Background()
	var v T
	if err := Col(v).FindOne(ctx, filter).Decode(&v); err != nil {
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
	ctx := context.Background()
	var v T
	vs = []*T{}
	cur, err := Col(v).Find(ctx, filter)
	if err != nil {
		return
	}

	defer cur.Close(ctx)

	err = cur.All(ctx, &vs)
	if err != nil {
		return
	}

	return vs, nil
}

// FindWithLimit ...
func FindWithLimit[T any](filter bson.M, limit int64) (vs []*T, err error) {
	ctx := context.Background()
	var v T
	vs = []*T{}
	cur, err := Col(v).Find(ctx, filter, options.Find().SetLimit(limit))
	if err != nil {
		return
	}

	defer cur.Close(ctx)

	err = cur.All(ctx, &vs)
	if err != nil {
		return
	}

	return vs, nil
}

// Count ...
func Count[T any](filter bson.M) (int64, error) {
	ctx := context.Background()
	var v T
	return Col(v).CountDocuments(ctx, filter)
}

// Pipe ...
func Pipe[T any](pipeline []bson.M) (vs []*T, err error) {
	ctx := context.Background()
	var v T
	vs = []*T{}
	cur, err := Col(v).Aggregate(ctx, pipeline)
	if err != nil {
		return
	}

	defer cur.Close(ctx)

	err = cur.All(ctx, &vs)
	if err != nil {
		return
	}

	return vs, nil
}
