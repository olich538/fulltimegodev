package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/olich538/fulltimegodev/hotel-reservation/api"
	"github.com/olich538/fulltimegodev/hotel-reservation/db"
	"github.com/olich538/fulltimegodev/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client       *mongo.Client
	roomStore    db.RoomStore
	hotelStore   db.HotelStore
	userStore    db.UserStore
	bookingStore db.BookingStore
	ctx          = context.Background()
)

func seedRoom(size string, seaside bool, price float64, hotelID primitive.ObjectID) *types.Room {
	room := &types.Room{
		Size:    size,
		Seaside: seaside,
		Price:   price,
		HotelID: hotelID,
	}
	insertedRoom, err := roomStore.InsertRoom(context.Background(), room)
	if err != nil {
		log.Fatal(err)
	}
	return insertedRoom
}
func seedHotel(name string, location string, rating int) *types.Hotel {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	inserteHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}
	return inserteHotel
}
func seedBooking(userId, roomsId primitive.ObjectID, from, till time.Time) {
	booking := &types.Booking{
		UserID:   userId,
		RoomID:   roomsId,
		FromDate: from,
		TillDate: till,
	}
	resp, err := bookingStore.InsertBooking(context.Background(), booking)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("booking: ", resp.ID)

}
func seedUser(isAdmin bool, fname, lname, email, password string) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     email,
		FirstName: fname,
		LastName:  lname,
		Password:  password,
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = isAdmin
	_, err = userStore.InsertUser(ctx, user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s -> %s/n", user.Email, api.CreateTokenFromUser(user))
	return user
}

func main() {
	sunar := seedUser(false, "Sunal", "Booker", "sunario@gmek.com", "super_secret_password")
	seedUser(true, "admin", "admin", "admin@admin.com", "admin")
	seedHotel("Beluccia", "France", 4)
	seedHotel("The cozy", "Spain", 3)
	hotel := seedHotel("Kyiv", "Ukraine", 3)
	fmt.Println(hotel)
	seedRoom("small", true, 89.99, hotel.ID)
	seedRoom("medium", true, 100.99, hotel.ID)
	room := seedRoom("large", false, 120.99, hotel.ID)
	seedBooking(sunar.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 3))

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
	userStore = db.NewMongoUserStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	bookingStore = db.NewMongoBookingStore(client)
}
