package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fulltimegodev/hotel-reservation/api/middleware"
	"github.com/fulltimegodev/hotel-reservation/db/fixtures"
	"github.com/fulltimegodev/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func TestAdminGetBookings(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	var (
		adminUser = fixtures.AddUser(db.Store, "Adam", "Super", true)
		user      = fixtures.AddUser(db.Store, "james", "foo", false)
		hotel     = fixtures.AddHotel(db.Store, "bar hotel", "a", 4, nil)
		room      = fixtures.AddRoom(db.Store, "small", true, 4.4, hotel.ID)
		from      = time.Now()
		till      = from.AddDate(0, 0, 5)
		booking   = fixtures.AddBooking(db.Store, user.ID, room.ID, from, till)
		app       = fiber.New()
	)
	_ = booking
	admin := app.Group("/", middleware.JWTAuthentication(db.User), middleware.AdminAuth)

	bookingHandler := NewBookingHandler(db.Store)
	admin.Get("/", bookingHandler.HandleGetBookings)
	req := httptest.NewRequest("GET", "/", nil)

	req.Header.Add("X-Api-Token", CreateTokenFromUser(adminUser))

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("non 200 responce %d", resp.StatusCode)
	}
	var bookings []*types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}
	if len(bookings) != 1 {
		t.Fatalf("expecte 1 booking got %d", len(bookings))
	}
	// if !reflect.DeepEqual(booking, bookings[0]) {
	// 	fmt.Printf("%+v\n", booking)
	// 	fmt.Printf("%+v\n", bookings[0])

	// 	t.Fatal("expected bookings to be equal got")
	// }
	if booking.ID != bookings[0].ID || booking.UserID != bookings[0].UserID {
		t.Fatal("Booking ID and/or User ID not match")
	}
	fmt.Println(bookings)

	// test non admin cannot access the bookings
	rn := httptest.NewRequest("GET", "/", nil)

	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))

	resp, err = app.Test(rn)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("expected a non 200 status code got  %d", resp.StatusCode)
	}

}

func TestUserGetBookings(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	var (
		nonAuthUser    = fixtures.AddUser(db.Store, "Fake", "User", false)
		user           = fixtures.AddUser(db.Store, "james", "foo", false)
		hotel          = fixtures.AddHotel(db.Store, "bar hotel", "a", 4, nil)
		room           = fixtures.AddRoom(db.Store, "small", true, 4.4, hotel.ID)
		from           = time.Now()
		till           = from.AddDate(0, 0, 5)
		booking        = fixtures.AddBooking(db.Store, user.ID, room.ID, from, till)
		app            = fiber.New()
		bookingHandler = NewBookingHandler(db.Store)
		route          = app.Group("/", middleware.JWTAuthentication(db.User))
	)

	route.Get("/:id", bookingHandler.HandleGetBooking)
	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)

	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("non 200 responce %d", resp.StatusCode)
	}
	var bookingResp *types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookingResp); err != nil {
		t.Fatal(err)
	}
	if bookingResp.ID != booking.ID {
		t.Fatal("Booking ID does not match")
	}
	if bookingResp.UserID != booking.UserID {
		t.Fatal("Booking User ID does not match")
	}
	req = httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(nonAuthUser))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("non 200 responce %d", resp.StatusCode)
	}

}
