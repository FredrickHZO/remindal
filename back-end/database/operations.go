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
WIP - Opens connection to the database and retrieves an array of items that have [value] as its primary key / unique key value.
[key] is the primary key / unique key to search for in the database

	{ _id: email@person.com, name: John }
	key = _id
	value = email@person.com

Usage

	var destination []mySchema
	err := GetMany(myCollection, &destination)

returns [ErrInternalServerError] in case a connection with the database can't be established.

returns [ErrInternalServerError] in case the retrieval operations fails.

returns [ErrNoDocumentsFound] in case no item matches the value provided.
*/
func GetMany(collection string, dest any) error {
	client, err := openConnection()
	if err != nil {
		return remerr.ErrInternalServerError
	}
	defer closeConnection(client)

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

/*
Opens connection to the database and retrieves an item that has [value] as its primary key / unique key value.
[key] is the primary key / unique key to search for in the database

	{ _id: email@person.com, name: John }
	key = _id
	value = email@person.com

Usage

	var destination mySchema
	err := GetOne(myCollection, &destination)

returns [ErrInternalServerError] in case a connection with the database can't be established.

returns [ErrInternalServerError] in case the retrieval operations fails.

returns [ErrNoDocumentsFound] in case no item matches the value provided.
*/
func GetOne(collection string, key string, value string, dest any) error {
	client, err := openConnection()
	if err != nil {
		return remerr.ErrInternalServerError
	}
	defer closeConnection(client)

	c := client.Database(DB_NAME).Collection(collection)

	doc := c.FindOne(CTX, bson.D{{Key: key, Value: value}})
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
Opens connection to the database and inserts the document provided.

returns [ErrInternalServerError] in case a connection with the database can't be established.

returns [ErrItemAlreadyPresent] in case there is a collision with the primary key of an item in the database.
*/
func PutOne(collection string, doc any) error {
	client, err := openConnection()
	if err != nil {
		return remerr.ErrInternalServerError
	}
	defer closeConnection(client)

	c := client.Database(DB_NAME).Collection(collection)

	// this is not a correct way to do this, must be changed.
	_, err = c.InsertOne(CTX, doc)
	if err != nil {
		return remerr.ErrItemAlreadyPresent
	}
	return nil
}

/*
Opens connection to the database and deletes an item that has [value] as its primary key / unique key value.
[key] is the primary key / unique key to search for in the database

	{ _id: email@person.com, name: Jhon }
	key = _id
	value = email@person.com

returns [ErrInternalServerError] in case a connection with the database can't be established.

returns [ErrInternalServerError] in case the delete operations fails.

returns [ErrNoDocumentsFound] in case no item matches the value provided.
*/
func DeleteOne(collection string, key string, value string) error {
	client, err := openConnection()
	if err != nil {
		return remerr.ErrInternalServerError
	}
	defer closeConnection(client)

	c := client.Database(DB_NAME).Collection(collection)

	delres, err := c.DeleteOne(CTX, bson.D{{Key: key, Value: value}})
	if delres.DeletedCount == 0 {
		return remerr.ErrNoDocumentsFound
	}
	if err != nil {
		return remerr.ErrInternalServerError
	}
	return nil
}
