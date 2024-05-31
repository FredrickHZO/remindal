package main

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type User struct {
	Email    string `bson:"_id" json:"_id,omitempty" validate:"required,email"`
	Password string `bson:"password" json:"password,omitempty" validate:"required"`
	Name     string `bson:"name" json:"name,omitempty" validate:"required"`
	Surname  string `bson:"surname" json:"surname,omitempty" validate:"required"`
	Age      uint8  `bson:"age,omitempty" json:"age,omitempty"`
}

type Date struct {
	ID string `bson:"_id,omitempty" json:"_id,omitempty"`

	Labels []string `bson:"labels,omitempty" json:"labels,omitempty"`
	Type   string   `bson:"type" json:"type" validate:"required"`
	Desc   string   `bson:"description,omitempty" json:"description,omitempty"`

	Year    int16 `bson:"year" json:"year" validate:"required"`
	Month   int8  `bson:"month" json:"month" validate:"required,min=1,max=12"`
	Day     int8  `bson:"day" json:"day" validate:"required,day_validation"`
	Hours   int8  `bson:"hours,omitempty" json:"hours,omitempty" validate:"min=0,max=23"`
	Minutes int8  `bson:"minutes,omitempty" json:"minutes,omitempty" validate:"min=0,max=59"`
}

// Checks if the date is a valid one by getting the total days in
// the specified month for the specific year and comparing the day in input with
// the total days value. If the input day value is less then 0 or greater
// than the total days in that given month, return false, true otherwise
func dayValidation(fl validator.FieldLevel) bool {
	year := fl.Parent().FieldByName("Year").Int()
	month := fl.Parent().FieldByName("Month").Int()
	day := fl.Field().Int()

	// 32 is greater than any possible value, it will be automatically normalized.
	// 2024-02-32 becomes 2024-03-03 after normalization. 32 - 3 = 29
	t := time.Date(int(year), time.Month(month), 32, 0, 0, 0, 0, time.UTC)
	maxDays := 32 - t.Day()

	if day < 1 || int(day) > maxDays {
		return false
	}
	return true
}

// returns a new validator with registered custom validators for the object Date
func newCustomDateValidator() *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("day_validation", dayValidation)
	return validate
}

type Calendar struct {
	Alias string `bson:"alias"`
	Dates []Date `bson:"dates"`
}
