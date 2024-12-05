/**
 *  EventRepository defines the interface for data access operations related to events.
 *  It abstracts the database layer, allowing the application to interact with event data
 *  without being tied to a specific database implementation.
 *
 *  @interface EventRepository
 *  @inherits None
 *
 *  @methods
 *  - CreateEvent(ctx, event)                - Creates a new event in the database.
 *  - GetEvent(ctx, userEmail, eventID)      - Retrieves a specific event by its ID and the user's email.
 *  - UpdateEvent(ctx, event)                - Updates an existing event in the database.
 *  - DeleteEvent(ctx, userEmail, eventID)   - Deletes an event by its ID and the user's email.
 *  - GetAllEvents(ctx, userEmail)           - Fetches all events associated with a specific user.
 *
 *  @dependencies
 *  - models.Event: Defines the structure of an event object.
 *  - context.Context: Used for managing request-scoped values, deadlines, and cancellation signals.
 *
 *  @file      event_repository.go
 *  @project   DailyVerse
 *  @framework Go Interface for Repository Pattern
 *  @authors
 *      - Aayush
 *      - Tung
 *      - Boss
 *      - Majd
 */

package repositories

import (
	"context"
	"proh2052-group6/pkg/models"
)

// EventRepository defines the interface for event-related data operations.
type EventRepository interface {
	// CreateEvent inserts a new event into the database.
	CreateEvent(ctx context.Context, event *models.Event) error

	// GetEvent retrieves a specific event by its ID and the associated user's email.
	GetEvent(ctx context.Context, userEmail, eventID string) (*models.Event, error)

	// UpdateEvent updates an existing event in the database.
	UpdateEvent(ctx context.Context, event *models.Event) error

	// DeleteEvent removes an event from the database by its ID and the user's email.
	DeleteEvent(ctx context.Context, userEmail, eventID string) error

	// GetAllEvents fetches all events associated with a specific user's email.
	GetAllEvents(ctx context.Context, userEmail string) ([]models.Event, error)
}
