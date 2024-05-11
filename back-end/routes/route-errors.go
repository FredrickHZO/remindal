package routes

import "errors"

var (
	ErrInternalServerError = errors.New("internal server error")

	// database response errors
	ErrItemAlreadyPresent = errors.New("item already in database")
	ErrNoDocumentsFound   = errors.New("no documents found")

	// request parsing errors
	ErrNoBodyProvided = errors.New("no body for request provided")

	// user route specific errors
	ErrNoEmailProvided = errors.New("no email provided")
)
