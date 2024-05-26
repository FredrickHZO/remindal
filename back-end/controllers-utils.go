package main

import (
	"net/url"
	db "remindal/internal/database"
	"strconv"
)

var (
	MULTI_SEL_SEPARATOR = ","
)

// Checks if a param is not an empty string. e.g. if it exists
func exists(s string) bool {
	return s != ""
}

// Wrapper to string conversion to integer
func paramToi(s string) (any, error) {
	return strconv.Atoi(s)
}

// Checks if all the possible User filters are inside the HTTP URL query
func buildUserQuery(q url.Values, qb *db.QueryBuilder) {
	email := q.Get("_id")
	if exists(email) {
		qb.AddField("_id", email)
	}
	password := q.Get("password")
	if exists(password) {
		qb.AddField("password", password)
	}
	names := q.Get("name")
	if exists(names) {
		qb.AddMultiSelectField("name", names, MULTI_SEL_SEPARATOR)
	}
	surnames := q.Get("surname")
	if exists(surnames) {
		qb.AddMultiSelectField("surname", surnames, MULTI_SEL_SEPARATOR)
	}
	minAge := q.Get("minage")
	if exists(minAge) {
		qb.AddRangeFieldC("age", minAge, true, paramToi)
	}
	maxAge := q.Get("maxage")
	if exists(maxAge) {
		qb.AddRangeFieldC("age", maxAge, false, paramToi)
	}
	if !exists(minAge) && !exists(maxAge) {
		age := q.Get("age")
		if exists(age) {
			qb.AddFieldC("age", age, paramToi)
		}
	}
}
