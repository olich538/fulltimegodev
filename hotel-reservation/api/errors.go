package api

import "net/http"

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

func ErrUnauthorized() Error {
	return Error{
		Code: http.StatusUnauthorized,
		Err:  "unauthrized request",
	}
}
