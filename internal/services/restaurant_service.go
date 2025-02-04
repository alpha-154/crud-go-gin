package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/alpha-154/crud-go-gin/internal/config"
	"github.com/alpha-154/crud-go-gin/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var restaurantCollection *mongo.Collection

// Ensure DB is initialized before accessing the collection
func getRestaurantCollection() *mongo.Collection {
	if restaurantCollection == nil {
		client := config.DB
		if client == nil {
			log.Fatal("Database connection is not initialized")
		}
		restaurantCollection = config.GetCollection(client, "restaurants")
	}
	return restaurantCollection
}

func CreateRestaurant(restaurant models.Restaurant) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Generate a new MongoDB ObjectID
	restaurant.ID = primitive.NewObjectID()

	// Convert ObjectID to a string and store it in RestaurantID
	restaurant.RestaurantID = restaurant.ID.Hex()

	// Insert restaurant into MongoDB
	result, err := getRestaurantCollection().InsertOne(ctx, restaurant)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetAllRestaurants() ([]map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := getRestaurantCollection().Find(ctx, bson.M{}) // Fetch all documents
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var response []map[string]interface{}

	for cursor.Next(ctx) {
		var restaurant bson.M // Use `bson.M` instead of `models.Restaurant`
		if err := cursor.Decode(&restaurant); err != nil {
			return nil, err
		}

		// Convert `_id` to string if it exists
		if id, ok := restaurant["_id"].(primitive.ObjectID); ok {
			restaurant["_id"] = id.Hex() // Convert ObjectID to string
		}

		response = append(response, restaurant)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return response, nil
}

func GetRestaurantByID(id string) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var restaurant bson.M

	// Query using restaurant_id instead of _id
	filter := bson.M{"restaurant_id": id}
	result := getRestaurantCollection().FindOne(ctx, filter)

	fmt.Println("Searching for RestaurantID:", id)

	// Log query result
	if result.Err() != nil {
		fmt.Println("Error finding document:", result.Err()) // Debugging line
		return bson.M{}, result.Err()
	}

	// Decode the result into restaurant variable
	err := result.Decode(&restaurant)
	if err != nil {
		fmt.Println("Error decoding document:", err) // Debugging line
		return bson.M{}, err
	}

	fmt.Println("Found restaurant:", restaurant)
	return restaurant, nil
}

// func GetRestaurantByName(name string) (map[string]interface{}, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	log.Println("Querying database for restaurant with name:", name)

// 	var restaurant bson.M
// 	err := getRestaurantCollection().FindOne(ctx, bson.M{"name": name}).Decode(&restaurant)
// 	if err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			log.Println("No restaurant found with name:", name)
// 		} else {
// 			log.Println("Error querying database:", err)
// 		}
// 		return bson.M{}, err
// 	}

// 	log.Println("Restaurant found with name:", name)
// 	return restaurant, nil
// }

func UpdateRestaurant(id string, updatedData models.Restaurant) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(id)
	update := bson.M{
		"$set": updatedData,
	}
	result, err := getRestaurantCollection().UpdateOne(ctx, bson.M{"_id": objID}, update)
	return result, err
}

func DeleteRestaurant(id string) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(id)
	result, err := getRestaurantCollection().DeleteOne(ctx, bson.M{"_id": objID})
	return result, err
}
