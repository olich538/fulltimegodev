package main

import (
	"context"
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/olich538/fulltimegodev/hotel-reservation/api"
	"github.com/olich538/fulltimegodev/hotel-reservation/db"
	"github.com/olich538/fulltimegodev/hotel-reservation/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		// code := fiber.StatusInternalServerError
		return ctx.JSON(map[string]string{"error": err.Error()})
		// Retrieve the custom status code if it's a *fiber.Error
		// var e *fiber.Error
		// if errors.As(err, &e) {
		// 	code = e.Code
		// }
		// // Send custom error page
		// err = ctx.Status(code).SendFile(fmt.Sprintf("./%d.html", code))
		// if err != nil {
		// 	// in case sendfile failed
		// 	return ctx.Status(fiber.StatusInternalServerError).SendString("Internal server error")
		// }
		// return nil
	},
}

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	listenAddr := flag.String("listenAddr", ":5000", "The listen address of the API server")
	flag.Parse()

	// handler initialization
	var (
		hotelstore = db.NewMongoHotelStore(client)
		roomStore  = db.NewMongoRoomStore(client, hotelstore)
		userStore  = db.NewMongoUserStore(client)
		app        = fiber.New(config)
		store      = &db.Store{
			Hotel: hotelstore,
			Room:  roomStore,
			User:  userStore,
		}
		hotelHandler = api.NewHotelHandler(store)
		userHandler  = api.NewUserHandler(userStore)
		authHandler  = api.NewAuthHandler(userStore)
		auth         = app.Group("/api")
		apiv1        = app.Group("/api/v1", middleware.JWTAuthentication)
	)

	// auth
	auth.Post("/auth", authHandler.HandleAuthenticate)

	// user handlers
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)

	//hotel handlers
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotelByID)
	app.Listen(*listenAddr)
}
