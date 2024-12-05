/**
 *  FirestoreFriendRepository provides a concrete implementation of the FriendRepository
 *  interface, managing CRUD operations for friend requests and relationships using
 *  Firestore as the database.
 *
 *  @file       firestore_friend_repository.go
 *  @package    repositories
 *
 *  @properties
 *  - Client (*firestore.Client) - A Firestore client instance for database interactions.
 *
 *  @methods
 *  - NewFirestoreFriendRepository(client)                    - Initializes a new FirestoreFriendRepository with the Firestore client.
 *  - CreateFriendRequest(ctx, friend)                        - Creates a friend request document in Firestore.
 *  - GetFriendRequest(ctx, senderEmail, recipientEmail)      - Retrieves a specific friend request document.
 *  - UpdateFriendRequest(ctx, senderEmail, recipientEmail)   - Updates fields in an existing friend request document.
 *  - DeleteFriendRequest(ctx, senderEmail, recipientEmail)   - Deletes a specific friend request document.
 *  - GetFriends(ctx, userEmail)                              - Retrieves all friends for a user with an "accepted" status.
 *  - GetPendingFriendRequests(ctx, userEmail)                - Retrieves all pending friend requests for a user.
 *
 *  @behaviors
 *  - Ensures friend request documents are uniquely identified using a composite key: `<senderEmail>_<recipientEmail>`.
 *  - Allows querying both sent and received friend requests by filtering on `Email` or `FriendEmail` fields.
 *  - Supports updating only specific fields in friend request documents using Firestore's `MergeAll` option.
 *  - Handles Firestore errors gracefully, returning `nil` for `NotFound` errors in `GetFriendRequest`.
 *
 *  @examples
 *  Create a Friend Request:
 *  ```
 *  friend := &models.Friend{
 *      Email:       "user@example.com",
 *      FriendEmail: "friend@example.com",
 *      Status:      "pending",
 *  }
 *  err := repository.CreateFriendRequest(ctx, friend)
 *  ```
 *
 *  Get Friends:
 *  ```
 *  friends, err := repository.GetFriends(ctx, "user@example.com")
 *  if err != nil {
 *      log.Fatal(err)
 *  }
 *  for _, friend := range friends {
 *      fmt.Println(friend.FriendEmail)
 *  }
 *  ```
 *
 *  @dependencies
 *  - Firestore client: Provides CRUD and query operations for Firestore collections.
 *  - "google.golang.org/grpc/status": For handling Firestore-specific errors.
 *  - "google.golang.org/api/iterator": To iterate through Firestore query results.
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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"proh2052-group6/pkg/models"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// FirestoreFriendRepository manages friend-related operations in Firestore.
type FirestoreFriendRepository struct {
	Client *firestore.Client // Firestore client instance.
}

// NewFirestoreFriendRepository initializes a new FirestoreFriendRepository.
func NewFirestoreFriendRepository(client *firestore.Client) FriendRepository {
	return &FirestoreFriendRepository{Client: client}
}

// CreateFriendRequest creates a new friend request document in Firestore.
func (fr *FirestoreFriendRepository) CreateFriendRequest(ctx context.Context, friend *models.Friend) error {
	docID := friend.Email + "_" + friend.FriendEmail
	_, err := fr.Client.Collection("friends").Doc(docID).Set(ctx, friend)
	return err
}

// GetFriendRequest retrieves a specific friend request document by sender and recipient emails.
func (fr *FirestoreFriendRepository) GetFriendRequest(ctx context.Context, senderEmail, recipientEmail string) (*models.Friend, error) {
	docID := senderEmail + "_" + recipientEmail
	doc, err := fr.Client.Collection("friends").Doc(docID).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil // Return nil if document not found.
		}
		return nil, err
	}
	var friend models.Friend
	if err := doc.DataTo(&friend); err != nil {
		return nil, err
	}
	return &friend, nil
}

// UpdateFriendRequest updates specific fields in an existing friend request document.
func (fr *FirestoreFriendRepository) UpdateFriendRequest(ctx context.Context, senderEmail, recipientEmail string, updates map[string]interface{}) error {
	docID := senderEmail + "_" + recipientEmail
	_, err := fr.Client.Collection("friends").Doc(docID).Set(ctx, updates, firestore.MergeAll)
	return err
}

// DeleteFriendRequest deletes a specific friend request document from Firestore.
func (fr *FirestoreFriendRepository) DeleteFriendRequest(ctx context.Context, senderEmail, recipientEmail string) error {
	docID := senderEmail + "_" + recipientEmail
	_, err := fr.Client.Collection("friends").Doc(docID).Delete(ctx)
	return err
}

// GetFriends retrieves all accepted friends for a user.
func (fr *FirestoreFriendRepository) GetFriends(ctx context.Context, userEmail string) ([]models.Friend, error) {
	var friends []models.Friend

	// Query for friends where the user is the sender.
	iter := fr.Client.Collection("friends").Where("Email", "==", userEmail).Where("Status", "==", "accepted").Documents(ctx)
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var friend models.Friend
		if err := doc.DataTo(&friend); err != nil {
			continue
		}
		friends = append(friends, friend)
	}

	// Query for friends where the user is the recipient.
	iter = fr.Client.Collection("friends").Where("FriendEmail", "==", userEmail).Where("Status", "==", "accepted").Documents(ctx)
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var friend models.Friend
		if err := doc.DataTo(&friend); err != nil {
			continue
		}
		friends = append(friends, friend)
	}

	return friends, nil
}

// GetPendingFriendRequests fetches all pending friend requests for a user.
func (fr *FirestoreFriendRepository) GetPendingFriendRequests(ctx context.Context, userEmail string) ([]models.Friend, error) {
	var friends []models.Friend

	// Query where FriendEmail is userEmail and Status is "pending".
	iter := fr.Client.Collection("friends").Where("FriendEmail", "==", userEmail).Where("Status", "==", "pending").Documents(ctx)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var friend models.Friend
		if err := doc.DataTo(&friend); err != nil {
			continue
		}

		friends = append(friends, friend)
	}

	return friends, nil
}
