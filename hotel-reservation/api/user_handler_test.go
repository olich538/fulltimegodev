package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/olich538/fulltimegodev/hotel-reservation/db"
	"github.com/olich538/fulltimegodev/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	testdburi = "mongodb://localhost:27017"
	dbname    = "hotel-reservation-test"
)

type testdb struct {
	db.UserStore
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testdburi))
	if err != nil {
		log.Fatal(err)
	}
	return &testdb{
		UserStore: db.NewMongoUserStore(client, dbname),
	}

}

func TestCreateUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.UserStore)
	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		Email:     "some@test.foo",
		FirstName: "James",
		LastName:  "Foo",
		Password:  "dfvwerfgwerfgv",
	}

	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	// var m map[string]any
	// err = json.NewDecoder(resp.Request.Body).Decode(&m)
	// if err != nil {
	// 	t.Error(err)
	// }
	// fmt.Println(resp.Status)
	bb, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(string(bb))

}
