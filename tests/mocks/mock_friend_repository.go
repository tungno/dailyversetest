/**
 *  MockFriendRepository is a mock implementation of the FriendRepository interface.
 *  It is used for testing friend-related functionalities without relying on a database.
 *
 *  @file       mock_friend_repository.go
 *  @package    mocks
 *
 *  @methods
 *  - NewMockFriendRepository(friends)                              - Creates a new instance of MockFriendRepository.
 *  - CreateFriendRequest(ctx, friend)                              - Simulates creating a friend request.
 *  - GetFriendRequest(ctx, senderEmail, recipientEmail)            - Simulates fetching a friend request by sender and recipient emails.
 *  - UpdateFriendRequest(ctx, senderEmail, recipientEmail, updates) - Simulates updating a friend request's details.
 *  - DeleteFriendRequest(ctx, senderEmail, recipientEmail)          - Simulates deleting a friend request.
 *  - GetFriends(ctx, userEmail)                                    - Simulates retrieving all accepted friends for a user.
 *  - GetPendingFriendRequests(ctx, userEmail)                      - Simulates retrieving pending friend requests for a user.
 *
 *  @behaviors
 *  - All methods manipulate an in-memory map to mimic database behavior.
 *  - Friend requests are uniquely identified by a combination of sender and recipient email addresses.
 *  - Provides filtering for accepted and pending friend requests.
 *
 *  @dependencies
 *  - models.Friend: Represents the structure of a friend or friend request.
 *
 *  @example
 *  ```
 *  friends := make(map[string]*models.Friend)
 *  repo := NewMockFriendRepository(friends)
 *  ctx := context.Background()
 *
 *  friendRequest := &models.Friend{
 *      Email:        "user1@example.com",
 *      FriendEmail:  "user2@example.com",
 *      Status:       "pending",
 *  }
 *  err := repo.CreateFriendRequest(ctx, friendRequest)
 *  if err != nil {
 *      log.Fatal("Failed to create friend request:", err)
 *  }
 *  ```
 *
 *  @errors
 *  - Returns errors when friend requests are not found for fetching, updating, or deleting.
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
	"errors"
	"proh2052-group6/pkg/models"
)

// MockFriendRepository provides an in-memory implementation of the FriendRepository interface.
type MockFriendRepository struct {
	Friends map[string]*models.Friend // In-memory store for friend requests.
}

// NewMockFriendRepository initializes a new MockFriendRepository instance.
func NewMockFriendRepository(friends map[string]*models.Friend) *MockFriendRepository {
	return &MockFriendRepository{Friends: friends}
}

// CreateFriendRequest simulates creating a friend request.
func (mfr *MockFriendRepository) CreateFriendRequest(ctx context.Context, friend *models.Friend) error {
	docID := friend.Email + "_" + friend.FriendEmail
	mfr.Friends[docID] = friend
	return nil
}

// GetFriendRequest simulates retrieving a specific friend request by sender and recipient emails.
func (mfr *MockFriendRepository) GetFriendRequest(ctx context.Context, senderEmail, recipientEmail string) (*models.Friend, error) {
	docID := senderEmail + "_" + recipientEmail
	friend, exists := mfr.Friends[docID]
	if !exists {
		return nil, errors.New("friend request not found")
	}
	return friend, nil
}

// UpdateFriendRequest simulates updating the details of a specific friend request.
func (mfr *MockFriendRepository) UpdateFriendRequest(ctx context.Context, senderEmail, recipientEmail string, updates map[string]interface{}) error {
	docID := senderEmail + "_" + recipientEmail
	friend, exists := mfr.Friends[docID]
	if !exists {
		return errors.New("friend request not found")
	}
	if status, ok := updates["Status"].(string); ok {
		friend.Status = status
	}
	return nil
}

// DeleteFriendRequest simulates deleting a specific friend request.
func (mfr *MockFriendRepository) DeleteFriendRequest(ctx context.Context, senderEmail, recipientEmail string) error {
	docID := senderEmail + "_" + recipientEmail
	delete(mfr.Friends, docID)
	return nil
}

// GetFriends simulates retrieving all accepted friends for a given user.
func (mfr *MockFriendRepository) GetFriends(ctx context.Context, userEmail string) ([]models.Friend, error) {
	var friends []models.Friend
	for _, friend := range mfr.Friends {
		if (friend.Email == userEmail || friend.FriendEmail == userEmail) && friend.Status == "accepted" {
			friends = append(friends, *friend)
		}
	}
	return friends, nil
}

// GetPendingFriendRequests simulates retrieving all pending friend requests for a given user.
func (mfr *MockFriendRepository) GetPendingFriendRequests(ctx context.Context, userEmail string) ([]models.Friend, error) {
	var pendingRequests []models.Friend
	for _, friend := range mfr.Friends {
		if friend.FriendEmail == userEmail && friend.Status == "pending" {
			pendingRequests = append(pendingRequests, *friend)
		}
	}
	return pendingRequests, nil
}
