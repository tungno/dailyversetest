/**
 *  MockEventService simulates an event service for testing purposes.
 *  It provides in-memory operations to manage events, including creation, retrieval,
 *  updating, and deletion. This mock implementation allows testing of handlers
 *  and services without requiring an actual database.
 *
 *  @file       mock_event_service.go
 *  @package    mocks
 *
 *  @structs
 *  - MockEventService: Simulates an event service with an in-memory store for events.
 *
 *  @methods
 *  - NewMockEventService: Initializes a new instance of MockEventService.
 *  - CreateEvent(ctx, event): Simulates creating a new event.
 *  - GetEvent(ctx, userEmail, eventID): Simulates retrieving an event by ID and user email.
 *  - UpdateEvent(ctx, event): Simulates updating an event.
 *  - DeleteEvent(ctx, userEmail, eventID): Simulates deleting an event.
 *  - GetAllEvents(ctx, userEmail): Simulates retrieving all events for a user.
 *
 *  @example
 *  ```
 *  mockService := mocks.NewMockEventService()
 *  event := &models.Event{
 *      EventID: "1",
 *      Email:   "user@example.com",
 *      Title:   "Sample Event",
 *  }
 *
 *  err := mockService.CreateEvent(context.Background(), event)
 *  if err != nil {
 *      t.Fatalf("Failed to create event: %v", err)
 *  }
 *  ```
 *
 *  @dependencies
 *  - pkg/models: Provides the Event model for use in the mock service.
 *
 *  @limitations
 *  - MockEventService is in-memory and does not persist data across tests.
 *  - MockEventService does not implement complex querying capabilities.
 *
 *  @authors
 *      - Aayush
 *      - Tung
 *      - Boss
 *      - Majd
 */

package mocks

import (
	"context"
	"fmt"
	"proh2052-group6/pkg/models"
)

// MockEventService simulates an event service for testing.
type MockEventService struct {
	Events map[string]*models.Event // In-memory store for events.
}

// NewMockEventService initializes a new instance of MockEventService.
func NewMockEventService() *MockEventService {
	return &MockEventService{
		Events: make(map[string]*models.Event),
	}
}

// CreateEvent simulates creating a new event.
func (mes *MockEventService) CreateEvent(ctx context.Context, event *models.Event) error {
	if _, exists := mes.Events[event.EventID]; exists {
		return fmt.Errorf("event already exists")
	}
	mes.Events[event.EventID] = event
	return nil
}

// GetEvent simulates retrieving an event by ID and user email.
func (mes *MockEventService) GetEvent(ctx context.Context, userEmail, eventID string) (*models.Event, error) {
	event, exists := mes.Events[eventID]
	if !exists || event.Email != userEmail {
		return nil, fmt.Errorf("event not found")
	}
	return event, nil
}

// UpdateEvent simulates updating an existing event.
func (mes *MockEventService) UpdateEvent(ctx context.Context, event *models.Event) error {
	existingEvent, exists := mes.Events[event.EventID]
	if !exists || existingEvent.Email != event.Email {
		return fmt.Errorf("event not found")
	}
	mes.Events[event.EventID] = event
	return nil
}

// DeleteEvent simulates deleting an event by ID and user email.
func (mes *MockEventService) DeleteEvent(ctx context.Context, userEmail, eventID string) error {
	event, exists := mes.Events[eventID]
	if !exists || event.Email != userEmail {
		return fmt.Errorf("event not found")
	}
	delete(mes.Events, eventID)
	return nil
}

// GetAllEvents simulates retrieving all events for a specific user.
func (mes *MockEventService) GetAllEvents(ctx context.Context, userEmail string) ([]models.Event, error) {
	var events []models.Event
	for _, event := range mes.Events {
		if event.Email == userEmail {
			events = append(events, *event)
		}
	}
	return events, nil
}
