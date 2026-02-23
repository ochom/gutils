package nosql

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Create inserts a new document into the collection.
// The collection is automatically determined from the type of v.
//
// Example:
//
//	user := &User{ID: uuid.New(), Name: "Alice", Email: "alice@example.com"}
//	err := nosql.Create(user)
//	if err != nil {
//		log.Error("Failed to create user: %v", err)
//	}
func Create[T any](v *T) error {
	ctx := context.Background()
	_, err := Col(v).InsertOne(ctx, v)
	return err
}

// Update replaces an existing document with the provided document.
// The document must have an ID field set for identification.
//
// Example:
//
//	user, _ := nosql.FindOneByID[User]("user-id")
//	user.Name = "New Name"
//	err := nosql.Update(user)
func Update[T any](v *T) error {
	ctx := context.Background()
	id, err := getIDField(v)
	if err != nil {
		return err
	}

	_, err = Col(v).ReplaceOne(ctx, bson.M{"_id": id}, v)
	return err
}

// Delete removes documents matching the filter.
//
// Example:
//
//	// Delete by custom filter
//	err := nosql.Delete[User](bson.M{"email": "old@example.com"})
//
//	// Delete inactive users
//	err := nosql.Delete[User](bson.M{"active": false})
func Delete[T any](filter bson.M) error {
	ctx := context.Background()
	var v T
	_, err := Col(v).DeleteOne(ctx, filter)
	return err
}

// DeleteByID removes a document by its _id field.
//
// Example:
//
//	err := nosql.DeleteByID[User]("user-id-to-delete")
//	if err != nil {
//		log.Error("Failed to delete user: %v", err)
//	}
func DeleteByID[T any](id string) error {
	return Delete[T](bson.M{"_id": id})
}

// FindOne retrieves a single document matching the filter.
// Returns nil and an error if no document is found.
//
// Example:
//
//	// Find by email
//	user, err := nosql.FindOne[User](bson.M{"email": "alice@example.com"})
//	if err != nil {
//		// User not found or error
//	}
//
//	// Find with multiple conditions
//	user, err := nosql.FindOne[User](bson.M{
//		"email":  "alice@example.com",
//		"active": true,
//	})
func FindOne[T any](filter bson.M) (*T, error) {
	ctx := context.Background()
	var v T
	if err := Col(v).FindOne(ctx, filter).Decode(&v); err != nil {
		return nil, err
	}

	return &v, nil
}

// FindOneByID retrieves a document by its _id field.
// Convenience wrapper around FindOne with _id filter.
//
// Example:
//
//	user, err := nosql.FindOneByID[User]("507f1f77bcf86cd799439011")
//	if err != nil {
//		return errors.NotFound("user not found")
//	}
func FindOneByID[T any](id string) (*T, error) {
	return FindOne[T](bson.M{"_id": id})
}

// FindAll retrieves all documents matching the filter.
// Returns an empty slice if no documents are found.
//
// Example:
//
//	// Find all active users
//	users, err := nosql.FindAll[User](bson.M{"active": true})
//	for _, user := range users {
//		fmt.Println(user.Name)
//	}
//
//	// Find all documents
//	products, err := nosql.FindAll[Product](bson.M{})
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

// FindWithLimit retrieves documents matching the filter with a limit.
// Useful for pagination or getting top N results.
//
// Example:
//
//	// Get the 10 most recent orders
//	orders, err := nosql.FindWithLimit[Order](bson.M{}, 10)
//
//	// Get 5 active users
//	users, err := nosql.FindWithLimit[User](bson.M{"active": true}, 5)
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

// Count returns the number of documents matching the filter.
//
// Example:
//
//	// Count active users
//	count, err := nosql.Count[User](bson.M{"active": true})
//	fmt.Printf("Active users: %d\n", count)
//
//	// Count all documents
//	total, err := nosql.Count[Product](bson.M{})
func Count[T any](filter bson.M) (int64, error) {
	ctx := context.Background()
	var v T
	return Col(v).CountDocuments(ctx, filter)
}

// Pipe executes an aggregation pipeline and returns the results.
// Use for complex queries that require grouping, lookups, or transformations.
//
// Example:
//
//	// Group orders by status and count
//	pipeline := []bson.M{
//		{"$group": bson.M{
//			"_id":   "$status",
//			"count": bson.M{"$sum": 1},
//		}},
//	}
//	results, err := nosql.Pipe[OrderStats](pipeline)
//
//	// Lookup related documents
//	pipeline := []bson.M{
//		{"$lookup": bson.M{
//			"from":         "orders",
//			"localField":   "_id",
//			"foreignField": "user_id",
//			"as":           "orders",
//		}},
//	}
//	users, err := nosql.Pipe[UserWithOrders](pipeline)
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
