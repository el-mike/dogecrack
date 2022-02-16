package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BaseModel - model with base fields for other models to embed.
type BaseModel struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt NullableTime       `bson:"createdAt" json:"createdAt"`
	UpdatedAt NullableTime       `bson:"updatedAt" json:"updatedAt"`
}
