// tests/mocks/mock_event_service.go
package mocks

import (
	"context"
	"fmt"
	"proh2052-group6/pkg/models"
)

type MockEventService struct {
	Events map[string]*models.Event
}

func NewMockEventService() *MockEventService {
	return &MockEventService{
		Events: make(map[string]*models.Event),
	}
}

func (mes *MockEventService) CreateEvent(ctx context.Context, event *models.Event) error {
	if _, exists := mes.Events[event.EventID]; exists {
		return fmt.Errorf("event already exists")
	}
	mes.Events[event.EventID] = event
	return nil
}

func (mes *MockEventService) GetEvent(ctx context.Context, userEmail, eventID string) (*models.Event, error) {
	event, exists := mes.Events[eventID]
	if !exists || event.Email != userEmail {
		return nil, fmt.Errorf("event not found")
	}
	return event, nil
}

func (mes *MockEventService) UpdateEvent(ctx context.Context, event *models.Event) error {
	existingEvent, exists := mes.Events[event.EventID]
	if !exists || existingEvent.Email != event.Email {
		return fmt.Errorf("event not found")
	}
	mes.Events[event.EventID] = event
	return nil
}

func (mes *MockEventService) DeleteEvent(ctx context.Context, userEmail, eventID string) error {
	event, exists := mes.Events[eventID]
	if !exists || event.Email != userEmail {
		return fmt.Errorf("event not found")
	}
	delete(mes.Events, eventID)
	return nil
}

func (mes *MockEventService) GetAllEvents(ctx context.Context, userEmail string) ([]models.Event, error) {
	var events []models.Event
	for _, event := range mes.Events {
		if event.Email == userEmail {
			events = append(events, *event)
		}
	}
	return events, nil
}
