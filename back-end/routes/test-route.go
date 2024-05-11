package routes

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"remindal/database"
	"remindal/res"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	CTX = context.Background()

	EMAIL_KEY_DB   = "_id"
	EMAIL_KEY_JSON = "email"
)

// gets the list of users in the database
func GetUsersListHandler(w http.ResponseWriter, r *http.Request) {
	client, err := database.OpenConnection()
	if err != nil {
		res.Err(w, err, http.StatusBadRequest)
		return
	}
	defer database.CloseConnection(client)

	users := client.Database(database.DB_NAME).Collection(database.USER_COLLECTION)
	cursor, err := users.Find(CTX, bson.D{{}})
	if err != nil {
		res.Err(w, err, http.StatusInternalServerError)
		return
	}

	var retrievedUserList []database.UserSchema
	if err := cursor.All(CTX, &retrievedUserList); err != nil {
		// this should also be errNoDocumentsFound?
		res.Err(w, ErrInternalServerError, http.StatusInternalServerError)
		return
	}
	res.Ok(w, retrievedUserList)
}

// gets a single user from database
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	client, err := database.OpenConnection()
	if err != nil {
		res.Err(w, err, http.StatusInternalServerError)
		return
	}
	defer database.CloseConnection(client)

	email := r.URL.Query().Get(EMAIL_KEY_JSON)
	if email == "" {
		res.Err(w, ErrNoEmailProvided, http.StatusBadRequest)
		return
	}

	retrievedUser, err := getUser(client, email)
	if err != nil {
		res.Err(w, err, http.StatusInternalServerError)
		return
	}
	res.Ok(w, retrievedUser)
}

// handles the put request to add a new user to database
func PutUserHandler(w http.ResponseWriter, r *http.Request) {
	client, err := database.OpenConnection()
	if err != nil {
		res.Err(w, err, http.StatusInternalServerError)
		return
	}
	defer database.CloseConnection(client)

	body, statusCode, err := decodeRequestBody(r.Body)
	if err != nil {
		res.Err(w, err, statusCode)
		return
	}

	var newuser database.UserSchema
	if err := json.Unmarshal(body, &newuser); err != nil {
		log.Println("PutUserHandle - json.Unmarshal ", err)
		res.Err(w, ErrInternalServerError, http.StatusInternalServerError)
		return
	}

	if err := putNewUser(client, newuser); err != nil {
		res.Err(w, err, http.StatusBadRequest)
		return
	}
	res.Ok(w, nil)
}

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
		return nil, http.StatusInternalServerError, ErrInternalServerError
	}
	if len(body) <= 0 {
		return nil, http.StatusBadRequest, ErrNoBodyProvided
	}
	return body, http.StatusOK, nil
}

// helper - opens collection and gets a single user from database
func getUser(client *mongo.Client, name string) (database.UserSchema, error) {
	users := client.Database(database.DB_NAME).Collection(database.USER_COLLECTION)

	var retrieved database.UserSchema
	doc := users.FindOne(CTX, bson.D{{Key: EMAIL_KEY_DB, Value: name}})
	err := doc.Decode(&retrieved)
	if err == nil {
		return retrieved, nil
	}
	if err == mongo.ErrNoDocuments {
		return database.UserSchema{}, ErrNoDocumentsFound
	}
	return database.UserSchema{}, ErrInternalServerError
}

// helper - opens collection and puts a new user in the database
func putNewUser(client *mongo.Client, newuser database.UserSchema) error {
	users := client.Database(database.DB_NAME).Collection(database.USER_COLLECTION)

	// this is not a correct way to do this, must be changed.
	_, err := users.InsertOne(CTX, newuser)
	if err != nil {
		return ErrItemAlreadyPresent
	}
	return nil
}
