package routes

import (
	"io"
	"log"
	"net/http"
	remerr "remindal/errors"
)

/*
Reads the body of an HTTP request and returns it as a byte array.

If there is an error reading the request body, it returns [nil] for the byte array, [remerr.ErrInternalServerError] as the error.
If the body is empty, it returns [nil] for the byte array, [remerr.ErrNoBodyProvided] as the error.
If the body is read successfully, it returns the body as a byte array and [nil] for the error.
*/
func decodeRequestBody(b io.Reader) ([]byte, error) {
	body, err := io.ReadAll(b)
	if err != nil {
		log.Println("decodeRequestBody - io.ReadAll ", err)
		return nil, remerr.ErrInternalServerError
	}
	if len(body) <= 0 {
		return nil, remerr.ErrNoBodyProvided
	}
	return body, nil
}

/*
Returns the appropriate HTTP status code based on the given error.
*/
func statusError(err error) int {
	var status int
	if err == remerr.ErrInternalServerError {
		status = http.StatusInternalServerError
	} else {
		status = http.StatusBadRequest
	}
	return status
}
