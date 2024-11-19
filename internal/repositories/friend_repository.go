// internal/repositories/friend_repository.go
package repositories

import (
	"context"
	"proh2052-group6/pkg/models"
)

type FriendRepository interface {
	CreateFriendRequest(ctx context.Context, friend *models.Friend) error
	GetFriendRequest(ctx context.Context, senderEmail, recipientEmail string) (*models.Friend, error)
	UpdateFriendRequest(ctx context.Context, senderEmail, recipientEmail string, updates map[string]interface{}) error
	DeleteFriendRequest(ctx context.Context, senderEmail, recipientEmail string) error
	GetFriends(ctx context.Context, userEmail string) ([]models.Friend, error)
	GetPendingFriendRequests(ctx context.Context, userEmail string) ([]models.Friend, error)
}
