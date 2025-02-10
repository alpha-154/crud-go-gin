package routes

import (
	"github.com/alpha-154/crud-go-gin/internal/controllers"
	"github.com/alpha-154/crud-go-gin/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")

	// Auth routes
	auth := api.Group("/auth")
	{
		auth.POST("/signup", controllers.SignUp)
		auth.POST("/signin", controllers.SignIn)
	}

	// Protected routes
	protected := api.Group("")
	protected.Use(middlewares.AuthMiddleware())
	{
		// User routes
		protected.GET("/users/:id", controllers.GetUser)
		protected.GET("/users", middlewares.AdminOnly(), controllers.GetAllUsers)

		// Restaurant routes
		protected.POST("/restaurants", controllers.CreateRestaurant)
		protected.GET("/restaurants", controllers.GetAllRestaurants)
		protected.GET("/restaurants/:id", controllers.GetRestaurant)
		protected.PUT("/restaurants/:id", controllers.UpdateRestaurant)
		protected.DELETE("/restaurants/:id", controllers.DeleteRestaurant)
	}
}
