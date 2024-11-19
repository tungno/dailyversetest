// internal/repositories/firestore_event_repository.go
package repositories

import (
	"context"
	"fmt"
	"proh2052-group6/pkg/models"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type FirestoreEventRepository struct {
	Client *firestore.Client
}

func NewFirestoreEventRepository(client *firestore.Client) EventRepository {
	return &FirestoreEventRepository{Client: client}
}

func (er *FirestoreEventRepository) CreateEvent(ctx context.Context, event *models.Event) error {
	userEventsCollection := er.Client.Collection("users").Doc(event.Email).Collection("events")
	docRef, _, err := userEventsCollection.Add(ctx, event)
	if err != nil {
		return fmt.Errorf("Failed to create event: %v", err)
	}

	event.EventID = docRef.ID
	_, err = docRef.Set(ctx, event)
	if err != nil {
		return fmt.Errorf("Failed to update event with EventID: %v", err)
	}

	return nil
}

func (er *FirestoreEventRepository) GetEvent(ctx context.Context, userEmail, eventID string) (*models.Event, error) {
	docRef := er.Client.Collection("users").Doc(userEmail).Collection("events").Doc(eventID)
	doc, err := docRef.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("Event not found: %v", err)
	}

	var event models.Event
	err = doc.DataTo(&event)
	if err != nil {
		return nil, fmt.Errorf("Error parsing event data: %v", err)
	}

	return &event, nil
}

func (er *FirestoreEventRepository) UpdateEvent(ctx context.Context, event *models.Event) error {
	docRef := er.Client.Collection("users").Doc(event.Email).Collection("events").Doc(event.EventID)
	_, err := docRef.Set(ctx, event)
	if err != nil {
		return fmt.Errorf("Failed to update event: %v", err)
	}
	return nil
}

func (er *FirestoreEventRepository) DeleteEvent(ctx context.Context, userEmail, eventID string) error {
	docRef := er.Client.Collection("users").Doc(userEmail).Collection("events").Doc(eventID)
	_, err := docRef.Delete(ctx)
	if err != nil {
		return fmt.Errorf("Failed to delete event: %v", err)
	}
	return nil
}

func (er *FirestoreEventRepository) GetAllEvents(ctx context.Context, userEmail string) ([]models.Event, error) {
	var events []models.Event

	iter := er.Client.Collection("users").Doc(userEmail).Collection("events").Documents(ctx)
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Failed to fetch user's events: %v", err)
		}

		var event models.Event
		err = doc.DataTo(&event)
		if err != nil {
			return nil, fmt.Errorf("Error parsing event data: %v", err)
		}

		event.EventID = doc.Ref.ID
		events = append(events, event)
	}

	return events, nil
}
