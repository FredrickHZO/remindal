package routes

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	db "remindal/database"
	remerr "remindal/errors"
	"remindal/res"
)

var (
	CTX = context.Background()

	EMAIL_KEY_DB   = "_id"
	EMAIL_KEY_JSON = "email"
)

func GetUsersListHandler(w http.ResponseWriter, r *http.Request) {
	client, err := db.OpenConnection()
	if err != nil {
		res.Err(w, err, http.StatusInternalServerError)
		return
	}
	defer db.CloseConnection(client)

	var retrievedUserList []db.UserSchema
	err = db.GetMany(client, db.USER_COLLECTION, &retrievedUserList)
	if err != nil {
		res.Err(w, err, http.StatusInternalServerError)
	}
	res.Ok(w, retrievedUserList)
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	client, err := db.OpenConnection()
	if err != nil {
		res.Err(w, err, http.StatusInternalServerError)
		return
	}
	defer db.CloseConnection(client)

	userEmail := r.URL.Query().Get(EMAIL_KEY_JSON)
	if userEmail == "" {
		res.Err(w, remerr.ErrNoEmailProvided, http.StatusBadRequest)
		return
	}

	var retrievedUser db.UserSchema
	err = db.GetOne(
		client,
		db.USER_COLLECTION,
		EMAIL_KEY_DB,
		userEmail,
		&retrievedUser,
	)
	if err != nil {
		res.Err(w, err, http.StatusInternalServerError)
		return
	}
	res.Ok(w, retrievedUser)
}

func PutUserHandler(w http.ResponseWriter, r *http.Request) {
	client, err := db.OpenConnection()
	if err != nil {
		res.Err(w, err, http.StatusInternalServerError)
		return
	}
	defer db.CloseConnection(client)

	body, statusCode, err := decodeRequestBody(r.Body)
	if err != nil {
		res.Err(w, err, statusCode)
		return
	}

	var newuser db.UserSchema
	if err := json.Unmarshal(body, &newuser); err != nil {
		log.Println("PutUserHandle - json.Unmarshal ", err)
		res.Err(w, remerr.ErrInternalServerError, http.StatusInternalServerError)
		return
	}

	if err := db.PutOne(client, db.USER_COLLECTION, newuser); err != nil {
		res.Err(w, err, http.StatusBadRequest)
		return
	}
	res.Ok(w, nil)
}

func DelUserHandler(w http.ResponseWriter, r *http.Request) {
	client, err := db.OpenConnection()
	if err != nil {
		res.Err(w, err, http.StatusInternalServerError)
	}
	defer db.CloseConnection(client)

	userEmail := r.URL.Query().Get(EMAIL_KEY_JSON)
	if userEmail == "" {
		res.Err(w, remerr.ErrNoEmailProvided, http.StatusBadRequest)
		return
	}

	err = db.DeleteOne(client, db.USER_COLLECTION, EMAIL_KEY_DB, userEmail)
	if err != nil {
		var status int
		if err == remerr.ErrInternalServerError {
			status = http.StatusInternalServerError
		} else {
			status = http.StatusBadRequest
		}
		res.Err(w, err, status)
		return
	}
	res.Ok(w, nil)
}
