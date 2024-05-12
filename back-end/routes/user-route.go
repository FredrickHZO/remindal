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
	// TODO: must work with any filter applied
	var retrievedUserList []db.UserSchema
	err := db.GetMany(db.USER_COLLECTION, &retrievedUserList)
	if err != nil {
		res.Err(w, err, http.StatusInternalServerError)
	}
	res.Ok(w, retrievedUserList)
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	userEmail := r.URL.Query().Get(EMAIL_KEY_JSON)
	if userEmail == "" {
		res.Err(w, remerr.ErrNoEmailProvided, http.StatusBadRequest)
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
		status := statusError(err)
		res.Err(w, err, status)
		return
	}
	res.Ok(w, retrievedUser)
}

func PutUserHandler(w http.ResponseWriter, r *http.Request) {
	body, err := decodeRequestBody(r.Body)
	if err != nil {
		status := statusError(err)
		res.Err(w, err, status)
		return
	}

	var newuser db.UserSchema
	if err := json.Unmarshal(body, &newuser); err != nil {
		log.Println("PutUserHandle - json.Unmarshal ", err)
		res.Err(w, remerr.ErrInternalServerError, http.StatusInternalServerError)
		return
	}

	if err := db.PutOne(db.USER_COLLECTION, newuser); err != nil {
		res.Err(w, err, http.StatusBadRequest)
		return
	}
	res.Ok(w, nil)
}

func DelUserHandler(w http.ResponseWriter, r *http.Request) {
	userEmail := r.URL.Query().Get(EMAIL_KEY_JSON)
	if userEmail == "" {
		res.Err(w, remerr.ErrNoEmailProvided, http.StatusBadRequest)
		return
	}

	err := db.DeleteOne(db.USER_COLLECTION, EMAIL_KEY_DB, userEmail)
	if err != nil {
		status := statusError(err)
		res.Err(w, err, status)
		return
	}
	res.Ok(w, nil)
}
