// tests/services/friend_service_test.go
package services_test

import (
	"context"
	"testing"

	"proh2052-group6/internal/services"
	"proh2052-group6/pkg/models"
	"proh2052-group6/tests/mocks"
)

func TestSendFriendRequest(t *testing.T) {
	mockDB := &mocks.MockDB{
		Users: map[string]*models.User{
			"user1@example.com": {Email: "user1@example.com", Username: "user1"},
			"user2@example.com": {Email: "user2@example.com", Username: "user2"},
		},
		Friends: make(map[string]*models.Friend),
	}
	friendService := services.NewFriendService(mockDB)

	err := friendService.SendFriendRequest(context.Background(), "user1@example.com", "user2")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(mockDB.Friends) != 1 {
		t.Errorf("Expected 1 friend request, got %d", len(mockDB.Friends))
	}
}

// Add similar tests for other methods like AcceptFriendRequest, GetFriendsList, etc.
