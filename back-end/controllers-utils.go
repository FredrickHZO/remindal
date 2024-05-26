package main

import (
	"net/url"
	db "remindal/internal/database"
	"strconv"
)

var (
	MULTI_SEL_SEPARATOR = ","
)

// Checks if all the possible User filters are inside the HTTP URL query
func buildUserQuery(q url.Values, qb *db.QueryBuilder) {
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
	if minAge == "" && maxAge == "" {
		age := q.Get("age")
		if age != "" {
			qb.AddFieldC("age", age, func(s string) (any, error) {
				return strconv.Atoi(s)
			})
		}
	}
}
