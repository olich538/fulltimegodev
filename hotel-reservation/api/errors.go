package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {

	if apiError, ok := err.(Error); ok {
		return c.Status(apiError.Code).JSON(apiError)
	}
	apiErr := NewError(http.StatusInternalServerError, err.Error())
	return c.Status(apiErr.Code).JSON(apiErr)

}

type Error struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}

// error implement the error interface
func (e Error) Error() string {
	return e.Err
}
func NewError(code int, err string) Error {
	return Error{
		Code: code,
		Err:  err,
	}
}
func ErrInvalidID() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err:  "invalid ID",
	}
}
func ErrBadRequest() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err:  "invalid JSON request",
	}
}
func ErrNotFound(res string) Error {
	return Error{
		Code: http.StatusNotFound,
		Err:  res + " resource not found",
	}
}

func ErrUnauthorized() Error {
	return Error{
		Code: http.StatusUnauthorized,
		Err:  "unauthrized request",
	}
}
