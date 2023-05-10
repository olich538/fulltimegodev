package db

import (
	"context"

	"github.com/olich538/fulltimegodev/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
)

const hotelColl = "hotels"

type HotelStore interface {
	InsertHotel(context.Context, *types.Hotel) (*types.Hotel, error)
}

type MongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHotelStore(c *mongo.Client, dbname string) *MongoHotelStore {
	return &MongoHotelStore{
		client: c,
		coll:   c.Database(dbname).Collection(hotelColl),
	}
}
