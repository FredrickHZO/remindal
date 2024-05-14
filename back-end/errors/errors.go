package errors

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

	// HTTP request errors
	ErrTooManyParametersInRange = errors.New("too many parameters in range")
	ErrRangeValueNotNumber      = errors.New("non numeric value for range")
	ErrInvalidRangeValues       = errors.New("invalid range values")
	ErrNotRangeable             = errors.New("this field cannot be used as a range")
	ErrNotMultipleSelection     = errors.New("this field cannot be used as a multiple selection filter")
)
