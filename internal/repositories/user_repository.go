// internal/repositories/user_repository.go
package repositories

import (
	"context"
	"proh2052-group6/pkg/models"
)

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
	UpdateUser(ctx context.Context, email string, updates map[string]interface{}) error
	SearchUsersByUsername(ctx context.Context, query string) ([]*models.User, error) // Added method
}
