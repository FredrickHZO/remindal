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

Iterates over the query parameters, identifies the type of filter (range, list, or single value),
and constructs a corresponding BSON filter.
*/
func ToMongoQuery(query url.Values) (bson.D, error) {
	var queryDoc bson.D
	for k, v := range query {
		key := getBsonDatabaseKey(k)
		val := v[0]
		switch {
		case strings.Contains(val, LIST):
			filter, err := multiSelectFilter(key, val)
			if err != nil {
				return queryDoc, err
			}
			queryDoc = append(queryDoc, filter)

		case strings.Contains(val, RANGE):
			filter, err := rangeFilter(key, val)
			if err != nil {
				return queryDoc, err
			}
			queryDoc = append(queryDoc, filter)

		default:
			queryDoc = append(
				queryDoc,
				primitive.E{Key: key, Value: val},
			)
		}
	}
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

	key=val1%,val2%,val3%,val4%,val5

Splits the string containing the values and creates the appropriate query using the "$or" operator.
*/
func multiSelectFilter(k string, v string) (primitive.E, error) {
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
Handles a range in the HTTP request. A range takes the form:

	key=val1%-val2

Splits the string containing the range and creates the appropriate query.
If only one value is present in the range, the resulting query will include all values
greater than or equal to (or less than or equal to) the specified value.
*/
func rangeFilter(k string, v string) (primitive.E, error) {
	if !contains(rangeKeys, k) {
		return primitive.E{}, remerr.ErrNotRangeable
	}

	rng := strings.Split(v, RANGE)
	if len(rng) > 2 {
		return primitive.E{}, remerr.ErrInvalidRangeValues
	}

	var filter bson.D
	switch {
	case strings.HasPrefix(v, RANGE):
		filter, err := singleRange(rng[1], filter, "$lte")
		if err != nil {
			return primitive.E{}, err
		}
		return primitive.E{Key: k, Value: filter}, nil

	case strings.HasSuffix(v, RANGE):
		filter, err := singleRange(rng[0], filter, "$gte")
		if err != nil {
			return primitive.E{}, err
		}
		return primitive.E{Key: k, Value: filter}, nil

	default:
		filter, err := fullRange(rng, filter)
		if err != nil {
			return primitive.E{}, err
		}
		return primitive.E{Key: k, Value: filter}, nil
	}
}

/*
Helper - processes a range that has a single value.
The condition for the query must be specified, either "$gte" or "$lte".
*/
func singleRange(v string, query bson.D, cond string) (bson.D, error) {
	lte, err := strconv.Atoi(v)
	if err != nil {
		return query, remerr.ErrRangeValueNotNumber
	}
	query = append(query, primitive.E{Key: cond, Value: lte})
	return query, nil
}

/*
Helper - processes a range and creates a query with the correct range format.
*/
func fullRange(v []string, query bson.D) (bson.D, error) {
	gte, err := strconv.Atoi(v[0])
	if err != nil {
		return query, remerr.ErrInternalServerError
	}

	lte, err := strconv.Atoi(v[1])
	if err != nil {
		return query, remerr.ErrInternalServerError
	}

	if gte > lte {
		return query, remerr.ErrInvalidRangeValues
	}

	query = append(query, primitive.E{Key: "$gte", Value: gte})
	query = append(query, primitive.E{Key: "$lte", Value: lte})
	return query, nil
}
