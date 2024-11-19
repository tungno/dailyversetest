// tests/mocks/mock_friend_service.go
package mocks

import (
	"context"

	"proh2052-group6/pkg/models"
)

type MockFriendService struct {
	// Add fields to simulate service behavior
}

func (mfs *MockFriendService) SendFriendRequest(ctx context.Context, userEmail, username string) error {
	// Simulate sending friend request
	return nil
}

func (mfs *MockFriendService) AcceptFriendRequest(ctx context.Context, userEmail, username string) error {
	// Simulate accepting friend request
	return nil
}

func (mfs *MockFriendService) GetFriendsList(ctx context.Context, userEmail string) ([]models.User, error) {
	// Simulate retrieving friends list
	return []models.User{}, nil
}

func (mfs *MockFriendService) RemoveFriend(ctx context.Context, userEmail, username string) error {
	// Simulate removing friend
	return nil
}

func (mfs *MockFriendService) GetPendingFriendRequests(ctx context.Context, userEmail string) ([]models.User, error) {
	// Simulate getting pending friend requests
	return []models.User{}, nil
}

func (mfs *MockFriendService) DeclineFriendRequest(ctx context.Context, userEmail, username string) error {
	// Simulate declining friend request
	return nil
}

func (mfs *MockFriendService) CancelFriendRequest(ctx context.Context, userEmail, username string) error {
	// Simulate canceling friend request
	return nil
}
