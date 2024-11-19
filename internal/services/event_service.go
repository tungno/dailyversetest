// internal/services/event_service.go
package services

import (
	"context"
	"fmt"
	"strings"
	"time"

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
	DB DatabaseInterface
}

func NewEventService(db DatabaseInterface) EventServiceInterface {
	return &EventService{DB: db}
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

	// Add the event to Firestore under the user's events subcollection
	userDocRef := es.DB.Collection("users").Doc(event.Email).Collection("events")
	docRef, _, err := userDocRef.Add(ctx, event)
	if err != nil {
		return fmt.Errorf("Failed to create event")
	}

	event.EventID = docRef.ID
	_, err = docRef.Set(ctx, event)
	if err != nil {
		return fmt.Errorf("Failed to update event with EventID")
	}

	return nil
}

func (es *EventService) GetEvent(ctx context.Context, userEmail, eventID string) (*models.Event, error) {
	docRef := es.DB.Collection("users").Doc(userEmail).Collection("events").Doc(eventID)
	doc, err := docRef.Get(ctx)
	if err != nil || !doc.Exists() {
		return nil, fmt.Errorf("Event not found")
	}

	var event models.Event
	err = doc.DataTo(&event)
	if err != nil {
		return nil, fmt.Errorf("Error parsing event data")
	}

	if event.Email != userEmail {
		return nil, fmt.Errorf("Unauthorized to access this event")
	}

	return &event, nil
}

func (es *EventService) UpdateEvent(ctx context.Context, event *models.Event) error {
	docRef := es.DB.Collection("users").Doc(event.Email).Collection("events").Doc(event.EventID)
	_, err := docRef.Set(ctx, event)
	if err != nil {
		return fmt.Errorf("Failed to update event")
	}
	return nil
}

func (es *EventService) DeleteEvent(ctx context.Context, userEmail, eventID string) error {
	docRef := es.DB.Collection("users").Doc(userEmail).Collection("events").Doc(eventID)
	_, err := docRef.Delete(ctx)
	if err != nil {
		return fmt.Errorf("Failed to delete event")
	}
	return nil
}

func (es *EventService) GetAllEvents(ctx context.Context, userEmail string) ([]models.Event, error) {
	var events []models.Event

	// Query the user's own events
	userEventsDocs, err := es.DB.Collection("users").Doc(userEmail).Collection("events").Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch user's events")
	}

	for _, doc := range userEventsDocs {
		var event models.Event
		err := doc.DataTo(&event)
		if err != nil {
			return nil, fmt.Errorf("Error parsing event data")
		}

		event.EventID = doc.Ref.ID
		events = append(events, event)
	}

	// You can implement fetching public events from friends here...

	return events, nil
}
