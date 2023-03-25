package modal

import "go.mongodb.org/mongo-driver/bson/primitive"

type Saree struct {
	ID    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name  string             `json:"name,omitempty" bson:"name,omitempty"`
	Price int                `json:"price,omitempty" bson:"price,omitempty"`
	Image byte               `json:"image,omitempty" bson:"image,omitempty"`
	Type  string             `json:"type,omitempty" bson:"type,omitempty"`
	Color string             `json:"color,omitempty" bson:"color,omitempty"`
}
