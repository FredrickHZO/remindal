package main

type User struct {
	Email    string `bson:"_id" json:"email"`
	Password string `bson:"password" validate:"required,password"`
	Name     string `bson:"name,omitempty"`
	Surname  string `bson:"surname,omitempty"`
	Age      uint8  `bson:"age,omitempty"`
}

type Date struct {
	Labels []string `bson:"label"`
	Type   string   `bson:"type"`

	Year  string `bson:"year"`
	Month string `bson:"month"`
	Day   string `bson:"day"`
}

type Calendar struct {
	Alias string `bson:"alias"`
	Dates []Date `bson:"dates"`
}
