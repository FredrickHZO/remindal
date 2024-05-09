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
	mongoURI        string
	DB_NAME         = "calendars"
	USER_COLLECTION = "users"
)

func CloseConnection(client *mongo.Client) {
	err := client.Disconnect(context.Background())
	if err != nil {
		log.Fatal("database - CloseConnection ", err)
	}
}

func OpenConnection() *mongo.Client {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		log.Fatal("database - OpenConnection ", err)
	}
	return client
}
