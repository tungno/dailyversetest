/**
 *  MockFriendService provides a mock implementation of the FriendServiceInterface for testing purposes.
 *  This mock enables simulating friend-related operations, such as sending requests, accepting friends,
 *  and retrieving friend lists without using the actual service implementation.
 *
 *  @struct   MockFriendService
 *  @inherits FriendServiceInterface
 *
 *  @methods
 *  - SendFriendRequest(ctx, userEmail, username) (error): Simulates sending a friend request.
 *  - AcceptFriendRequest(ctx, userEmail, username) (error): Simulates accepting a friend request.
 *  - GetFriendsList(ctx, userEmail) ([]models.User, error): Simulates retrieving the user's friends list.
 *  - RemoveFriend(ctx, userEmail, username) (error): Simulates removing a friend.
 *  - GetPendingFriendRequests(ctx, userEmail) ([]models.User, error): Simulates retrieving pending friend requests.
 *  - DeclineFriendRequest(ctx, userEmail, username) (error): Simulates declining a friend request.
 *  - CancelFriendRequest(ctx, userEmail, username) (error): Simulates canceling a friend request.
 *
 *  @example
 *  ```
 *  // Initialize the mock friend service
 *  mockFriendService := &MockFriendService{}
 *
 *  // Simulate sending a friend request
 *  err := mockFriendService.SendFriendRequest(context.Background(), "user1@example.com", "user2")
 *  if err != nil {
 *      t.Errorf("Expected no error, got %v", err)
 *  }
 *
 *  // Simulate retrieving the user's friends list
 *  friends, err := mockFriendService.GetFriendsList(context.Background(), "user1@example.com")
 *  if err != nil {
 *      t.Errorf("Expected no error, got %v", err)
 *  }
 *  fmt.Println(friends) // Output: []
 *  ```
 *
 *  @file      mock_friend_service.go
 *  @project   DailyVerse
 *  @framework Go Testing with Mock Services
 */

package mocks

import (
	"context"
	"proh2052-group6/pkg/models"
)

// MockFriendService is a mock implementation of the FriendServiceInterface.
type MockFriendService struct {
	// Add fields to simulate service behavior, e.g., store friend requests, relationships, etc.
}

// SendFriendRequest simulates sending a friend request.
// Parameters:
// - ctx (context.Context): The request context.
// - userEmail (string): The email of the user sending the request.
// - username (string): The username of the user to whom the request is being sent.
//
// Returns:
// - error: Always returns nil in this mock, simulating successful request sending.
func (mfs *MockFriendService) SendFriendRequest(ctx context.Context, userEmail, username string) error {
	// Simulate sending friend request
	return nil
}

// AcceptFriendRequest simulates accepting a friend request.
// Parameters:
// - ctx (context.Context): The request context.
// - userEmail (string): The email of the user accepting the request.
// - username (string): The username of the friend being accepted.
//
// Returns:
// - error: Always returns nil in this mock, simulating successful request acceptance.
func (mfs *MockFriendService) AcceptFriendRequest(ctx context.Context, userEmail, username string) error {
	// Simulate accepting friend request
	return nil
}

// GetFriendsList simulates retrieving the user's friends list.
// Parameters:
// - ctx (context.Context): The request context.
// - userEmail (string): The email of the user whose friends list is being requested.
//
// Returns:
// - []models.User: A slice of users representing the friends list.
// - error: Always returns nil in this mock.
func (mfs *MockFriendService) GetFriendsList(ctx context.Context, userEmail string) ([]models.User, error) {
	// Simulate retrieving friends list
	return []models.User{}, nil
}

// RemoveFriend simulates removing a friend.
// Parameters:
// - ctx (context.Context): The request context.
// - userEmail (string): The email of the user removing the friend.
// - username (string): The username of the friend being removed.
//
// Returns:
// - error: Always returns nil in this mock, simulating successful removal.
func (mfs *MockFriendService) RemoveFriend(ctx context.Context, userEmail, username string) error {
	// Simulate removing friend
	return nil
}

// GetPendingFriendRequests simulates retrieving pending friend requests.
// Parameters:
// - ctx (context.Context): The request context.
// - userEmail (string): The email of the user whose pending friend requests are being retrieved.
//
// Returns:
// - []models.User: A slice of users representing the pending friend requests.
// - error: Always returns nil in this mock.
func (mfs *MockFriendService) GetPendingFriendRequests(ctx context.Context, userEmail string) ([]models.User, error) {
	// Simulate getting pending friend requests
	return []models.User{}, nil
}

// DeclineFriendRequest simulates declining a friend request.
// Parameters:
// - ctx (context.Context): The request context.
// - userEmail (string): The email of the user declining the request.
// - username (string): The username of the friend request being declined.
//
// Returns:
// - error: Always returns nil in this mock, simulating successful decline.
func (mfs *MockFriendService) DeclineFriendRequest(ctx context.Context, userEmail, username string) error {
	// Simulate declining friend request
	return nil
}

// CancelFriendRequest simulates canceling a friend request.
// Parameters:
// - ctx (context.Context): The request context.
// - userEmail (string): The email of the user canceling the request.
// - username (string): The username of the friend request being canceled.
//
// Returns:
// - error: Always returns nil in this mock, simulating successful cancellation.
func (mfs *MockFriendService) CancelFriendRequest(ctx context.Context, userEmail, username string) error {
	// Simulate canceling friend request
	return nil
}
