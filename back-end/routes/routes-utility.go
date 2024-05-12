package routes

import (
	"io"
	"log"
	"net/http"
	remerr "remindal/errors"
)

/*
Decodes the body in the HTTP request and returns it as a byte array.

If there is an error trying to read the request body, returns [nil] for the byte array, [InternalServerError]
as code status and [errInternalServerError] as error.

Should the body be empty, returns [nil] for the byte array, [StatusBadRequest] as code status and
[errNoBodyProvided] as error.

In case the body is correctly read it will return the body as a byte array, [StatusOK] as code status
and [nil] for the error.

TODO: move function.
*/
func decodeRequestBody(b io.Reader) ([]byte, int, error) {
	body, err := io.ReadAll(b)
	if err != nil {
		log.Println("decodeRequestBody - io.ReadAll ", err)
		return nil, http.StatusInternalServerError, remerr.ErrInternalServerError
	}
	if len(body) <= 0 {
		return nil, http.StatusBadRequest, remerr.ErrNoBodyProvided
	}
	return body, http.StatusOK, nil
}
