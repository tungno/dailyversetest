// tests/mocks/mock_friend_repository.go
package mocks

import (
	"context"
	"errors"
	"proh2052-group6/pkg/models"
)

type MockFriendRepository struct {
	Friends map[string]*models.Friend
}

func NewMockFriendRepository(friends map[string]*models.Friend) *MockFriendRepository {
	return &MockFriendRepository{Friends: friends}
}

func (mfr *MockFriendRepository) CreateFriendRequest(ctx context.Context, friend *models.Friend) error {
	docID := friend.Email + "_" + friend.FriendEmail
	mfr.Friends[docID] = friend
	return nil
}

func (mfr *MockFriendRepository) GetFriendRequest(ctx context.Context, senderEmail, recipientEmail string) (*models.Friend, error) {
	docID := senderEmail + "_" + recipientEmail
	friend, exists := mfr.Friends[docID]
	if !exists {
		return nil, errors.New("friend request not found")
	}
	return friend, nil
}

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

func (mfr *MockFriendRepository) DeleteFriendRequest(ctx context.Context, senderEmail, recipientEmail string) error {
	docID := senderEmail + "_" + recipientEmail
	delete(mfr.Friends, docID)
	return nil
}

func (mfr *MockFriendRepository) GetFriends(ctx context.Context, userEmail string) ([]models.Friend, error) {
	var friends []models.Friend
	for _, friend := range mfr.Friends {
		if (friend.Email == userEmail || friend.FriendEmail == userEmail) && friend.Status == "accepted" {
			friends = append(friends, *friend)
		}
	}
	return friends, nil
}

func (mfr *MockFriendRepository) GetPendingFriendRequests(ctx context.Context, userEmail string) ([]models.Friend, error) {
	var pendingRequests []models.Friend
	for _, friend := range mfr.Friends {
		if friend.FriendEmail == userEmail && friend.Status == "pending" {
			pendingRequests = append(pendingRequests, *friend)
		}
	}
	return pendingRequests, nil
}
