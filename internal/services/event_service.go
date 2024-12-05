/**
 *  EventService provides business logic for managing events. It acts as an intermediary
 *  between the repositories and handlers, ensuring proper validation and formatting of event data.
 *
 *  @interface EventServiceInterface
 *  @methods
 *  - CreateEvent(ctx, event)                  - Creates a new event with validation.
 *  - GetEvent(ctx, userEmail, eventID)        - Retrieves a specific event by its ID.
 *  - UpdateEvent(ctx, event)                  - Updates an existing event.
 *  - DeleteEvent(ctx, userEmail, eventID)     - Deletes a specific event by its ID.
 *  - GetAllEvents(ctx, userEmail)             - Retrieves all events for a given user.
 *
 *  @struct   EventService
 *  @inherits EventServiceInterface
 *
 *  @methods
 *  - NewEventService(eventRepo)              - Initializes a new EventService with the given repository.
 *  - CreateEvent(ctx, event)                 - Implements event creation logic.
 *  - GetEvent(ctx, userEmail, eventID)       - Implements event retrieval logic.
 *  - UpdateEvent(ctx, event)                 - Implements event update logic.
 *  - DeleteEvent(ctx, userEmail, eventID)    - Implements event deletion logic.
 *  - GetAllEvents(ctx, userEmail)            - Implements logic to retrieve all events for a user.
 *
 *  @behaviors
 *  - Validates event data (e.g., EventTypeID, Date format) before creating an event.
 *  - Ensures only authorized users can access or modify their events.
 *  - Handles errors gracefully and returns meaningful messages on failure.
 *
 *  @dependencies
 *  - repositories.EventRepository: Repository for interacting with event data in the database.
 *  - models.Event: Struct representing the event entity.
 *
 *  @example
 *  ```
 *  // Create a new event
 *  event := &models.Event{
 *      Email: "user@example.com",
 *      Title: "Meeting",
 *      Date: "2024-12-01",
 *      EventTypeID: "public",
 *  }
 *  err := eventService.CreateEvent(ctx, event)
 *
 *  // Fetch an event
 *  event, err := eventService.GetEvent(ctx, "user@example.com", "eventID123")
 *  ```
 *
 *  @file      event_service.go
 *  @project   DailyVerse
 *  @framework Go HTTP Server & Firestore API
 */

package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"proh2052-group6/internal/repositories"
	"proh2052-group6/pkg/models"
)

// EventServiceInterface defines methods for managing events.
type EventServiceInterface interface {
	CreateEvent(ctx context.Context, event *models.Event) error
	GetEvent(ctx context.Context, userEmail, eventID string) (*models.Event, error)
	UpdateEvent(ctx context.Context, event *models.Event) error
	DeleteEvent(ctx context.Context, userEmail, eventID string) error
	GetAllEvents(ctx context.Context, userEmail string) ([]models.Event, error)
}

// EventService provides implementations for EventServiceInterface.
type EventService struct {
	EventRepo repositories.EventRepository
}

// NewEventService initializes a new EventService with the given EventRepository.
func NewEventService(eventRepo repositories.EventRepository) EventServiceInterface {
	return &EventService{EventRepo: eventRepo}
}

// CreateEvent validates and creates a new event.
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

// GetEvent retrieves a specific event by its ID and ensures the user is authorized to access it.
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

// UpdateEvent updates an existing event in the repository.
func (es *EventService) UpdateEvent(ctx context.Context, event *models.Event) error {
	return es.EventRepo.UpdateEvent(ctx, event)
}

// DeleteEvent deletes a specific event by its ID for a user.
func (es *EventService) DeleteEvent(ctx context.Context, userEmail, eventID string) error {
	return es.EventRepo.DeleteEvent(ctx, userEmail, eventID)
}

// GetAllEvents retrieves all events for a specific user from the repository.
func (es *EventService) GetAllEvents(ctx context.Context, userEmail string) ([]models.Event, error) {
	return es.EventRepo.GetAllEvents(ctx, userEmail)
}
