// internal/repositories/event_repository.go
package repositories

import (
	"context"
	"proh2052-group6/pkg/models"
)

type EventRepository interface {
	CreateEvent(ctx context.Context, event *models.Event) error
	GetEvent(ctx context.Context, userEmail, eventID string) (*models.Event, error)
	UpdateEvent(ctx context.Context, event *models.Event) error
	DeleteEvent(ctx context.Context, userEmail, eventID string) error
	GetAllEvents(ctx context.Context, userEmail string) ([]models.Event, error)
}
