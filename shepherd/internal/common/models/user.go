package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// UserCredentials - describes a User credentials for login operation.
type UserCredentials struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// User - describes a single user of a Shepherd system.
type User struct {
	BaseModel `bson:",inline"`

	Name     string `bson:"name" json:"name"`
	Email    string `bson:"email" json:"email"`
	Password string `bson:"password" json:"password"`
	Role     string `bson:"role" json:"role"`
}

// NewUser - returns a new User instance.
func NewUser(name, password string) *User {
	user := &User{
		Name:     name,
		Password: password,
	}

	user.ID = primitive.NewObjectID()

	return user
}
