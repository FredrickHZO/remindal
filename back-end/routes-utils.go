package main

import (
	"io"
	"log"
)

/*
Reads the body of an HTTP request and returns it as a byte array.

If there is an error reading the request body, it returns [nil] for the byte array, [remerr.ErrInternalServerError] as the error.
If the body is empty, it returns [nil] for the byte array, [remerr.ErrNoBodyProvided] as the error.
If the body is read successfully, it returns the body as a byte array and [nil] for the error.
*/
func decodeRequestBody(b io.Reader) ([]byte, *HttpError) {
	body, err := io.ReadAll(b)
	if err != nil {
		log.Println("decodeRequestBody - io.ReadAll ", err)
		return nil, Err500(err)
	}
	return body, nil
}
