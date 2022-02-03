package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BaseModel - model with base fields for other models to embed.
type BaseModel struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
