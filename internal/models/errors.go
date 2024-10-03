package models

import "errors"

var (
	// ErrInternalServer is a generic error message returned by a server
	// in case of an internal server error when we don't want to expose
	// the real error to the client. All the internal server errors
	// should be logged and fixed. Don't use this error if it's
	// something that can be fixed by the client.
	ErrInternalServer = errors.New("don't worry, we are working on it")

	// ErrUnauthorized is returned when the user credentials are invalid.
	ErrUnauthorized = errors.New("invalid credentials")

	// ErrUserNotFound is returned when the user is not found in the database.
	ErrUserNotFound = errors.New("user with specified id not found")
)
