package apperr

import (
	"fmt"
	"net/http"
)

type Error struct {
	Status  int
	Message string
	kind    *Error
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) WithMessage(format string, args ...any) *Error {
	return &Error{
		Status:  e.Status,
		Message: fmt.Sprintf(format, args...),
		kind:    e.sentinel(),
	}
}

func (e *Error) Wrap(err error) error {
	if err == nil {
		return e
	}
	return fmt.Errorf("%w: %v", e, err)
}

func (e *Error) sentinel() *Error {
	if e.kind != nil {
		return e.kind
	}
	return e
}

func New(status int, message string) *Error {
	return &Error{Status: status, Message: message}
}

var (
	ErrAirportNotFound      = New(http.StatusNotFound, "airport not found")
	ErrAirportAlreadyExists = New(http.StatusConflict, "airport with this iataCode already exists")
	ErrInternal             = New(http.StatusInternalServerError, "internal server error")

	ErrMissingIATACode = New(http.StatusBadRequest, "missing airport iataCode")
	ErrInvalidIATACode = New(http.StatusBadRequest, "invalid iataCode format")
	ErrInvalidCountry  = New(http.StatusBadRequest, "invalid country format")
	ErrInvalidName     = New(http.StatusBadRequest, "invalid name format")

	ErrDecodeJSON = New(http.StatusBadRequest, "failed to decode request")
	ErrEncodeJSON = New(http.StatusInternalServerError, "failed to encode response")
)
