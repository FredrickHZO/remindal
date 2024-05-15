package routes

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	db "remindal/database"
	remerr "remindal/errors"
	"remindal/res"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	CTX = context.Background()

	EMAIL_KEY_DB   = "_id"
	EMAIL_KEY_JSON = "email"
)

/*
GetUsersListHandler handles requests to retrieve a list of users based on query parameters.

It converts the query parameters to a MongoDB query, retrieves the matching users from the database,
and writes the result as a JSON response. If an error occurs, it responds with the appropriate error message and status code.
*/
func GetUsersListHandler(w http.ResponseWriter, r *http.Request) {
	query, err := db.ToMongoQuery(r.URL.Query())
	if err != nil {
		res.Err(w, err)
		return
	}

	var retrievedUserList []db.UserSchema
	sort := bson.D{{Key: "age", Value: 1}}
	err = db.GetMany(db.USER_COLLECTION, query, sort, &retrievedUserList)
	if err != nil {
		res.Err(w, err)
		return
	}
	res.Ok(w, retrievedUserList)
}

/*
GetUserHandler handles requests to retrieve a single user based on their email.

It retrieves the email from the query parameters, fetches the user from the database,
and writes the result as a JSON response. If an error occurs, it responds with the appropriate error message and status code.
*/
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	userEmail := r.URL.Query().Get(EMAIL_KEY_JSON)
	if userEmail == "" {
		res.Err(w, remerr.ErrNoEmailProvided)
		return
	}

	var retrievedUser db.UserSchema
	err := db.GetOne(
		db.USER_COLLECTION,
		EMAIL_KEY_DB,
		userEmail,
		&retrievedUser,
	)
	if err != nil {
		res.Err(w, err)
		return
	}
	res.Ok(w, retrievedUser)
}

/*
PutUserHandler handles requests to add a new user to the database.

It reads the request body, unmarshals the JSON into a UserSchema, and inserts the user into the database.
If an error occurs, it responds with the appropriate error message and status code.
*/
func PutUserHandler(w http.ResponseWriter, r *http.Request) {
	body, err := decodeRequestBody(r.Body)
	if err != nil {
		res.Err(w, err)
		return
	}

	var newuser db.UserSchema
	if err := json.Unmarshal(body, &newuser); err != nil {
		log.Println("PutUserHandle - json.Unmarshal ", err)
		res.Err(w, remerr.ErrInternalServerError)
		return
	}

	if err := db.PutOne(db.USER_COLLECTION, newuser); err != nil {
		res.Err(w, err)
		return
	}
	res.Ok(w, nil)
}

/*
DelUserHandler handles requests to delete a user from the database based on their email.

It retrieves the email from the query parameters and deletes the user from the database.
If an error occurs, it responds with the appropriate error message and status code.
*/
func DelUserHandler(w http.ResponseWriter, r *http.Request) {
	userEmail := r.URL.Query().Get(EMAIL_KEY_JSON)
	if userEmail == "" {
		res.Err(w, remerr.ErrNoEmailProvided)
		return
	}

	err := db.DeleteOne(db.USER_COLLECTION, EMAIL_KEY_DB, userEmail)
	if err != nil {
		res.Err(w, err)
		return
	}
	res.Ok(w, nil)
}
