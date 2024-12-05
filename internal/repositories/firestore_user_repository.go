/**
 *  FirestoreUserRepository provides methods to interact with the Firestore database for user-related operations.
 *  This repository encapsulates CRUD operations for managing user accounts and searching users.
 *
 *  @struct   FirestoreUserRepository
 *  @inherits None
 *
 *  @methods
 *  - NewFirestoreUserRepository(client)    - Initializes a new FirestoreUserRepository with a Firestore client.
 *  - GetUserByEmail(ctx, email)            - Fetches a user by their email address.
 *  - GetUserByUsername(ctx, username)      - Fetches a user by their username.
 *  - CreateUser(ctx, user)                 - Creates a new user in Firestore.
 *  - UpdateUser(ctx, email, updates)       - Updates a user's details in Firestore.
 *  - SearchUsersByUsername(ctx, query)     - Searches users by a username substring query.
 *
 *  @behaviors
 *  - Uses Firestore's document-based structure to store and query user data under `users/{email}`.
 *  - Supports case-insensitive username search with prefix matching using Firestore queries.
 *  - Handles error scenarios and returns meaningful messages for failed operations.
 *
 *  @dependencies
 *  - cloud.google.com/go/firestore: Firestore client for database operations.
 *  - google.golang.org/api/iterator: Iterator for traversing Firestore query results.
 *  - models.User: Struct representing user data.
 *
 *  @example
 *  ```
 *  // Fetch a user by email
 *  user, err := repository.GetUserByEmail(ctx, "user@example.com")
 *
 *  // Create a new user
 *  user := &models.User{
 *      Email: "user@example.com",
 *      Username: "JohnDoe",
 *      UsernameLower: "johndoe",
 *  }
 *  err := repository.CreateUser(ctx, user)
 *  ```
 *
 *  @file      firestore_user_repository.go
 *  @project   DailyVerse
 *  @framework Firestore Client (Go) API
 *  @authors
 *      - Aayush
 *      - Tung
 *      - Boss
 *      - Majd
 */

package repositories

import (
	"context"
	"fmt"
	"proh2052-group6/pkg/models"
	"strings"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// FirestoreUserRepository implements the UserRepository interface for Firestore.
type FirestoreUserRepository struct {
	Client *firestore.Client
}

// NewFirestoreUserRepository initializes a new FirestoreUserRepository with the given Firestore client.
func NewFirestoreUserRepository(client *firestore.Client) UserRepository {
	return &FirestoreUserRepository{Client: client}
}

// GetUserByEmail retrieves a user by their email address.
func (ur *FirestoreUserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	doc, err := ur.Client.Collection("users").Doc(email).Get(ctx)
	if err != nil {
		return nil, err
	}
	var user models.User
	if err := doc.DataTo(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername retrieves a user by their username.
func (ur *FirestoreUserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	iter := ur.Client.Collection("users").Where("UsernameLower", "==", strings.ToLower(username)).Limit(1).Documents(ctx)
	defer iter.Stop()

	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}

	var user models.User
	if err := doc.DataTo(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser creates a new user in Firestore.
func (ur *FirestoreUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	_, err := ur.Client.Collection("users").Doc(user.Email).Set(ctx, user)
	return err
}

// UpdateUser updates a user's details in Firestore with the provided key-value pairs.
func (ur *FirestoreUserRepository) UpdateUser(ctx context.Context, email string, updates map[string]interface{}) error {
	_, err := ur.Client.Collection("users").Doc(email).Set(ctx, updates, firestore.MergeAll)
	return err
}

// SearchUsersByUsername searches for users with a username matching the given query (prefix match, case-insensitive).
func (ur *FirestoreUserRepository) SearchUsersByUsername(ctx context.Context, query string) ([]*models.User, error) {
	iter := ur.Client.Collection("users").
		Where("UsernameLower", ">=", strings.ToLower(query)).
		Where("UsernameLower", "<=", strings.ToLower(query)+"\uf8ff").
		Documents(ctx)
	defer iter.Stop()

	var users []*models.User
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var user models.User
		if err := doc.DataTo(&user); err != nil {
			continue
		}
		users = append(users, &user)
	}

	return users, nil
}
