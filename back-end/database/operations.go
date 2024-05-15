package database

import (
	"context"
	remerr "remindal/errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	CTX = context.Background()
)

/*
Opens a connection to the database and retrieves an array of items that match the provided query.
Fetches multiple documents based on the specified query filter and unmarshals the results into the provided destination.

Usage:

	var destination []mySchema
	err := GetMany(myCollection, bson.D{{Key: "_id", Value: "email@person.com"}}, &destination)

[ErrInternalServerError]: If a connection to the database cannot be established or if the retrieval operation fails.
[ErrNoDocumentsFound]: If no documents match the query.
*/
func GetMany(collectionName string, query bson.D, sort bson.D, dest any) error {
	client, err := openConnection()
	if err != nil {
		return remerr.ErrInternalServerError
	}
	defer closeConnection(client)

	opts := options.Find().SetSort(sort)
	coll := client.Database(DB_NAME).Collection(collectionName)
	cursor, err := coll.Find(CTX, query, opts)
	if err != nil {
		return remerr.ErrInternalServerError
	}
	if err := cursor.All(CTX, dest); err != nil {
		return remerr.ErrInternalServerError
	}
	return nil
}

/*
Opens a connection to the database and retrieves a single document that matches the provided key-value pair.
Fetches a document based on the specified key and value and unmarshals the result into the provided destination.

Usage:

	var destination mySchema
	err := GetOne(myCollection, "_id", "email@person.com", &destination)

[ErrInternalServerError]: If a connection to the database cannot be established or if the retrieval operation fails.
[ErrNoDocumentsFound]: If no document matches the key-value pair.
*/
func GetOne(collectionName string, key string, value string, dest any) error {
	client, err := openConnection()
	if err != nil {
		return remerr.ErrInternalServerError
	}
	defer closeConnection(client)

	coll := client.Database(DB_NAME).Collection(collectionName)
	doc := coll.FindOne(CTX, bson.D{{Key: key, Value: value}})
	err = doc.Decode(dest)
	if err == nil {
		return nil
	}
	if err == mongo.ErrNoDocuments {
		return remerr.ErrNoDocumentsFound
	}
	return remerr.ErrInternalServerError
}

/*
Opens a connection to the database and inserts the provided document into the specified collection.

[ErrInternalServerError]: If a connection to the database cannot be established.
[ErrItemAlreadyPresent]: If there is a collision with the primary key of an existing item in the database.
*/
func PutOne(collectionName string, doc any) error {
	client, err := openConnection()
	if err != nil {
		return remerr.ErrInternalServerError
	}
	defer closeConnection(client)

	coll := client.Database(DB_NAME).Collection(collectionName)
	// this is not a correct way to do this, must be changed.
	_, err = coll.InsertOne(CTX, doc)
	if err != nil {
		return remerr.ErrItemAlreadyPresent
	}
	return nil
}

/*
Opens a connection to the database and deletes a document that matches the provided key-value pair.

[ErrInternalServerError]: If a connection to the database cannot be established or if the delete operation fails.
[ErrNoDocumentsFound]: If no document matches the key-value pair.
*/
func DeleteOne(collectionName string, key string, value string) error {
	client, err := openConnection()
	if err != nil {
		return remerr.ErrInternalServerError
	}
	defer closeConnection(client)

	coll := client.Database(DB_NAME).Collection(collectionName)
	delres, err := coll.DeleteOne(CTX, bson.D{{Key: key, Value: value}})
	if delres.DeletedCount == 0 {
		return remerr.ErrNoDocumentsFound
	}
	if err != nil {
		return remerr.ErrInternalServerError
	}
	return nil
}
