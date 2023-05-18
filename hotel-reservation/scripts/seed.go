package main

import (
	"context"
	"log"

	"github.com/olich538/fulltimegodev/hotel-reservation/db"
	"github.com/olich538/fulltimegodev/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	roomStore  db.RoomStore
	hotelStore db.HotelStore
	ctx        = context.Background()
)

func seedHotel(name string, location string, rating int) {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	rooms := []types.Room{
		{Size: "normal",
			Price: 88.6},
		{Size: "small",
			Price: 222.6},
	}

	inserteHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	for _, room := range rooms {
		room.HotelID = inserteHotel.ID
		_, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func main() {
	seedHotel("Beluccia", "France", 4)
	seedHotel("The cozy", "Spain", 3)
	seedHotel("Kyiv", "Ukraine", 3)

}

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore = db.NewMongoHotelStore(client)
	// roomStore = db.NewMongoRoomStore(client, hotelStore)
}
