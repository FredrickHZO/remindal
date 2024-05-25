package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	db "remindal/internal/database"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	EMAIL_KEY_DB   = "_id"
	EMAIL_KEY_JSON = "email"

	errNoEmailProvided = errors.New("no email provided")
	errInvalidUserInfo = errors.New("invalid or missing user info")
)

/*
GetUsersListHandler handles requests to retrieve a list of users based on query parameters.

It converts the query parameters to a MongoDB query, retrieves the matching users from the database,
and writes the result as a JSON response. If an error occurs, it responds with the appropriate error message and status code.
*/
func GetUsersListHandler(w http.ResponseWriter, r *http.Request) {
	//query := r.URL.Query()

	var retrievedUserList []User
	sort := bson.D{{Key: "age", Value: 1}}
	query := bson.D{{}}

	err := db.GetMany(db.USER_COLLECTION, query, sort, &retrievedUserList)
	if err != nil {
		Eres(w, Err500(err))
		return
	}
	Okres(w, retrievedUserList)
}

/*
GetUserHandler handles requests to retrieve a single user based on their email.

It retrieves the email from the query parameters, fetches the user from the database,
and writes the result as a JSON response. If an error occurs, it responds with the appropriate error message and status code.
*/
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	userEmail := r.URL.Query().Get(EMAIL_KEY_JSON)
	if userEmail == "" {
		Eres(w, Err400(errNoEmailProvided))
		return
	}

	var retrievedUser User
	err := db.GetOne(db.USER_COLLECTION, EMAIL_KEY_DB,
		userEmail, &retrievedUser)
	if err != nil {
		Eres(w, Err400(err))
		return
	}
	Okres(w, retrievedUser)
}

/*
PutUserHandler handles requests to add a new user to the database.

It reads the request body, unmarshals the JSON into a UserSchema, and inserts the user into the database.
If an error occurs, it responds with the appropriate error message and status code.
*/
func PutUserHandler(w http.ResponseWriter, r *http.Request) {
	body, herr := decodeRequestBody(r.Body)
	if herr != nil {
		Eres(w, Err500(herr))
		return
	}

	var newuser User
	if err := json.Unmarshal(body, &newuser); err != nil {
		log.Println("PutUserHandle - json.Unmarshal ", err)
		Eres(w, Err500(err))
		return
	}

	validate := validator.New()
	err := validate.Struct(newuser)
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

/*
DelUserHandler handles requests to delete a user from the database based on their email.

It retrieves the email from the query parameters and deletes the user from the database.
If an error occurs, it responds with the appropriate error message and status code.
*/
func DelUserHandler(w http.ResponseWriter, r *http.Request) {
	userEmail := r.URL.Query().Get(EMAIL_KEY_JSON)
	if userEmail == "" {
		Eres(w, Err400(errNoEmailProvided))
		return
	}

	err := db.DeleteOne(db.USER_COLLECTION, EMAIL_KEY_DB, userEmail)
	if err != nil {
		Eres(w, Err400(err))
		return
	}
	Okres(w, nil)
}
