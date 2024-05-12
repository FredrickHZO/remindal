package routes

import (
	"io"
	"log"
	"net/http"
	remerr "remindal/errors"
)

/*
Decodes the body in the HTTP request and returns it as a byte array.

If there is an error trying to read the request body, returns [nil] for the byte array, [http.StatusInternalServerError]
as code status and [ErrInternalServerError] as error.

Should the body be empty, returns [nil] for the byte array, [http.StatusBadRequest] as code status and
[ErrNoBodyProvided] as error.

In case the body is correctly read it will return the body as a byte array, [http.StatusOK] as code status
and [nil] for the error.
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
Given an error, it returns the correct HTTP error code status.
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
