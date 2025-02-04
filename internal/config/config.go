package config

import (
	"context"
	//"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var DB *mongo.Client

// ConnectDB initializes the MongoDB client and establishes a connection.
func ConnectDB() *mongo.Client {

	// Retrieve the MongoDB URI from environment variables
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI is not set in the environment variables")
	}
	//fmt.Println("MONGODB_URI:", os.Getenv("MONGODB_URI"))

	// Set client options with the MongoDB URI
	// This line creates a MongoDB client options instance by calling options.Client().
	clientOpts := options.Client().ApplyURI(mongoURI)

	// Connect to MongoDB
	// This line initializes a new MongoDB client instance using the
	// Note: In MongoDB v2 Go Driver, we don't need to pass context.TODO() or context.Background() to mongo.Connect() because it's now managed internally.
	client, err := mongo.Connect(clientOpts)
	if err != nil {
		log.Fatal("Failed to create MongoDB client:", err)
	}

	// Ping the database to verify the connection
	// In MongoDB, context ensures queries do not hang indefinitely.
	// context.WithTimeout() ensures that the database operation does not run indefinitely and gets canceled automatically if it takes longer than 10 seconds.
	//Breaking it down:
	//context.Background(): A parent context that is typically used when you’re not deriving from another context.
	//context.WithTimeout(): Creates a new context with a deadline.
	//cancel: A function that releases resources when the operation is done.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//Ensures that if MongoDB takes too long to respond, the request is canceled to prevent the system from hanging.
	defer cancel() // Always call cancel() to free up resources

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	log.Println("Successfully connected to MongoDB ☭")
	DB = client
	return client
}

// GetCollection returns a reference to a specific collection in the database.
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatal("DB_NAME is not set in the environment variables")
	}
	return client.Database(dbName).Collection(collectionName)
}
