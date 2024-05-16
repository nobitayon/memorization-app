package model

import (
	"errors"
	"fmt"
	"net/http"
)

type CodeError string

const (
	Authorization   CodeError = "AUTHORIZATION"
	BadRequest      CodeError = "BADREQUEST"
	Conflict        CodeError = "CONFLICT"
	Internal        CodeError = "INTERNAL"
	NotFound        CodeError = "NOTFOUND"
	PayloadTooLarge CodeError = "PAYLOADTOOLARGE"
)

type Error struct {
	Type    CodeError `json:"type"`
	Message string    `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Status() int {
	switch e.Type {
	case Authorization:
		return http.StatusUnauthorized
	case BadRequest:
		return http.StatusBadRequest
	case Conflict:
		return http.StatusConflict
	case Internal:
		return http.StatusInternalServerError
	case NotFound:
		return http.StatusNotFound
	case PayloadTooLarge:
		return http.StatusRequestEntityTooLarge
	default:
		return http.StatusInternalServerError
	}
}

func Status(err error) int {
	var e *Error
	if errors.As(err, &e) {
		return e.Status()
	}
	return http.StatusInternalServerError
}

// error factory
func NewAuthorization(reason string) *Error {
	return &Error{
		Type:    Authorization,
		Message: reason,
	}
}

func NewBadRequest(reason string) *Error {
	return &Error{
		Type:    BadRequest,
		Message: fmt.Sprintf("Bad request. Reason: %v", reason),
	}
}

func NewConflict(name string, value string) *Error {
	return &Error{
		Type:    Conflict,
		Message: fmt.Sprintf("resource: %v with value: %v already exists", name, value),
	}
}

func NewInternal(name string, value string) *Error {
	return &Error{
		Type:    Internal,
		Message: "internal server error",
	}
}

func NewNotFound(name string, value string) *Error {
	return &Error{
		Type:    NotFound,
		Message: fmt.Sprintf("resource: %v with value: %v not found", name, value),
	}
}

func NewPayloadTooLarge(maxBodySize int64, contentLength int64) *Error {
	return &Error{
		Type:    NotFound,
		Message: fmt.Sprintf("Max payload size of %v extended. Actual payload size: %v", maxBodySize, contentLength),
	}
}
