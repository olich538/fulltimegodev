package main

import (
	"context"
	"fmt"
	"log"

	"github.com/olich538/fulltimegodev/hotel-reservation/db"
	"github.com/olich538/fulltimegodev/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client)
	roomStore := db.NewMongoRoomStore(client, hotelStore)

	hotel := types.Hotel{
		Name:     "Bellucia",
		Location: "Norway",
		Rooms:    []primitive.ObjectID{},
	}

	// room := types.Room{
	// 	Type: ,
	// }
	inserteHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}
	rooms := []types.Room{
		{Type: types.SeaSideRoomType,
			BasePrice: 88.6},
		{Type: types.DeluxRoomType,
			BasePrice: 199.6},
	}
	roomA := types.Room{
		Type:      types.SeaSideRoomType,
		BasePrice: 88.6,
	}

	for _, room := range rooms {
		room.HotelID = inserteHotel.ID
		insertedRoom, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(insertedRoom)

	}
	roomA.HotelID = inserteHotel.ID
	insertedRoom, err := roomStore.InsertRoom(ctx, &roomA)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(inserteHotel)
	fmt.Println(insertedRoom)

}
