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

/*
Opens a connection to the database and retrieves an array of items that match the provided query.
Fetches multiple documents based on the specified query filter and unmarshals the results into the provided destination.

Usage:

	var destination []mySchema
	err := GetMany(myCollection, bson.D{{Key: "_id", Value: "email@person.com"}}, &destination)

@param collectionName: The name of the MongoDB collection to query.
@param query: The BSON query filter to apply.
@param dest: A pointer to the variable where the results will be stored.
@return error: An error object if any error occurs during the operation.

Possible errors:
- [ErrInternalServerError]: If a connection to the database cannot be established or if the retrieval operation fails.
- [ErrNoDocumentsFound]: If no documents match the query.
*/
func GetMany(collectionName string, query bson.D, dest any) error {
	client, err := openConnection()
	if err != nil {
		return remerr.ErrInternalServerError
	}
	defer closeConnection(client)

	coll := client.Database(DB_NAME).Collection(collectionName)
	cursor, err := coll.Find(CTX, query)
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

@param collectionName: The name of the MongoDB collection to query.
@param key: The key to search for in the database (e.g., "_id").
@param value: The value of the key to match in the query.
@param dest: A pointer to the variable where the result will be stored.
@return error: An error object if any error occurs during the operation.

Possible errors:
- [ErrInternalServerError]: If a connection to the database cannot be established or if the retrieval operation fails.
- [ErrNoDocumentsFound]: If no document matches the key-value pair.
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

Usage:

	err := PutOne(myCollection, myDocument)

@param collectionName: The name of the MongoDB collection to insert the document into.
@param doc: The document to be inserted.
@return error: An error object if any error occurs during the operation.

Possible errors:
- [ErrInternalServerError]: If a connection to the database cannot be established.
- [ErrItemAlreadyPresent]: If there is a collision with the primary key of an existing item in the database.
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

Usage:

	err := DeleteOne(myCollection, "_id", "email@person.com")

@param collectionName: The name of the MongoDB collection to delete the document from.
@param key: The key to search for in the database (e.g., "_id").
@param value: The value of the key to match in the query.
@return error: An error object if any error occurs during the operation.

Possible errors:
- [ErrInternalServerError]: If a connection to the database cannot be established or if the delete operation fails.
- [ErrNoDocumentsFound]: If no document matches the key-value pair.
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
