package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Opens a connection to the database and retrieves an array of items that match the provided query.
// Fetches multiple documents based on the specified query filter and unmarshals the results into the provided destination.
//
// [ErrInternalServerError]: If a connection to the database cannot be established or if the retrieval operation fails.
// [ErrNoDocumentsFound]: If no documents match the query.
func GetMany(collectionName string, query bson.D, sort bson.D, dest any) error {
	client, err := openConnection()
	if err != nil {
		return err
	}
	defer closeConnection(client)

	opts := options.Find().SetSort(sort)
	coll := client.Database(DB_NAME).Collection(collectionName)
	cursor, err := coll.Find(context.TODO(), query, opts)
	if err != nil {
		return err
	}
	if err := cursor.All(context.TODO(), dest); err != nil {
		return err
	}
	return nil
}

// Opens a connection to the database and retrieves a single document that matches the provided key-value pair.
// Fetches a document based on the specified key and value and unmarshals the result into the provided destination.
//
// [ErrInternalServerError]: If a connection to the database cannot be established or if the retrieval operation fails.
// [ErrNoDocumentsFound]: If no document matches the key-value pair.
func GetOne(collectionName string, key string, value string, dest any) error {
	client, err := openConnection()
	if err != nil {
		return err
	}
	defer closeConnection(client)

	coll := client.Database(DB_NAME).Collection(collectionName)
	doc := coll.FindOne(context.TODO(), bson.D{{Key: key, Value: value}})
	err = doc.Decode(dest)
	if err == nil || err == mongo.ErrNoDocuments {
		return nil
	}
	return err
}

// Opens a connection to the database and inserts the provided document into the specified collection.
//
// [ErrInternalServerError]: If a connection to the database cannot be established.
// [ErrItemAlreadyPresent]: If there is a collision with the primary key of an existing item in the database.
func PutOne(collectionName string, doc any) error {
	client, err := openConnection()
	if err != nil {
		return err
	}
	defer closeConnection(client)

	coll := client.Database(DB_NAME).Collection(collectionName)
	_, err = coll.InsertOne(context.TODO(), doc)
	if err != nil {
		return err
	}
	return nil
}

//Opens a connection to the database and deletes a document that matches the provided key-value pair.

// [ErrInternalServerError]: If a connection to the database cannot be established or if the delete operation fails.
// [ErrNoDocumentsFound]: If no document matches the key-value pair.
func DeleteOne(collectionName string, key string, value string) error {
	client, err := openConnection()
	if err != nil {
		return err
	}
	defer closeConnection(client)

	coll := client.Database(DB_NAME).Collection(collectionName)
	_, err = coll.DeleteOne(context.TODO(), bson.D{{Key: key, Value: value}})
	if err != nil {
		return err
	}
	return nil
}

// creates a sort document to correctly sort a mongoDB resul
func CreateSort(k string, v int) bson.D {
	return bson.D{{Key: k, Value: v}}
}
