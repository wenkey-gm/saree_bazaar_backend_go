package domain

import "github.com/google/uuid"

type User struct {
	ID        uuid.UUID `json:"id" bson:"_id"`
	Username  string    `json:"username" bson:"username"`
	Password  string    `json:"password" bson:"password"`
	Email     string    `json:"email" bson:"email"`
	Phone     string    `json:"phone" bson:"phone"`
	Role      string    `json:"role" bson:"role"`
	FirstName string    `json:"firstName" bson:"firstName"`
	LastName  string    `json:"lastName" bson:"lastName"`
}
