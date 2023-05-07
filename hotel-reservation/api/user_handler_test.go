package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)
	if user.FirstName != params.FirstName {
		t.Errorf("expected firstname %s but got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("expected lastname %s but got %s", params.LastName, user.LastName)
	}
	if len(user.EncryptedPassword) > 0 {
		t.Errorf("expected Encrypted Password not to be included in json")
	}
	if user.Email != params.Email {
		t.Errorf("expected email %s but got %s", params.Email, user.Email)
	}
	if len(user.ID) == 0 {
		t.Errorf("expecting user id to be set")
	}
	// bb, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	t.Error(err)
	// }

	fmt.Println(user)
}
