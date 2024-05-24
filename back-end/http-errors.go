package main

import (
	"fmt"
)

type HttpError struct {
	err    error
	status int
}

func (he *HttpError) Error() string {
	return fmt.Sprintf("%d: %s", he.status, he.err)
}

func Err400(err error) *HttpError {
	return &HttpError{
		err:    err,
		status: 400,
	}
}

func Err403(err error) *HttpError {
	return &HttpError{
		err:    err,
		status: 403,
	}
}

func Err500(err error) *HttpError {
	return &HttpError{
		err:    err,
		status: 500,
	}
}
