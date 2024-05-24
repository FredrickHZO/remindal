package database

// TODO: Rework / remove

import (
	"net/url"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	keyMap = map[string]string{
		"email": "_id", // userSchema
		// ...
	}

	multiSelKeys = []string{"name"}    // test
	minRangeKeys = []string{"min_age"} // test
	maxRangeKeys = []string{"max_age"} // test
)

/*
ToMongoQuery converts the query parameters from an HTTP request into a MongoDB compatible query.

Iterates over the query parameters, identifies the type of filter (range, list, or single value),
and constructs a corresponding BSON filter.
*/
func ToMongoQuery(query url.Values) (bson.D, error) {
	queryDoc := bson.D{}
	conds := bson.A{}
	for k, v := range query {
		k = getBsonDatabaseKey(k)
		switch {
		case contains(multiSelKeys, k):
			queryDoc = append(queryDoc, multiSelFilter(k, v))

		case contains(minRangeKeys, k):
			filter, err := rangeFilter(k, v, "$gte")
			if err != nil {
				return bson.D{}, err
			}
			queryDoc = append(queryDoc, filter)

		case contains(maxRangeKeys, k):
			filter, err := rangeFilter(k, v, "$lte")
			if err != nil {
				return bson.D{}, err
			}
			queryDoc = append(queryDoc, filter)

		default:
			conds = append(conds, bson.D{{Key: k, Value: v}})
		}
	}
	queryDoc = append(queryDoc, bson.E{Key: "$or", Value: conds})

	return queryDoc, nil
}

/*
Converts the key name in the HTTP request to the corresponding database key name.
*/
func getBsonDatabaseKey(k string) string {
	if item, ok := keyMap[k]; ok {
		return item
	}
	return k
}

/*
Handles a multiple selection in the HTTP request. A multiple selection takes the form:

	key=valkey=val2key=val3key=val4key=val5

Splits the string containing the values and creates the appropriate query using the "$or" operator.
*/
func multiSelFilter(k string, v []string) primitive.E {
	var orConditions bson.A
	for _, item := range v {
		orConditions = append(orConditions, bson.D{{Key: k, Value: item}})
	}
	return primitive.E{Key: "$or", Value: orConditions}
}

/*
contains checks if a specific key is present in an array of keys.

Validates that a key from the HTTP request is allowed
to have certain types of values (e.g., multiple selection or range values).
By ensuring the key is within the predefined allowed keys, it helps maintain
the integrity of the query formation process.
*/
func contains(allowedKeys []string, k string) bool {
	for _, item := range allowedKeys {
		if item == k {
			return true
		}
	}
	return false
}

/*
Reimplement
*/
func rangeFilter(k string, v []string, c string) (bson.E, error) {
	val, err := strconv.Atoi(v[0])
	if err != nil {
		return bson.E{}, err
	}

	rangeCond := bson.E{
		Key:   k,
		Value: bson.D{bson.E{Key: c, Value: val}},
	}
	return rangeCond, nil
}
