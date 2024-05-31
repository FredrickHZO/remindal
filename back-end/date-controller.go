package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	db "remindal/internal/database"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var errNoDateIDProvided = errors.New("no id provided for the date")

// Handles requests to retrieve a list of dates based on query parameters.
//
// Converts the query parameters to a MongoDB query, retrieves the matching
// users from the database and writes the result as a JSON response.
// If an error occurs, it responds with the appropriate error message and status code.
func GetDateListHandler(w http.ResponseWriter, r *http.Request) {
	var (
		builder = db.NewQueryBuilder()
		query   = r.URL.Query()
	)
	buildDateQuery(query, &builder)
	err := builder.Err()
	if err != nil {
		Eres(w, Err400(err))
		return
	}

	sort := db.CreateSort("year", -1)
	d := []Date{}
	err = db.GetMany(db.CALENDAR_COLLECTION, builder.Query(), sort, &d)
	if err != nil {
		log.Println("GetDateListHandler - db.GetMany ", err)
		Eres(w, Err500(err))
		return
	}
	Okres(w, d)
}

// Handles requests to delete a date from the database based on its id.
//
// Retrieves the id from the query parameters and deletes the date from the database.
// If an error occurs, it responds with the appropriate error message and status code.
func DelDateHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("_id")
	if id == "" {
		Eres(w, Err400(errNoDateIDProvided))
		return
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("DelDateHandler - primitive.ObjectIDFromHex ", err)
		Eres(w, Err500(err))
		return
	}

	err = db.DeleteOne(db.CALENDAR_COLLECTION, "_id", objID)
	if err != nil {
		Eres(w, Err400(err))
		return
	}
	Okres(w, nil)
}

// Handles requests to add a new date to the database.
//
// Reads the request body, unmarshals the JSON into a Date, and inserts the date
// into the database. If an error occurs, it responds with the appropriate error
// message and status code.
func PutDateHandler(w http.ResponseWriter, r *http.Request) {
	jsn, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("PutCalendarHandler - io.ReadAll ", err)
		Eres(w, Err500(err))
		return
	}

	var d Date
	err = json.Unmarshal(jsn, &d)
	if err != nil {
		log.Println("PutCalendarHandler - json.Unmarshal ", err)
		Eres(w, Err500(err))
		return
	}

	validate := newCustomDateValidator()
	err = validate.Struct(d)
	if err != nil {
		Eres(w, Err400(err))
		return
	}

	err = db.PutOne(db.CALENDAR_COLLECTION, d)
	if err != nil {
		log.Println("PutCalendarHandler - db.PutOne ", err)
		Eres(w, Err500(err))
		return
	}
	Okres(w, nil)
}
