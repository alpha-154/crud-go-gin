package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Restaurant struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	RestaurantID string             `bson:"restaurant_id" json:"restaurant_id"` // Store ObjectID as a string
	Name         string             `json:"name"`
	Address      string             `json:"address"`
	Email        string             `json:"email"`
	Cuisine      string             `json:"cuisine"`
}
