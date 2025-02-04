package main

import (
	"log"
	"os"

	"github.com/alpha-154/crud-go-gin/internal/config"
	"github.com/alpha-154/crud-go-gin/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from the .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to MongoDB
	config.ConnectDB()

	// Create a new Gin router
	router := gin.Default()

	// Set up routes
	routes.SetupRoutes(router)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	log.Println("Server is running on port", port)
	err := router.Run(port)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
