package nosql_test

import (
	"context"
	"testing"

	"github.com/ochom/gutils/nosql"
	"github.com/ochom/gutils/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func initDB(t *testing.T) {
	if err := nosql.Init("mongodb://root:example@localhost:27017/?authSource=admin", "test"); err != nil {
		t.Fatal(err)
	}
}

type User struct {
	ID   string `bson:"_id"`
	Name string `bson:"name"`
}

func TestCRUD(t *testing.T) {
	initDB(t)

	// test create
	user := &User{ID: uuid.NewString(), Name: "John Doe"}
	if err := nosql.Create(context.Background(), user); err != nil {
		t.Fatalf("[create] %s", err.Error())
	}

	// test update
	user.Name = "John Doom"
	if err := nosql.Update(context.Background(), user); err != nil {
		t.Fatalf("[update] %s", err.Error())
	}

	// test find one
	u, err := nosql.FindOne[User](context.Background(), bson.M{"_id": user.ID})
	if err != nil {
		t.Fatalf("[find one] %s", err.Error())
	}

	if u.Name != "John Doom" {
		t.Fatalf("[find one] expected name to be John Doom, got %s", u.Name)
	}

	// test find all
	users, err := nosql.FindAll[User](context.Background(), bson.M{"_id": user.ID})
	if err != nil {
		t.Fatalf("[find all] %s", err.Error())
	}

	if len(users) != 1 {
		t.Fatalf("[find all] expected 1 user, got %d", len(users))
	}

	// test delete
	if err := nosql.Delete[User](context.Background(), bson.M{"_id": user.ID}); err != nil {
		t.Fatalf("[delete] %s", err.Error())
	}

}
