package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/olich538/fulltimegodev/hotel-reservation/api"
	"github.com/olich538/fulltimegodev/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dburi          = "mongodb://localhost:27017"
	dbName         = "hotel-reservation"
	userCollection = "users"
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
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(client)
	// ctx := context.Background()
	// coll := client.Database(dbName).Collection(userCollection)
	// user := types.User{
	// 	FirstName: "James",
	// 	LastName:  "Mott",
	// }

	// _, err = coll.InsertOne(ctx, user)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// var james types.User
	// if err := coll.FindOne(ctx, bson.M{}).Decode(&james); err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(james)

	listenAddr := flag.String("listenAddr", ":5000", "The listen address of the API server")
	flag.Parse()

	app := fiber.New(config)
	userHandler := api.NewUserHandler(db.NewMongoUserStore(client, dbName))

	apiv1 := app.Group("/api/v1")
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)

	app.Listen(*listenAddr)
}
