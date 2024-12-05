/**
 *  FriendRepository defines the interface for managing friend-related operations,
 *  including creating, retrieving, updating, and deleting friend requests, as well
 *  as fetching friends and pending friend requests.
 *
 *  @file       friend_repository.go
 *  @package    repositories
 *
 *  @methods
 *  - CreateFriendRequest(ctx, friend)                   - Creates a new friend request.
 *  - GetFriendRequest(ctx, senderEmail, recipientEmail) - Retrieves a specific friend request.
 *  - UpdateFriendRequest(ctx, senderEmail, recipientEmail, updates) - Updates fields of an existing friend request.
 *  - DeleteFriendRequest(ctx, senderEmail, recipientEmail) - Deletes a specific friend request.
 *  - GetFriends(ctx, userEmail)                         - Fetches all friends for a user with the "accepted" status.
 *  - GetPendingFriendRequests(ctx, userEmail)           - Fetches all pending friend requests for a user.
 *
 *  @behavior
 *  - Provides a contract for repository implementations to ensure consistency.
 *  - Focuses on operations for friend requests and relationships.
 *
 *  @example
 *  ```
 *  type FirestoreFriendRepository struct {
 *      Client *firestore.Client
 *  }
 *
 *  func (fr *FirestoreFriendRepository) CreateFriendRequest(ctx context.Context, friend *models.Friend) error {
 *      // Firestore implementation here
 *  }
 *
 *  func (fr *FirestoreFriendRepository) GetFriendRequest(ctx context.Context, senderEmail, recipientEmail string) (*models.Friend, error) {
 *      // Firestore implementation here
 *  }
 *  ```
 *
 *  @dependencies
 *  - Context: For passing request-scoped values and managing timeouts or deadlines.
 *  - models.Friend: Represents the data structure for friend requests.
 *
 *  @authors
 *      - Aayush
 *      - Tung
 *      - Boss
 *      - Majd
 */

package repositories

import (
	"context"
	"proh2052-group6/pkg/models"
)

// FriendRepository defines the interface for friend-related operations.
type FriendRepository interface {
	// CreateFriendRequest creates a new friend request.
	CreateFriendRequest(ctx context.Context, friend *models.Friend) error

	// GetFriendRequest retrieves a specific friend request based on sender and recipient emails.
	GetFriendRequest(ctx context.Context, senderEmail, recipientEmail string) (*models.Friend, error)

	// UpdateFriendRequest updates specific fields in an existing friend request.
	UpdateFriendRequest(ctx context.Context, senderEmail, recipientEmail string, updates map[string]interface{}) error

	// DeleteFriendRequest deletes a specific friend request.
	DeleteFriendRequest(ctx context.Context, senderEmail, recipientEmail string) error

	// GetFriends retrieves all friends for a user with the "accepted" status.
	GetFriends(ctx context.Context, userEmail string) ([]models.Friend, error)

	// GetPendingFriendRequests retrieves all pending friend requests for a user.
	GetPendingFriendRequests(ctx context.Context, userEmail string) ([]models.Friend, error)
}
