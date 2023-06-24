package api

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/olich538/fulltimegodev/hotel-reservation/db"
	"github.com/olich538/fulltimegodev/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type AuthResponse struct {
	User  *types.User `json:"user"`
	Token string      `json:"token"`
}
type genericResponce struct {
	Type string `json: "type"`
	Msg  string `json: "msg"`
}

func invalidCredentials(c *fiber.Ctx) error {
	return c.Status(http.StatusBadRequest).JSON(genericResponce{
		Type: "error",
		Msg:  "invalid credentials",
	})
}

// a handler should only do:
// - serialization of the incomming requests
// - do some data fetching from DB
// - return the data back to user
// - call some bussines logic

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var params AuthParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	fmt.Println(params)
	user, err := h.userStore.GetUserByEmail(c.Context(), params.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// return fmt.Errorf("invalid credentials")
			return invalidCredentials(c)
		}
		return err
	}
	if !types.IsValidPassword(user.EncryptedPassword, params.Password) {
		return invalidCredentials(c)
	}
	// err = bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(authParams.Password))
	// if err != nil {
	// 	return fmt.Errorf("invalid credentials")
	// }
	resp := AuthResponse{
		User:  user,
		Token: CreateTokenFromUser(user),
	}
	return c.JSON(resp)
}

func CreateTokenFromUser(user *types.User) string {
	now := time.Now()
	expires := now.Add(time.Hour * 5).Unix()
	claims := jwt.MapClaims{
		"id":      user.ID,
		"email":   user.Email,
		"expires": expires,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s := os.Getenv("JWT_SECRET")
	tokenStr, err := token.SignedString([]byte(s))
	if err != nil {
		fmt.Println("failed tp sign token with the secret")
	}
	return tokenStr

}
