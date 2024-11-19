// internal/services/event_service.go
package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"proh2052-group6/internal/repositories"
	"proh2052-group6/pkg/models"
)

type EventServiceInterface interface {
	CreateEvent(ctx context.Context, event *models.Event) error
	GetEvent(ctx context.Context, userEmail, eventID string) (*models.Event, error)
	UpdateEvent(ctx context.Context, event *models.Event) error
	DeleteEvent(ctx context.Context, userEmail, eventID string) error
	GetAllEvents(ctx context.Context, userEmail string) ([]models.Event, error)
}

type EventService struct {
	EventRepo repositories.EventRepository
}

func NewEventService(eventRepo repositories.EventRepository) EventServiceInterface {
	return &EventService{EventRepo: eventRepo}
}

func (es *EventService) CreateEvent(ctx context.Context, event *models.Event) error {
	// Validate EventTypeID
	event.EventTypeID = strings.ToLower(event.EventTypeID)
	if event.EventTypeID != "public" && event.EventTypeID != "private" {
		return fmt.Errorf("Invalid event type")
	}

	// Parse and format the date
	eventDate, err := time.Parse("2006-01-02", event.Date)
	if err != nil {
		return fmt.Errorf("Invalid date format. Please use YYYY-MM-DD.")
	}
	event.Date = eventDate.Format("2006-01-02")

	// Delegate to repository
	return es.EventRepo.CreateEvent(ctx, event)
}

func (es *EventService) GetEvent(ctx context.Context, userEmail, eventID string) (*models.Event, error) {
	event, err := es.EventRepo.GetEvent(ctx, userEmail, eventID)
	if err != nil {
		return nil, err
	}

	if event.Email != userEmail {
		return nil, fmt.Errorf("Unauthorized to access this event")
	}

	return event, nil
}

func (es *EventService) UpdateEvent(ctx context.Context, event *models.Event) error {
	return es.EventRepo.UpdateEvent(ctx, event)
}

func (es *EventService) DeleteEvent(ctx context.Context, userEmail, eventID string) error {
	return es.EventRepo.DeleteEvent(ctx, userEmail, eventID)
}

func (es *EventService) GetAllEvents(ctx context.Context, userEmail string) ([]models.Event, error) {
	return es.EventRepo.GetAllEvents(ctx, userEmail)
}
