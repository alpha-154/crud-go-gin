package controllers

import (
	"net/http"
	"time"

	"github.com/alpha-154/crud-go-gin/internal/dto"
	"github.com/alpha-154/crud-go-gin/internal/helpers"
	"github.com/alpha-154/crud-go-gin/internal/models"
	"github.com/alpha-154/crud-go-gin/internal/services"
	"github.com/gin-gonic/gin"
)

// SignUp handles the request to sign up a new user
func SignUp(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := services.SignUp(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// SignIn handles the request to sign in a user
func SignIn(c *gin.Context) {
	var input dto.SignInInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, err := services.SignIn(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	//c.JSON(http.StatusOK, tokens)

	// Set new refresh token in HTTP-only cookie
	// c.SetCookie("refresh_token", tokens.RefreshToken,
	// 	int(time.Now().Add(time.Hour*24*7).Unix()), // Expires in 7 days
	// 	"/",
	// 	"yourdomain.com", // Domain
	// 	true,             // Secure (only sent over HTTPS)
	// 	true,             // HTTP-only (not accessible via JavaScript)
	// )

	c.SetCookie("refresh_token", tokens.RefreshToken,
		int(time.Now().Add(time.Hour*24*7).Unix()),
		"/",
		"",    // ✅ Set empty domain for localhost
		false, // ✅ Disable Secure for HTTP testing
		true,  // Keep HTTP-only enabled
	)

	c.JSON(http.StatusOK, gin.H{"access_token": tokens.AccessToken})
}

// GetUser retrieves a user by ID
func GetUser(c *gin.Context) {
	userID := c.Param("id")
	user, err := services.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetAllUsers retrieves all users
func GetAllUsers(c *gin.Context) {
	users, err := services.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// Logout handles the request to log out a user
func Logout(c *gin.Context) {
	userID := c.Param("user_id")

	// Remove refresh token from database
	err := services.InvalidateUserTokens(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "logout failed"})
		return
	}

	// Clear refresh token cookie
	//c.SetCookie("refresh_token", "", -1, "/", "yourdomain.com", true, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

// Refresh token endpoint
func RefreshToken(c *gin.Context) {
	// Get refresh token from HTTP-only cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token required"})
		return
	}

	// Validate refresh token
	claims, err := helpers.ValidateRefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	// Get user from database to verify token hasn't been revoked
	user, err := services.GetUserByID(claims["user_id"].(string))
	if err != nil || user.RefreshToken != refreshToken {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token revoked"})
		return
	}

	// Generate new token pair
	tokens, err := helpers.GenerateTokenPair(user.UserID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate tokens"})
		return
	}

	// Set new refresh token in HTTP-only cookie
	// c.SetCookie("refresh_token", tokens.RefreshToken,
	// 	int(time.Now().Add(time.Hour*24*7).Unix()),
	// 	"/",
	// 	"yourdomain.com",
	// 	true, // Secure
	// 	true, // HTTP only
	// )
	c.SetCookie("refresh_token", tokens.RefreshToken,
		int(time.Now().Add(time.Hour*24*7).Unix()),
		"/",
		"",
		false, // Secure
		true,  // HTTP only
	)

	c.JSON(http.StatusOK, gin.H{"access_token": tokens.AccessToken})
}
