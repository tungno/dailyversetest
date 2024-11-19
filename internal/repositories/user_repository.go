// internal/repositories/user_repository.go
package repositories

import (
	"context"
	"proh2052-group6/pkg/models"
)

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
}
