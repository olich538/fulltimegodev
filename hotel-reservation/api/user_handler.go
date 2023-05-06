package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/olich538/fulltimegodev/hotel-reservation/db"
	"github.com/olich538/fulltimegodev/hotel-reservation/types"
)

type UserHandlerStore struct {
	userStore db.UserStore
}

// constructor function
func NewUserHandler(userStore db.UserStore) *UserHandlerStore {
	return &UserHandlerStore{
		userStore: userStore,
	}
}

func (h *UserHandlerStore) HandleGetUser(c *fiber.Ctx) error {
	var (
		id = c.Params("id")
	)
	user, err := h.userStore.GetUserByID(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(user)
}
func (h *UserHandlerStore) HandleGetUsers(c *fiber.Ctx) error {
	// u := types.User{
	// 	FirstName: "James",
	// 	LastName:  "WaterKooler",
	// }
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(&users)
}

func (h *UserHandlerStore) HandlePostUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if errors := params.ValidateUserParams(); len(errors) > 0 {
		return c.JSON(errors)
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}

	insertedUser, err := h.userStore.InsertUser(c.Context(), user)
	if err != nil {
		return err
	}
	return c.JSON(insertedUser)
}
