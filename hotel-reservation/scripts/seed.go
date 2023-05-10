package main

import (
	"context"
	"fmt"
	"log"

	"github.com/olich538/fulltimegodev/hotel-reservation/db"
	"github.com/olich538/fulltimegodev/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client, "hotels")

	hotel := types.Hotel{
		Name:     "Bellucia",
		Location: "Norway",
	}

	inserteHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	room := types.Room{
		Type:      types.SeaSideRoomType,
		BasePrice: 88.6,
	}
	_ = room
	fmt.Println(inserteHotel)
}
