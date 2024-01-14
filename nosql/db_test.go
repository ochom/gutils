package nosql_test

import (
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
	if err := nosql.Create(user); err != nil {
		t.Fatalf("[create] %s", err.Error())
	}

	// test update
	user.Name = "John Does"
	if err := nosql.Update(user); err != nil {
		t.Fatalf("[update] %s", err.Error())
	}

	// test find one
	u, err := nosql.FindOne[User](bson.M{"_id": user.ID})
	if err != nil {
		t.Fatalf("[find one] %s", err.Error())
	}

	if u.Name != "John Does" {
		t.Fatalf("[find one] expected name to be John Doom, got %s", u.Name)
	}

	// test find all
	users, err := nosql.FindAll[User](bson.M{"_id": user.ID})
	if err != nil {
		t.Fatalf("[find all] %s", err.Error())
	}

	if len(users) != 1 {
		t.Fatalf("[find all] expected 1 user, got %d", len(users))
	}

	// test delete
	if err := nosql.Delete[User](bson.M{"_id": user.ID}); err != nil {
		t.Fatalf("[delete] %s", err.Error())
	}

}

func TestGetTableName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Student",
			args{"Student"},
			"students",
		},
		{
			"Students",
			args{"Students"},
			"students",
		},
		{
			"StudentInfo",
			args{"StudentInfo"},
			"student_infos",
		},
		{
			"Library",
			args{"Library"},
			"libraries",
		},
		{
			"LibraryBook",
			args{"LibraryBook"},
			"library_books",
		},
		{
			"LibraryBooks",
			args{"LibraryBooks"},
			"library_books",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nosql.GetTableName(tt.args.name); got != tt.want {
				t.Errorf("GetTableName() = %v, want %v", got, tt.want)
			}
		})
	}
}
