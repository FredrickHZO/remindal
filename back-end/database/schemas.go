package database

type UserSchema struct {
	Email   string `bson:"_id" json:"email"`
	Name    string `bson:"name,omitempty"`
	Surname string `bson:"surname,omitempty"`
	Age     int    `bson:"age,omitempty"`
}
