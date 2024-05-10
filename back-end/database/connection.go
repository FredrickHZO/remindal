package database

import (
	"context"
	_ "embed"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	//go:embed mongouri
	mongoURI            string
	DB_NAME             = "remindalDB"
	CALENDAR_COLLECTION = "calendar"
	USER_COLLECTION     = "users"

	errCouldNotReachDatabase = errors.New("could not reach database")
)

func CloseConnection(client *mongo.Client) {
	err := client.Disconnect(context.Background())
	if err != nil {
		log.Println("database - CloseConnection ", err)
	}
}

/*
Opens connection to the Remindal database.

Should the attempt to establish a connection fail, returns [nil] for the client and
[errCouldNotReachDatabase] error.
*/
func OpenConnection() (*mongo.Client, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		log.Println("database - OpenConnection ", err)
		return nil, errCouldNotReachDatabase
	}
	return client, nil
}
