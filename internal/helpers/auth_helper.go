package helpers

import (
	"context"
	"os"
	"time"

	"github.com/alpha-154/crud-go-gin/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection

func getUserCollection() *mongo.Collection {
	if userCollection == nil {
		userCollection = config.GetCollection(config.DB, "users")
	}
	return userCollection
}

// IsEmailTaken checks if the provided email is already taken in the database.
func IsEmailTaken(ctx context.Context, email string) (bool, error) {
	count, err := getUserCollection().CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// HashPassword hashes the provided password using bcrypt.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// GenerateRefreshToken generates a new JWT refresh token for the user.
func GenerateRefreshToken(userID string) (string, error) {
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["user_id"] = userID
	rtClaims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix() // Set expiration to 7 days

	// Sign the token with the secret key from environment variables
	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return refreshTokenString, nil
}
