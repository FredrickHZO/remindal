package main

import (
	"net/url"
	db "remindal/internal/database"
	"strconv"
)

// alias for conversion functions
type convFunc func(s string) (any, error)

// character used to separate list of item in http query
const MULTI_SEL_SEPARATOR = ","

// Possible keys that make up a calendar query
const (
	DATE_TYPE = "type"
	LABELS    = "labels"

	MIN_YEAR = "minyear"
	MAX_YEAR = "maxyear"
	YEAR     = "year"

	MIN_MONTH = "minmonth"
	MAX_MONTH = "maxmonth"
	MONTH     = "month"

	MIN_DAY = "minday"
	MAX_DAY = "maxday"
	DAY     = "day"

	MIN_HOURS = "minhours"
	MAX_HOURS = "maxhours"
	HOURS     = "hours"

	MIN_MINUTES = "minminutes"
	MAX_MINUTES = "maxminutes"
	MINUTES     = "minutes"
)

// Possible keys that make up a user query
const (
	EMAIL    = "_id"
	PASSWORD = "password"

	SURNAME = "surname"
	NAME    = "name"

	MIN_AGE = "minage"
	MAX_AGE = "maxage"
	AGE     = "age"
)

// Checks if a param is not an empty string. I.E. if it has value
func hasValue(s string) bool {
	return s != ""
}

// Wrapper for conversion from string to integer
func paramToi(s string) (any, error) {
	return strconv.Atoi(s)
}

// If the given parameters used as ranges do not exist, then a single filter is added
func addFilterIfNoRangeExists(k, v, min, max string, c convFunc, b *db.QueryBuilder) {
	if hasValue(min) || hasValue(max) || !hasValue(v) {
		return
	}
	b.AddFieldCnv(k, v, c)
}

// Adds a range filter with a min and max value to the query. If both values are empty strings
// this is a no operation and no filter is added.
func addRangeFilter(k, min, max string, c convFunc, b *db.QueryBuilder) {
	if hasValue(min) {
		b.AddRangeField(k, min, true, c)
	}
	if hasValue(max) {
		b.AddRangeField(k, max, false, c)
	}
}

func addSimpleFilter(k, v string, b *db.QueryBuilder) {
	if !hasValue(v) {
		return
	}
	b.AddField(k, v)
}

func addMultiSelectFilter(k, v, sep string, b *db.QueryBuilder) {
	if !hasValue(v) {
		return
	}
	b.AddMultiSelectField(k, v, sep)
}

// Checks the valid User filters inside the HTTP URL query and builds a mongoDB query.
func buildUserQuery(q url.Values, b *db.QueryBuilder) {
	email := q.Get(EMAIL)
	addSimpleFilter(EMAIL, email, b)

	password := q.Get(PASSWORD)
	addSimpleFilter(PASSWORD, password, b)

	names := q.Get(NAME)
	addMultiSelectFilter(NAME, names, MULTI_SEL_SEPARATOR, b)

	surnames := q.Get(SURNAME)
	addMultiSelectFilter(SURNAME, surnames, MULTI_SEL_SEPARATOR, b)

	minAge := q.Get(MIN_AGE)
	maxAge := q.Get(MAX_AGE)
	addRangeFilter(AGE, minAge, maxAge, paramToi, b)
	age := q.Get(AGE)
	addFilterIfNoRangeExists(AGE, age, minAge, maxAge, paramToi, b)
}

// Checks the valid Calendar filters inside the HTTP URL query and builds a mongoDB query.
func buildDateQuery(q url.Values, b *db.QueryBuilder) {
	labels := q.Get(LABELS)
	addMultiSelectFilter(LABELS, labels, MULTI_SEL_SEPARATOR, b)

	dateType := q.Get(DATE_TYPE)
	addSimpleFilter(DATE_TYPE, dateType, b)

	minYear := q.Get(MIN_YEAR)
	maxYear := q.Get(MAX_YEAR)
	addRangeFilter(YEAR, minYear, maxYear, paramToi, b)
	year := q.Get(YEAR)
	addFilterIfNoRangeExists(YEAR, year, minYear, maxYear, paramToi, b)

	minMonth := q.Get(MIN_MONTH)
	maxMonth := q.Get(MAX_MONTH)
	addRangeFilter(MONTH, minMonth, maxMonth, paramToi, b)
	month := q.Get(MONTH)
	addFilterIfNoRangeExists(MONTH, month, minMonth, maxMonth, paramToi, b)

	minDay := q.Get(MIN_DAY)
	maxDay := q.Get(MAX_DAY)
	addRangeFilter(DAY, minDay, maxDay, paramToi, b)
	day := q.Get(DAY)
	addFilterIfNoRangeExists(DAY, day, minDay, maxDay, paramToi, b)

	minHours := q.Get(MIN_HOURS)
	maxHours := q.Get(MAX_HOURS)
	addRangeFilter(HOURS, minHours, maxHours, paramToi, b)
	hours := q.Get(HOURS)
	addFilterIfNoRangeExists(HOURS, hours, minHours, maxHours, paramToi, b)

	minMinutes := q.Get(MIN_MINUTES)
	maxMinutes := q.Get(MAX_MINUTES)
	addRangeFilter(MINUTES, minMinutes, maxMinutes, paramToi, b)
	minutes := q.Get(MINUTES)
	addFilterIfNoRangeExists(MINUTES, minutes, minMinutes, maxMinutes, paramToi, b)
}
