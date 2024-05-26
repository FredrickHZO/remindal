package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	db "remindal/internal/database"

	"github.com/go-playground/validator/v10"
)

var (
	EMAIL_KEY = "_id"

	errNoEmailProvided = errors.New("no email provided")
	errInvalidUserInfo = errors.New("invalid or missing user info")
)

// Handles requests to retrieve a list of users based on query parameters.
//
// Converts the query parameters to a MongoDB query, retrieves the matching users from the database
// and writes the result as a JSON response. If an error occurs, it responds with the appropriate error message and status code.
func GetUsersListHandler(w http.ResponseWriter, r *http.Request) {
	var (
		qbuilder = db.NewQueryBuilder()
		query    = r.URL.Query()
	)
	constructUserQuery(query, &qbuilder)
	err := qbuilder.Err()
	if err != nil {
		Eres(w, Err400(err))
		return
	}

	var retrievedUserList []User
	sort := db.CreateSort("age", 1)
	err = db.GetMany(db.USER_COLLECTION, qbuilder.Query(), sort, &retrievedUserList)
	if err != nil {
		Eres(w, Err500(err))
		return
	}
	Okres(w, retrievedUserList)
}

// GetUserHandler handles requests to retrieve a single user based on their email.
//
// Retrieves the email from the query parameters, fetches the user from the database
// and writes the result as a JSON response. If an error occurs, it responds with the appropriate error message and status code.
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	userEmail := r.URL.Query().Get(EMAIL_KEY)
	if userEmail == "" {
		Eres(w, Err400(errNoEmailProvided))
		return
	}

	var retrievedUser User
	err := db.GetOne(db.USER_COLLECTION, EMAIL_KEY, userEmail, &retrievedUser)
	if err != nil {
		Eres(w, Err400(err))
		return
	}
	Okres(w, retrievedUser)
}

// Handles requests to add a new user to the database.
//
// Reads the request body, unmarshals the JSON into a UserSchema, and inserts the user into the database.
// If an error occurs, it responds with the appropriate error message and status code.
func PutUserHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("decodeRequestBody - io.ReadAll ", err)
		Eres(w, Err500(err))
		return
	}

	var newuser User
	if err := json.Unmarshal(body, &newuser); err != nil {
		log.Println("PutUserHandle - json.Unmarshal ", err)
		Eres(w, Err500(err))
		return
	}

	validate := validator.New()
	err = validate.Struct(newuser)
	if err != nil {
		Eres(w, Err400(errInvalidUserInfo))
		return
	}

	if err := db.PutOne(db.USER_COLLECTION, newuser); err != nil {
		Eres(w, Err500(err))
		return
	}
	Okres(w, nil)
}

// Handles requests to delete a user from the database based on their email.
//
// Retrieves the email from the query parameters and deletes the user from the database.
// If an error occurs, it responds with the appropriate error message and status code.
func DelUserHandler(w http.ResponseWriter, r *http.Request) {
	userEmail := r.URL.Query().Get(EMAIL_KEY)
	if userEmail == "" {
		Eres(w, Err400(errNoEmailProvided))
		return
	}

	err := db.DeleteOne(db.USER_COLLECTION, EMAIL_KEY, userEmail)
	if err != nil {
		Eres(w, Err400(err))
		return
	}
	Okres(w, nil)
}
