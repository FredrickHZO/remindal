package main

type User struct {
	Email    string `bson:"_id" json:"email,omitempty" validate:"required,email"`
	Password string `bson:"password" json:"password,omitempty" validate:"required"`
	Name     string `bson:"name" json:"name,omitempty" validate:"required"`
	Surname  string `bson:"surname" json:"surname,omitempty" validate:"required"`
	Age      uint8  `bson:"age,omitempty" json:"age,omitempty"`
}

type Date struct {
	Labels []string `bson:"labels,omitempty" json:"labels,omitempty"`
	Type   string   `bson:"type"`

	Year  uint16 `bson:"year"`
	Month uint8  `bson:"month"`
	Day   uint8  `bson:"day"`
}

type Calendar struct {
	Alias string `bson:"alias"`
	Dates []Date `bson:"dates"`
}
