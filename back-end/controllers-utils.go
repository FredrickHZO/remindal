package main

import (
	"net/url"
	db "remindal/internal/database"
	"strconv"
)

var MULTI_SEL_SEPARATOR = ","

// Possible keys that make up a calendar query
var (
	CALENDAR_TYPE = "type"
	LABELS        = "labels"

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
var (
	EMAIL    = "_id"
	PASSWORD = "password"

	SURNAME = "surname"
	NAME    = "name"

	MIN_AGE = "minage"
	MAX_AGE = "maxage"
	AGE     = "age"
)

// Checks if a param is not an empty string. I.E. if it exists
func exists(s string) bool {
	return s != ""
}

// Wrapper for conversion from string to integer
func paramToi(s string) (any, error) {
	return strconv.Atoi(s)
}

// If the given parameters used as ranges do not exist, then a single filter is added
func addFilterIfNoRangeExists(k string, v string, min string, max string, qb *db.QueryBuilder, cnv func(s string) (any, error)) {
	if exists(min) || exists(max) {
		return
	}
	if cnv != nil {
		qb.AddFieldC(k, v, cnv)
	} else {
		qb.AddField(k, v)
	}
}

func addRangeFilter(k string, min string, max string, qb *db.QueryBuilder, cnv func(s string) (any, error)) {
	if exists(min) {
		if cnv != nil {
			qb.AddRangeFieldC(k, min, true, cnv)
		} else {
			qb.AddRangeField(k, min, true)
		}
	}
	if exists(max) {
		if cnv != nil {
			qb.AddRangeFieldC(k, max, false, cnv)
		} else {
			qb.AddRangeField(k, min, false)
		}
	}
}

// Checks the valid User filters inside the HTTP URL query and builds a mongoDB query.
func buildUserQuery(q url.Values, qb *db.QueryBuilder) {
	email := q.Get(EMAIL)
	if exists(email) {
		qb.AddField(EMAIL, email)
	}
	password := q.Get(PASSWORD)
	if exists(password) {
		qb.AddField(PASSWORD, password)
	}
	names := q.Get(NAME)
	if exists(names) {
		qb.AddMultiSelectField(NAME, names, MULTI_SEL_SEPARATOR)
	}
	surnames := q.Get(SURNAME)
	if exists(surnames) {
		qb.AddMultiSelectField(SURNAME, surnames, MULTI_SEL_SEPARATOR)
	}

	minAge := q.Get(MIN_AGE)
	maxAge := q.Get(MAX_AGE)
	addRangeFilter(AGE, minAge, maxAge, qb, paramToi)
	age := q.Get(AGE)
	addFilterIfNoRangeExists(AGE, age, minAge, maxAge, qb, paramToi)
}

// Checks the valid User filters inside the HTTP URL query and builds a mongoDB query.
func buildCalendarQuery(q url.Values, qb *db.QueryBuilder) {
	labels := q.Get(LABELS)
	if exists(labels) {
		qb.AddMultiSelectField(LABELS, labels, MULTI_SEL_SEPARATOR)
	}
	calType := q.Get(CALENDAR_TYPE)
	if exists(calType) {
		qb.AddField(CALENDAR_TYPE, calType)
	}

	minYear := q.Get(MIN_YEAR)
	maxYear := q.Get(MAX_YEAR)
	addRangeFilter(YEAR, minYear, maxYear, qb, paramToi)
	year := q.Get(YEAR)
	addFilterIfNoRangeExists(YEAR, year, minYear, maxYear, qb, paramToi)

	minMonth := q.Get(MIN_MONTH)
	maxMonth := q.Get(MAX_MONTH)
	addRangeFilter(MONTH, minMonth, maxMonth, qb, paramToi)
	month := q.Get(MONTH)
	addFilterIfNoRangeExists(MONTH, month, minMonth, maxMonth, qb, paramToi)

	minDay := q.Get(MIN_DAY)
	maxDay := q.Get(MAX_DAY)
	addRangeFilter(DAY, minDay, maxDay, qb, paramToi)
	day := q.Get(DAY)
	addFilterIfNoRangeExists(DAY, day, minDay, maxDay, qb, paramToi)

	minHours := q.Get(MIN_HOURS)
	maxHours := q.Get(MAX_HOURS)
	addRangeFilter(HOURS, minHours, maxHours, qb, paramToi)
	hours := q.Get(HOURS)
	addFilterIfNoRangeExists(HOURS, hours, minHours, maxHours, qb, paramToi)

	minMinutes := q.Get(MIN_MINUTES)
	maxMinutes := q.Get(MAX_MINUTES)
	addRangeFilter(MINUTES, minMinutes, maxMinutes, qb, paramToi)
	minutes := q.Get(MINUTES)
	addFilterIfNoRangeExists(MINUTES, minutes, minMinutes, maxMinutes, qb, paramToi)
}
