package services

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/alpha-154/crud-go-gin/internal/config"
	"github.com/alpha-154/crud-go-gin/internal/helpers"

	"github.com/alpha-154/crud-go-gin/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// SignUp registers a new user and generates the required tokens.
func SignUp(user models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if email exists
	emailTaken, err := helpers.IsEmailTaken(ctx, user.Email)
	if err != nil {
		return nil, fmt.Errorf("error checking email: %w", err)
	}
	if emailTaken {
		return nil, errors.New("email already exists")
	}

	// Hash password using the helper function
	hashedPassword, err := helpers.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	// Set user fields
	user.ID = primitive.NewObjectID()
	user.UserID = user.ID.Hex()
	user.Password = hashedPassword
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	if user.Role == "" {
		user.Role = "user"
	}

	// Generate refresh token using the helper function
	refreshTokenString, err := helpers.GenerateRefreshToken(user.UserID)
	if err != nil {
		return nil, err
	}
	user.RefreshToken = refreshTokenString

	// Save user to the database
	_, err = getUserCollection().InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	// Don't return the password in the response
	user.Password = ""

	// Return the user with the refresh token
	return &user, nil
}

func SignIn(input models.SignInInput) (*models.TokenResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := getUserCollection().FindOne(ctx, bson.M{"email": input.Email}).Decode(&user)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate access token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.UserID
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Hour).Unix()

	accessToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, err
	}

	// Generate new refresh token
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["user_id"] = user.UserID
	rtClaims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()

	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, err
	}

	// Update refresh token in database
	update := bson.M{"$set": bson.M{"refresh_token": refreshTokenString, "updated_at": time.Now()}}
	_, err = getUserCollection().UpdateOne(ctx, bson.M{"user_id": user.UserID}, update)
	if err != nil {
		return nil, err
	}

	return &models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString,
	}, nil
}

func GetUserByID(userID string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := getUserCollection().FindOne(ctx, bson.M{"user_id": userID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	user.Password = "" // Remove password from response
	return &user, nil
}

func GetAllUsers() ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := getUserCollection().Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	// Remove passwords from response
	for i := range users {
		users[i].Password = ""
	}

	return users, nil
}
