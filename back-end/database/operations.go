package database

import (
	"context"
	remerr "remindal/errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	CTX = context.Background()
)

func GetMany(client *mongo.Client, collection string, dest any) error {
	c := client.Database(DB_NAME).Collection(collection)
	// TODO: must work with any filter applied
	cursor, err := c.Find(CTX, bson.D{{}})
	if err != nil {
		return remerr.ErrInternalServerError
	}

	if err := cursor.All(CTX, dest); err != nil {
		return remerr.ErrInternalServerError
	}
	return nil
}

// helper - opens collection and gets a single user from database
func GetOne(client *mongo.Client, collection string, key string, value string, dest any) error {
	c := client.Database(DB_NAME).Collection(collection)

	doc := c.FindOne(CTX, bson.D{{Key: key, Value: value}})
	err := doc.Decode(dest)
	if err == nil {
		return nil
	}
	if err == mongo.ErrNoDocuments {
		return remerr.ErrNoDocumentsFound
	}
	return remerr.ErrInternalServerError
}

// helper - opens collection and puts a new user in the database
func PutOne(client *mongo.Client, collection string, doc any) error {
	c := client.Database(DB_NAME).Collection(collection)

	// this is not a correct way to do this, must be changed.
	_, err := c.InsertOne(CTX, doc)
	if err != nil {
		return remerr.ErrItemAlreadyPresent
	}
	return nil
}

func DeleteOne(client *mongo.Client, collection string, key string, value string) error {
	c := client.Database(DB_NAME).Collection(collection)

	delres, err := c.DeleteOne(CTX, bson.D{{Key: key, Value: value}})
	if delres.DeletedCount == 0 {
		return remerr.ErrNoItemToDelete
	}
	if err != nil {
		return remerr.ErrInternalServerError
	}
	return nil
}
