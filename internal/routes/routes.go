package routes

import (
	"github.com/alpha-154/crud-go-gin/internal/controllers"
	"github.com/gin-gonic/gin"
)

// SetupRoutes initializes all API routes
func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.POST("/restaurants", controllers.CreateRestaurant)
		api.GET("/restaurants", controllers.GetAllRestaurants)
		api.GET("/restaurants/:id", controllers.GetRestaurant)
		api.PUT("/restaurants/:id", controllers.UpdateRestaurant)
		api.DELETE("/restaurants/:id", controllers.DeleteRestaurant)
	}
}
