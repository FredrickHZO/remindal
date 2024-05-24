package main

import (
	"io"
	"log"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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

/*
Reimplement
*/
func rangeFilter(k string, v []string, cond string) (bson.E, error) {
	val, err := strconv.Atoi(v[0])
	if err != nil {
		return bson.E{}, err
	}

	rangeCond := bson.E{
		Key:   k,
		Value: bson.D{{Key: cond, Value: val}},
	}
	log.Println(rangeCond)
	return rangeCond, nil
}

/*
Maybe needed
*/
func stringToDate(s string) (time.Time, error) {
	layout := "2021-11-22"
	t, err := time.Parse(layout, s)
	if err != nil {
		return time.Time{}, err
	}
	return t, err
}
