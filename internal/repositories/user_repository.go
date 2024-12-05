/**
 *  UserRepository defines an interface for interacting with user data in the database.
 *  This interface standardizes CRUD operations and search functionality for user entities.
 *
 *  @interface UserRepository
 *
 *  @methods
 *  - GetUserByEmail(ctx, email)                 - Retrieves a user by their email address.
 *  - GetUserByUsername(ctx, username)           - Retrieves a user by their username.
 *  - CreateUser(ctx, user)                      - Creates a new user in the database.
 *  - UpdateUser(ctx, email, updates)            - Updates a user's data in the database.
 *  - SearchUsersByUsername(ctx, query)          - Searches for users by a username substring (prefix match, case-insensitive).
 *
 *  @behaviors
 *  - Allows extensibility for implementing user management across different database systems.
 *  - Standardizes operations for retrieving and updating user-related data.
 *
 *  @dependencies
 *  - context.Context: Used for propagating deadlines, cancellation signals, and other request-scoped values.
 *  - models.User: Struct representing the user entity.
 *
 *  @example
 *  ```
 *  // Create a new user
 *  user := &models.User{
 *      Email: "user@example.com",
 *      Username: "JohnDoe",
 *  }
 *  err := userRepo.CreateUser(ctx, user)
 *
 *  // Fetch a user by email
 *  user, err := userRepo.GetUserByEmail(ctx, "user@example.com")
 *  ```
 *
 *  @file      user_repository.go
 *  @project   DailyVerse
 *  @framework Database Agnostic (e.g., Firestore, SQL, etc.)
 */

package repositories

import (
	"context"
	"proh2052-group6/pkg/models"
)

// UserRepository defines the interface for user-related data operations.
type UserRepository interface {
	// GetUserByEmail retrieves a user by their email address.
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)

	// GetUserByUsername retrieves a user by their username.
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)

	// CreateUser creates a new user in the database.
	CreateUser(ctx context.Context, user *models.User) error

	// UpdateUser updates a user's data in the database with the provided key-value pairs.
	UpdateUser(ctx context.Context, email string, updates map[string]interface{}) error

	// SearchUsersByUsername searches for users whose usernames match the given query.
	// The search supports prefix matching and is case-insensitive.
	SearchUsersByUsername(ctx context.Context, query string) ([]*models.User, error)
}
