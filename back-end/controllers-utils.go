package main

import (
	"net/url"
	db "remindal/internal/database"
	"strconv"
)

var (
	MULTI_SEL_SEPARATOR = ","
)

// Checks if there are range values in the HTTP URL query for the specified key.
// If one or both range values are present, prevents the single value to be
// added as a query filter.
func addSingleFieldIfNoRangePresent(min, max, single string, qb *db.QueryBuilder) {
	if min != "" && max != "" {
		return
	}
	if single != "" {
		qb.AddFieldC("age", single, func(s string) (any, error) {
			return strconv.Atoi(s)
		})
	}
}

// Checks if all the possible User filters are inside the HTTP URL query
func constructUserQuery(q url.Values, qb *db.QueryBuilder) {
	email := q.Get("_id")
	if email != "" {
		qb.AddField("_id", email)
	}
	password := q.Get("password")
	if password != "" {
		qb.AddField("password", password)
	}
	names := q.Get("name")
	if names != "" {
		qb.AddMultiSelectField("name", names, MULTI_SEL_SEPARATOR)
	}
	surnames := q.Get("surname")
	if surnames != "" {
		qb.AddMultiSelectField("surname", surnames, MULTI_SEL_SEPARATOR)
	}
	minAge := q.Get("minage")
	if minAge != "" {
		qb.AddRangeFieldC("age", minAge, true, func(s string) (any, error) {
			return strconv.Atoi(s)
		})
	}
	maxAge := q.Get("maxage")
	if maxAge != "" {
		qb.AddRangeFieldC("age", maxAge, false, func(s string) (any, error) {
			return strconv.Atoi(s)
		})
	}
	age := q.Get("age")
	addSingleFieldIfNoRangePresent(minAge, maxAge, age, qb)
}
