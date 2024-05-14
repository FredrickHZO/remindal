package database

import (
	"net/url"
	remerr "remindal/errors"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	keyMap = map[string]string{
		"email": "_id", // userSchema
		// ...
	}

	multipleSelectionKeys = []string{"name"} // test
	rangeKeys             = []string{"age"}  // test

	RANGE = "%-"
	LIST  = "%,"
)

/*
ToMongoQuery converts the query parameters from an HTTP request into a MongoDB compatible query.

The function iterates over the query parameters, identifies the type of filter (range, list, or single value),
and constructs a corresponding BSON filter.

@param query: The URL query parameters
@return bson.D: The BSON filters for the database query
@return error: Any error that occurred during conversion
*/
func ToMongoQuery(query url.Values) (bson.D, error) {
	var queryDoc bson.D
	for k, v := range query {
		key := getBsonDatabaseKey(k)
		val := v[0]
		switch {
		case strings.Contains(val, LIST):
			filter, err := processMultipleSelectionRequest(key, val)
			if err != nil {
				return queryDoc, err
			}
			queryDoc = append(queryDoc, filter)

		case strings.Contains(val, RANGE):
			filter, err := processRangeInRequest(key, val)
			if err != nil {
				return queryDoc, err
			}
			queryDoc = append(queryDoc, filter)

		default:
			queryDoc = append(queryDoc, primitive.E{Key: key, Value: val})
		}
	}
	return queryDoc, nil
}

/*
getBsonDatabaseKey converts the key name in the HTTP request to the corresponding database key name.

@param k: The key in the HTTP query
@return string: The corresponding database key name
*/
func getBsonDatabaseKey(k string) string {
	if item, ok := keyMap[k]; ok {
		return item
	}
	return k
}

/*
processMultipleSelectionRequest handles a multiple selection in the HTTP request. A multiple selection takes the form:

	key=val1%,val2%,val3%,val4%,val5

It splits the string containing the values and creates the appropriate query using the "$or" operator.

@param k: The key in the HTTP query
@param v: The value in the HTTP query containing the multiple selection
@return primitive.E: The BSON filter for the multiple selection
*/
func processMultipleSelectionRequest(k string, v string) (primitive.E, error) {
	if !contains(multipleSelectionKeys, k) {
		return primitive.E{}, remerr.ErrNotMultipleSelection
	}

	var orConditions bson.A
	arr := strings.Split(v, LIST)
	for _, item := range arr {
		orConditions = append(orConditions, bson.D{{Key: k, Value: item}})
	}
	return primitive.E{Key: "$or", Value: orConditions}, nil
}

/*
contains checks if a specific key is present in an array of keys.

This function is used to validate that a key from the HTTP request is allowed
to have certain types of values (e.g., multiple selection or range values).
By ensuring the key is within the predefined allowed keys, it helps maintain
the integrity of the query formation process.

@param allowedKeys: The array of allowed keys
@param key: The key to check for presence in the array
@return bool: True if the key is found in the array, otherwise false
*/
func contains(allowedKeys []string, key string) bool {
	for _, item := range allowedKeys {
		if item == key {
			return true
		}
	}
	return false
}

/*
processRangeInRequest handles a range in the HTTP request. A range takes the form:

	key=val1%-val2

It splits the string containing the range and creates the appropriate query.
If only one value is present in the range, the resulting query will include all values
greater than or equal to (or less than or equal to) the specified value.

@param k: The key in the HTTP query
@param v: The value in the HTTP query containing the range
@return primitive.E: The BSON filter for the range
@return error: Any error that occurred during conversion
*/
func processRangeInRequest(k string, v string) (primitive.E, error) {
	if !contains(rangeKeys, k) {
		return primitive.E{}, remerr.ErrNotRangeable
	}

	str := strings.Split(v, RANGE)
	if len(str) > 2 {
		return primitive.E{}, remerr.ErrInternalServerError
	}

	var query bson.D
	if strings.HasPrefix(v, RANGE) {
		query, err := singleValRange(str[1], query, "$lte")
		if err != nil {
			return primitive.E{}, err
		}
		return primitive.E{Key: k, Value: query}, nil
	}

	if strings.HasSuffix(v, RANGE) {
		query, err := singleValRange(str[0], query, "$gte")
		if err != nil {
			return primitive.E{}, err
		}
		return primitive.E{Key: k, Value: query}, nil
	}

	query, err := fullValRange(str, query)
	if err != nil {
		return primitive.E{}, err
	}
	return primitive.E{Key: k, Value: query}, nil
}

/*
singleValRange processes a range that has a single value.
The condition for the query must be specified, either "$gte" or "$lte".

@param v: The value in the range
@param cond: The BSON query document
@param logg: The condition for the range ("$gte" or "$lte")
@return bson.D: The updated BSON query document
@return error: Any error that occurred during conversion
*/
func singleValRange(v string, cond bson.D, logg string) (bson.D, error) {
	lte, err := strconv.Atoi(v)
	if err != nil {
		return cond, remerr.ErrRangeValueNotNumber
	}
	cond = append(cond, primitive.E{Key: logg, Value: lte})
	return cond, nil
}

/*
fullValRange processes a range and creates a query with the correct range format.

@param v: The values in the range
@param cond: The BSON query document
@return bson.D: The updated BSON query document
@return error: Any error that occurred during conversion
*/
func fullValRange(v []string, cond bson.D) (bson.D, error) {
	gte, err := strconv.Atoi(v[0])
	if err != nil {
		return cond, remerr.ErrInternalServerError
	}

	lte, err := strconv.Atoi(v[1])
	if err != nil {
		return cond, remerr.ErrInternalServerError
	}

	if gte > lte {
		return cond, remerr.ErrInvalidRangeValues
	}

	cond = append(cond, primitive.E{Key: "$gte", Value: gte})
	cond = append(cond, primitive.E{Key: "$lte", Value: lte})
	return cond, nil
}