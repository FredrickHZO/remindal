package main

import (
	"net/url"
	db "remindal/internal/database"
	"strconv"
)

var (
	MULTI_SEL_SEPARATOR = ","
)

func exists(s string) bool {
	return s != ""
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
		qb.AddRangeFieldC("age", minAge, true, func(s string) (any, error) {
			return strconv.Atoi(s)
		})
	}
	maxAge := q.Get("maxage")
	if exists(maxAge) {
		qb.AddRangeFieldC("age", maxAge, false, func(s string) (any, error) {
			return strconv.Atoi(s)
		})
	}
	if !exists(minAge) && !exists(maxAge) {
		age := q.Get("age")
		if age != "" {
			qb.AddFieldC("age", age, func(s string) (any, error) {
				return strconv.Atoi(s)
			})
		}
	}
}
