package database

// TODO: Rework / remove

import (
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

type QueryBuilder struct {
	err   error
	query bson.D
}

func NewQueryBuilder() QueryBuilder {
	return QueryBuilder{
		err:   nil,
		query: bson.D{},
	}
}

// adds a simple field filter to the query document
func (qb *QueryBuilder) AddField(k string, v string) {
	qb.query = append(qb.query, bson.E{Key: k, Value: v})
}

// adds a multi select filter to the query document
func (qb *QueryBuilder) AddMultiSelectField(k string, v string, sep string) {
	vals := strings.Split(v, sep)
	orConditions := bson.A{}
	for _, item := range vals {
		orConditions = append(orConditions, bson.D{{Key: k, Value: item}})
	}
	qb.query = append(qb.query, bson.E{Key: "$or", Value: orConditions})
}

// adds a range filter to the query document
func (qb *QueryBuilder) AddRangeField(k string, v string, c string) {
	rangeCond := bson.E{
		Key:   k,
		Value: bson.D{bson.E{Key: c, Value: v}},
	}
	qb.query = append(qb.query, rangeCond)
}

// converts the value and adds a simple field filter to the query document
// adds an error to the QueryBuilder if the convertion is unsuccessfull
func (qb *QueryBuilder) AddFeildC(k string, v string, cnv func(s string) (any, error)) {
	val, err := cnv(v)
	if err != nil {
		qb.err = errors.Join(err)
		return
	}
	qb.query = append(qb.query, bson.E{Key: k, Value: val})

}

// converts the values and adds a multi select filter to the query document
// adds an error to the QueryBuilder if the convertion is unsuccessfull
func (qb *QueryBuilder) AddMultiSelectC(k string, v string, sep string, cnv func(s string) (any, error)) {
	vals := strings.Split(v, sep)

	orConditions := bson.A{}
	for _, item := range vals {
		val, err := cnv(item)
		if err != nil {
			qb.err = errors.Join(err)
			return
		}
		orConditions = append(orConditions, bson.D{{Key: k, Value: val}})
	}
	qb.query = append(qb.query, bson.E{Key: "$or", Value: orConditions})
}

// converts the value and adds a range filter to the query document
// adds an error to the QueryBuilder if the convertion is unsuccessfull
func (qb *QueryBuilder) AddRangeFieldC(k string, v string, min bool, cnv func(s string) (any, error)) {
	val, err := cnv(v)
	if err != nil {
		qb.err = errors.Join(err)
		return
	}

	var cond string
	if min {
		cond = "$gte"
	} else {
		cond = "$lte"
	}
	rangeCond := bson.E{
		Key:   k,
		Value: bson.D{bson.E{Key: cond, Value: val}},
	}
	qb.query = append(qb.query, rangeCond)
}

// returns the error in the QueryBuilder
func (qb *QueryBuilder) Err() error {
	return qb.err
}

// return the built query
func (qb *QueryBuilder) Query() bson.D {
	return qb.query
}
