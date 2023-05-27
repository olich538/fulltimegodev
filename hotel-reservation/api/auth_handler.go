package api

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/olich538/fulltimegodev/hotel-reservation/db"
	"github.com/olich538/fulltimegodev/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
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

// a handler should only do:
// - serialization of the incomming requests
// - do some data fetching from DB
// - return the data back to user
// - call some bussines logic

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var authParams AuthParams

	if err := c.BodyParser(&authParams); err != nil {
		return err
	}

	user, err := h.userStore.GetUserByEmail(c.Context(), authParams.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("invalid credentials")
		}
		return err
	}
	if !types.IsValidPassword(user.EncryptedPassword, authParams.Password) {
		return fmt.Errorf("invalid password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(authParams.Password))
	if err != nil {
		return fmt.Errorf("invalid credentials")
	}
	token := createTokenFromUser(user)
	fmt.Println("authenticated ->", user)
	resp := AuthResponse{
		User:  user,
		Token: token,
	}
	return c.JSON(resp)
}

func createTokenFromUser(user *types.User) string {
	now := time.Now()
	validTill := now.Add(time.Hour * 5)
	claims := jwt.MapClaims{
		"id":        user.ID,
		"email":     user.Email,
		"validTill": validTill,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s := os.Getenv("JWT_SECRET")
	tokenStr, err := token.SignedString(s)
	if err != nil {
		fmt.Println("failed tp sign token with the secret")
	}
	return tokenStr

}
