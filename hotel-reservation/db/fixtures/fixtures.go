package fixtures

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/olich538/fulltimegodev/hotel-reservation/db"
	"github.com/olich538/fulltimegodev/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddUser(store *db.Store, fname, lname string, admin bool) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     fmt.Sprintf("%s@%sgmail.com", fname, lname),
		FirstName: fname,
		LastName:  lname,
		Password:  fmt.Sprintf("%s_%s", fname, lname),
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = admin
	var insertedUser *types.User
	insertedUser, err = store.User.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	return insertedUser
}

func AddHotel(s *db.Store, location, name string, rating int, rooms []primitive.ObjectID) *types.Hotel {
	var roomIDs = rooms
	if rooms == nil {
		roomIDs = []primitive.ObjectID{}
	}
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    roomIDs,
		Rating:   rating,
	}

	inserteHotel, err := s.Hotel.InsertHotel(context.TODO(), &hotel)
	if err != nil {
		log.Fatal(err)
	}
	return inserteHotel
}

func AddRoom(s *db.Store, size string, seaside bool, price float64, hid primitive.ObjectID) *types.Room {

	room := &types.Room{
		Size:    size,
		Seaside: seaside,
		Price:   price,
		HotelID: hid,
	}
	insertedRoom, err := s.Room.InsertRoom(context.Background(), room)
	if err != nil {
		log.Fatal(err)
	}
	return insertedRoom
}

func AddBoooking(s *db.Store, uid, rid primitive.ObjectID, from, till time.Time) *types.Booking {
	booking := &types.Booking{
		UserID:   uid,
		RoomID:   rid,
		FromDate: from,
		TillDate: till,
	}
	resp, err := s.Booking.InsertBooking(context.Background(), booking)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("booking: ", resp.ID)
	return resp
}
