package database

type UserSchema struct {
	Email   string `bson:"_id" json:"email"`
	Name    string `bson:"name"`
	Surname string `bson:"surname"`
	Age     int    `bson:"age"`
}
