package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/alpha-154/crud-go-gin/internal/config"
	"github.com/alpha-154/crud-go-gin/internal/helpers"

	"github.com/alpha-154/crud-go-gin/internal/dto"
	"github.com/alpha-154/crud-go-gin/internal/models"
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

// SignIn authenticates a user and generates the required tokens
func SignIn(input dto.SignInInput) (*dto.SignInServiceResponse, error) {
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

	// Generate access & refresh token
	token, err := helpers.GenerateTokenPair(user.UserID, user.Role)
	if err != nil {
		return nil, err
	}

	// Update refresh token in database
	updateResult, err := getUserCollection().UpdateOne(
		ctx,
		bson.M{"user_id": user.UserID}, // Ensure user.UserID is a string
		bson.M{"$set": bson.M{"refresh_token": token.RefreshToken}},
	)
	if err != nil {
		return nil, err
	}

	fmt.Println("Matched count:", updateResult.MatchedCount, "Modified count:", updateResult.ModifiedCount)

	// If no document was modified, return an error
	if updateResult.MatchedCount == 0 {
		return nil, errors.New("failed to update refresh token")
	}

	return &dto.SignInServiceResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}

// GetUserByID retrieves a user by ID
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

// GetAllUsers retrieves all users
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

// InvalidateUserTokens invalidates the refresh token for a user
func InvalidateUserTokens(userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := getUserCollection().UpdateOne(ctx, bson.M{"user_id": userID}, bson.M{"$set": bson.M{"refresh_token": ""}})
	if err != nil {
		return err
	}
	return nil
}
