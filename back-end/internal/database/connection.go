package database

import (
	"context"
	_ "embed"
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
)

// Closes connection to te Remindal database
func closeConnection(client *mongo.Client) {
	err := client.Disconnect(context.Background())
	if err != nil {
		log.Println("database.CloseConnection - client.Disconnect ", err)
	}
}

// Opens connection to the Remindal database.
//
// Should the attempt to establish a connection fail, returns [nil] for the client and [ErrInternalServerError] error.
func openConnection() (*mongo.Client, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		log.Println("database.OpenConnection - mongo.Connect ", err)
		return nil, err
	}
	return client, nil
}
