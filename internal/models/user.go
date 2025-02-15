package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID       string             `bson:"user_id" json:"user_id"`
	Name         string             `bson:"name" json:"name" binding:"required"`
	Email        string             `bson:"email" json:"email" binding:"required,email"`
	Password     string             `bson:"password" json:"password" binding:"required,min=6"`
	Role         string             `bson:"role" json:"role"`
	RefreshToken string             `bson:"refresh_token" json:"refresh_token,omitempty"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}
