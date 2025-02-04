package controllers

import (
	"net/http"

	"github.com/alpha-154/crud-go-gin/internal/models"
	"github.com/alpha-154/crud-go-gin/internal/services"
	"github.com/gin-gonic/gin"
)

// CreateRestaurant handles the request to create a new restaurant
func CreateRestaurant(c *gin.Context) {
	var restaurant models.Restaurant

	if err := c.ShouldBindJSON(&restaurant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.CreateRestaurant(restaurant)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)

}

// GetAllRestaurants retrieves all restaurants
func GetAllRestaurants(c *gin.Context) {
	restaurants, err := services.GetAllRestaurants()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, restaurants)
}

// GetRestaurant retrieves a restaurant by ID
func GetRestaurant(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing id parameter"})
		return
	}

	restaurant, err := services.GetRestaurantByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
		return
	}

	c.JSON(http.StatusOK, restaurant)
}

// // UpdateRestaurant updates restaurant details by ID
func UpdateRestaurant(c *gin.Context) {
	id := c.Param("id")
	var updatedRestaurant models.Restaurant

	if err := c.ShouldBindJSON(&updatedRestaurant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.UpdateRestaurant(id, updatedRestaurant)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Restaurant updated successfully", "data": result})
}

// DeleteRestaurant removes a restaurant by ID
func DeleteRestaurant(c *gin.Context) {
	id := c.Param("id")
	result, err := services.DeleteRestaurant(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Restaurant deleted successfully", "data": result})
}
