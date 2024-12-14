package domain

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `json:"id" bson:"_id"`
	Username string    `json:"username" bson:"username"`
	Password string    `json:"password" bson:"password"`
}
