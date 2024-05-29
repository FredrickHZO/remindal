package main

import (
	"time"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Email    string `bson:"_id" json:"_id,omitempty" validate:"required,email"`
	Password string `bson:"password" json:"password,omitempty" validate:"required"`
	Name     string `bson:"name" json:"name,omitempty" validate:"required"`
	Surname  string `bson:"surname" json:"surname,omitempty" validate:"required"`
	Age      uint8  `bson:"age,omitempty" json:"age,omitempty"`
}

type Date struct {
	ID primitive.ObjectID `bson:"_id" json:"_id"`

	Labels []string `bson:"labels,omitempty" json:"labels,omitempty"`
	Type   string   `bson:"type" json:"type" validate:"required"`

	Year    uint16 `bson:"year" json:"year" validate:"required"`
	Month   uint8  `bson:"month" json:"month" validate:"required,min=0,max=12"`
	Day     uint8  `bson:"day" json:"day" validate:"required,dayValidation"`
	Hours   uint8  `bson:"hours,omitempty" json:"hours,omitempty" validate:"min=0,max=23"`
	Minutes uint8  `bson:"minutes,omitempty" json:"minutes,omitempty" validate:"min=0,max=59"`
}

// Checks if the date is a valid one by getting the total days in
// the specified month and year and comparing the day in input with
// the total days value. If the input day value is less then 0 or greater
// than the total days in that given month, return false, true otherwise
func dayValidation(fl validator.FieldLevel) bool {
	year := fl.Parent().FieldByName("Year").Int()
	month := fl.Parent().FieldByName("Month").Int()
	day := fl.Parent().FieldByName("Day").Int()

	// 32 is greater than any possible value, it will get automatically normalized.
	// by getting the normalized day and subtracting it to 32, we get the actual correct
	// number of days for that month in that given year (needed for February)
	// 2024-02-32 becomes 2024-03-04. 32 - 4 = 29
	t := time.Date(int(year), time.Month(month), 32, 0, 0, 0, 0, time.UTC)
	maxDay := 32 - t.Day()

	if day < 1 {
		return false
	}
	if int(day) > maxDay {
		return false
	}
	return true
}

type Calendar struct {
	Alias string `bson:"alias"`
	Dates []Date `bson:"dates"`
}
