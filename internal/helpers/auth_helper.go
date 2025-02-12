package helpers

import (
	"context"
	"errors"

	"os"
	"time"

	"github.com/alpha-154/crud-go-gin/internal/config"
	"github.com/alpha-154/crud-go-gin/internal/dto"

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

// Generate AccessToken & RefreshToken for the user
func GenerateTokenPair(userID string, role string) (*dto.TokenPair, error) {
	// Generate access token
	accessToken := jwt.New(jwt.SigningMethodHS256)
	accessClaims := accessToken.Claims.(jwt.MapClaims)
	accessClaims["user_id"] = userID
	accessClaims["role"] = role
	accessClaims["exp"] = time.Now().Add(time.Hour).Unix()
	accessClaims["type"] = "access"

	// Generate refresh token
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshClaims["user_id"] = userID
	refreshClaims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix() // Set expiration to 7 days

	refreshClaims["type"] = "refresh"

	// Sign tokens
	accessTokenString, err := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, err
	}

	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, err
	}

	return &dto.TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

// ValidateRefreshToken validates the provided refresh token and returns the claims if valid.
func ValidateRefreshToken(tokenString string) (map[string]interface{}, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
