package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Note struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title,omitempty"`
	Description string             `json:"description"`
}
