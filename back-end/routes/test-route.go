package routes

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"remindal/database"
	"remindal/res"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	CTX = context.Background()

	errInternalServerError = errors.New("internal server error")
	errItemAlreadyPresent  = errors.New("item already in database")
	errNoDocumentsFound    = errors.New("no documents found")
	errNoNameProvided      = errors.New("no name provided")
	errNoBodyProvided      = errors.New("no body for request provided")
)

// gets the list of users in the database
func GetUsersListHandle(w http.ResponseWriter, r *http.Request) {
	client := database.OpenConnection()
	defer database.CloseConnection(client)

	users := client.Database(database.DB_NAME).Collection(database.USER_COLLECTION)
	cursor, err := users.Find(CTX, bson.D{{}})
	if err != nil {
		res.Err(w, err, 400)
		return
	}

	var retrievedUserList []database.UserSchema
	if err := cursor.All(CTX, &retrievedUserList); err != nil {
		// this should also be errNoDocumentsFound?
		res.Err(w, errInternalServerError, 500)
		return
	}
	res.Ok(w, retrievedUserList)
}

// gets a single user from database
func GetUserHandle(w http.ResponseWriter, r *http.Request) {
	client := database.OpenConnection()
	defer database.CloseConnection(client)

	name := r.URL.Query().Get("name")
	if name == "" {
		res.Err(w, errNoNameProvided, 400)
		return
	}

	retrievedUser, err := getUser(client, name)
	// not necessary? see helper function comment
	if err != nil {
		res.Err(w, errNoDocumentsFound, 400)
		return
	}
	res.Ok(w, retrievedUser)
}

// decodes the request body
func decodeRequestBody(b io.Reader) ([]byte, error) {
	body, err := io.ReadAll(b)
	if err != nil {
		return nil, err
	}
	if len(body) <= 0 {
		return nil, errNoBodyProvided
	}
	return body, nil
}

// handles the put request to add a new user to database
func PutUserHandle(w http.ResponseWriter, r *http.Request) {
	client := database.OpenConnection()
	defer database.CloseConnection(client)

	body, err := decodeRequestBody(r.Body)
	if err != nil {
		res.Err(w, err, 400)
		return
	}

	var newuser database.UserSchema
	if err := json.Unmarshal(body, &newuser); err != nil {
		res.Err(w, errInternalServerError, 500)
		return
	}

	if err := putNewUser(client, newuser); err != nil {
		res.Err(w, err, 500)
		return
	}
	res.Ok(w, nil)
}

// helper - opens collection and gets a single user from database
func getUser(client *mongo.Client, name string) (database.UserSchema, error) {
	users := client.Database(database.DB_NAME).Collection(database.USER_COLLECTION)

	var retrieved database.UserSchema
	doc := users.FindOne(CTX, bson.D{{Key: "name", Value: name}})
	err := doc.Decode(&retrieved)
	return retrieved, err
}

// helper - opens collection and puts a new user in the database
func putNewUser(client *mongo.Client, newuser database.UserSchema) error {
	users := client.Database(database.DB_NAME).Collection(database.USER_COLLECTION)

	// this is not a correct way to do this, must be changed.
	_, err := users.InsertOne(CTX, newuser)
	if err != nil {
		return errItemAlreadyPresent
	}
	return nil
}
