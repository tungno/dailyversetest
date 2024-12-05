/**
 *  MockUserRepository is a mock implementation of the UserRepository interface.
 *  It is used for testing user-related functionalities without relying on a database.
 *
 *  @file       mock_user_repository.go
 *  @package    mocks
 *
 *  @methods
 *  - NewMockUserRepository(users)                           - Creates a new instance of MockUserRepository.
 *  - GetUserByEmail(ctx, email)                             - Simulates retrieving a user by email.
 *  - GetUserByUsername(ctx, username)                       - Simulates retrieving a user by username.
 *  - CreateUser(ctx, user)                                  - Simulates creating a new user.
 *  - UpdateUser(ctx, email, updates)                        - Simulates updating user details.
 *  - SearchUsersByUsername(ctx, query)                      - Simulates searching for users by username prefix.
 *
 *  @behaviors
 *  - All methods manipulate an in-memory map to mimic database behavior.
 *  - Ensures unique user email for `CreateUser`.
 *  - Supports partial updates for user fields such as OTP, password, and verification status.
 *
 *  @dependencies
 *  - models.User: Represents the structure of a user.
 *
 *  @example
 *  ```
 *  users := make(map[string]*models.User)
 *  repo := NewMockUserRepository(users)
 *  ctx := context.Background()
 *
 *  newUser := &models.User{
 *      Email:    "user@example.com",
 *      Username: "testuser",
 *      Password: "hashedpassword",
 *      IsVerified: false,
 *  }
 *  err := repo.CreateUser(ctx, newUser)
 *  if err != nil {
 *      log.Fatal("Failed to create user:", err)
 *  }
 *  ```
 *
 *  @errors
 *  - Returns errors when a user is not found or when a user already exists.
 *
 *  @authors
 *      - Aayush
 *      - Tung
 *      - Boss
 *      - Majd
 */

package mocks

import (
	"context"
	"fmt"
	"proh2052-group6/pkg/models"
	"strings"
	"time"
)

// MockUserRepository provides an in-memory implementation of the UserRepository interface.
type MockUserRepository struct {
	Users map[string]*models.User // In-memory store for user data.
}

// NewMockUserRepository initializes a new MockUserRepository instance.
func NewMockUserRepository(users map[string]*models.User) *MockUserRepository {
	return &MockUserRepository{Users: users}
}

// GetUserByEmail simulates retrieving a user by email.
func (mur *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	if user, exists := mur.Users[email]; exists {
		return user, nil
	}
	return nil, fmt.Errorf("user not found")
}

// GetUserByUsername simulates retrieving a user by username (case-insensitive).
func (mur *MockUserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	for _, user := range mur.Users {
		if strings.ToLower(user.Username) == strings.ToLower(username) {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

// CreateUser simulates adding a new user to the repository.
func (mur *MockUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	if _, exists := mur.Users[user.Email]; exists {
		return fmt.Errorf("user already exists")
	}
	mur.Users[user.Email] = user
	return nil
}

// UpdateUser simulates updating a user's details.
func (mur *MockUserRepository) UpdateUser(ctx context.Context, email string, updates map[string]interface{}) error {
	user, exists := mur.Users[email]
	if !exists {
		return fmt.Errorf("user not found")
	}
	// Apply updates
	if otp, ok := updates["OTP"]; ok {
		user.OTP = otp.(string)
	}
	if otpExpiresAt, ok := updates["OTPExpiresAt"]; ok {
		user.OTPExpiresAt = otpExpiresAt.(time.Time)
	}
	if isVerified, ok := updates["IsVerified"]; ok {
		user.IsVerified = isVerified.(bool)
	}
	if password, ok := updates["Password"]; ok {
		user.Password = password.(string)
	}
	return nil
}

// SearchUsersByUsername simulates searching for users by username prefix (case-insensitive).
func (mur *MockUserRepository) SearchUsersByUsername(ctx context.Context, query string) ([]*models.User, error) {
	var users []*models.User
	queryLower := strings.ToLower(query)
	for _, user := range mur.Users {
		usernameLower := strings.ToLower(user.Username)
		if strings.HasPrefix(usernameLower, queryLower) {
			users = append(users, user)
		}
	}
	return users, nil
}
