package server

import (
	"fmt"
	"net/http"
	"net/url"
)

type InvalidQuery string

func (e InvalidQuery) Error() string {
	return fmt.Sprintf("invalid %s", string(e))
}

// InvalidBody TODO
type InvalidBody string

func (e InvalidBody) Error() string {
	return fmt.Sprintf("invalid %s", string(e))
}

type NotExistError string

func (e NotExistError) Error() string {
	return fmt.Sprintf("not exist %s", string(e))
}

// ExistError TODO
type ExistError string

func (e ExistError) Error() string {
	return fmt.Sprintf("already exist %s", string(e))
}

// UnavailableError TODO
type UnavailableError string

func (e UnavailableError) Error() string {
	return fmt.Sprintf("unavailable %s", string(e))
}

func ErrorCode(err error) int {
	var code int
	switch err.(type) {
	case *url.Error, InvalidQuery, InvalidBody:
		code = http.StatusBadRequest
	case NotExistError:
		code = http.StatusNotFound
	case UnavailableError:
		code = http.StatusNotAcceptable
	default:
		code = http.StatusInternalServerError
	}
	return code
}