package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Email    string `bson:"_id" json:"_id,omitempty" validate:"required,email"`
	Password string `bson:"password" json:"password,omitempty" validate:"required"`
	Name     string `bson:"name" json:"name,omitempty" validate:"required"`
	Surname  string `bson:"surname" json:"surname,omitempty" validate:"required"`
	Age      uint8  `bson:"age,omitempty" json:"age,omitempty"`
}

// TODO: better validation
type Date struct {
	ID primitive.ObjectID `bson:"_id" json:"_id"`

	Labels []string `bson:"labels,omitempty" json:"labels,omitempty"`
	Type   string   `bson:"type" json:"type" validate:"required"`

	Year    uint16 `bson:"year" json:"year" validate:"required"`
	Month   uint8  `bson:"month" json:"month" validate:"required,min=0,max=12"`
	Day     uint8  `bson:"day" json:"day" validate:"required,min=1,max=31"`
	Hours   uint8  `bson:"hours,omitempty" json:"hours,omitempty" validate:"min=0,max=23"`
	Minutes uint8  `bson:"minutes,omitempty" json:"minutes,omitempty" validate:"min=0,max=59"`
}

type Calendar struct {
	Alias string `bson:"alias"`
	Dates []Date `bson:"dates"`
}
