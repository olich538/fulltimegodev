package api

import (
	"context"
	"log"
	"testing"

	"github.com/olich538/fulltimegodev/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	testdburi = "mongodb://localhost:27017"
	dbname    = "hotel-reservation-test"
)

type testdb struct {
	client *mongo.Client
	*db.Store
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.client.Database(db.DBNAME).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testdburi))
	if err != nil {
		log.Fatal(err)
	}

	// hotelStore := db.NewMongoHotelStore(client)
	return &testdb{
		client: client,

		// User:    db.NewMongoUserStore(client),
		// Room:    db.NewMongoRoomStore(client, hotelStore),
		// Booking: db.NewMongoBookingStore(client),
		// Hotel:   db.NewMongoHotelStore(client),
	}

}
