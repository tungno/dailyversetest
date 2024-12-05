/**
 *  FirestoreEventRepository provides methods to interact with the Firestore database for event-related operations.
 *  This repository encapsulates CRUD operations for managing events tied to specific user accounts.
 *
 *  @struct   FirestoreEventRepository
 *  @inherits None
 *
 *  @methods
 *  - NewFirestoreEventRepository(client) - Initializes a new FirestoreEventRepository with a Firestore client.
 *  - CreateEvent(ctx, event)             - Creates a new event for a user in Firestore.
 *  - GetEvent(ctx, userEmail, eventID)   - Fetches a specific event for a user by its ID.
 *  - UpdateEvent(ctx, event)             - Updates an existing event in Firestore.
 *  - DeleteEvent(ctx, userEmail, eventID)- Deletes a specific event for a user by its ID.
 *  - GetAllEvents(ctx, userEmail)        - Retrieves all events for a user from Firestore.
 *
 *  @behaviors
 *  - Uses Firestore's hierarchical document structure to store user-specific events under `users/{userEmail}/events/{eventID}`.
 *  - Handles error scenarios and returns meaningful messages on failure.
 *  - Ensures seamless conversion between Firestore documents and the `models.Event` struct.
 *
 *  @dependencies
 *  - cloud.google.com/go/firestore: Firestore client for database operations.
 *  - google.golang.org/api/iterator: Iterator for traversing Firestore query results.
 *  - models.Event: Struct representing event data.
 *
 *  @example
 *  ```
 *  // Create a new event
 *  event := &models.Event{
 *      Email: "user@example.com",
 *      Title: "Meeting",
 *      Date: "2024-12-01",
 *  }
 *  err := repository.CreateEvent(ctx, event)
 *
 *  // Fetch all events for a user
 *  events, err := repository.GetAllEvents(ctx, "user@example.com")
 *  ```
 *
 *  @file      firestore_event_repository.go
 *  @project   DailyVerse
 *  @framework Firestore Client (Go) API
 *  @authors
 *      - Aayush
 *      - Tung
 *      - Boss
 *      - Majd
 */

package repositories

import (
	"context"
	"fmt"
	"proh2052-group6/pkg/models"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// FirestoreEventRepository implements the EventRepository interface for Firestore.
type FirestoreEventRepository struct {
	Client *firestore.Client
}

// NewFirestoreEventRepository initializes a new FirestoreEventRepository with the given Firestore client.
func NewFirestoreEventRepository(client *firestore.Client) EventRepository {
	return &FirestoreEventRepository{Client: client}
}

// CreateEvent creates a new event for a user in Firestore.
func (er *FirestoreEventRepository) CreateEvent(ctx context.Context, event *models.Event) error {
	userEventsCollection := er.Client.Collection("users").Doc(event.Email).Collection("events")
	docRef, _, err := userEventsCollection.Add(ctx, event)
	if err != nil {
		return fmt.Errorf("Failed to create event: %v", err)
	}

	// Assign the generated EventID back to the event object and update the Firestore document.
	event.EventID = docRef.ID
	_, err = docRef.Set(ctx, event)
	if err != nil {
		return fmt.Errorf("Failed to update event with EventID: %v", err)
	}

	return nil
}

// GetEvent retrieves a specific event for a user by its ID.
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

// UpdateEvent updates an existing event in Firestore.
func (er *FirestoreEventRepository) UpdateEvent(ctx context.Context, event *models.Event) error {
	docRef := er.Client.Collection("users").Doc(event.Email).Collection("events").Doc(event.EventID)
	_, err := docRef.Set(ctx, event)
	if err != nil {
		return fmt.Errorf("Failed to update event: %v", err)
	}
	return nil
}

// DeleteEvent deletes a specific event for a user by its ID.
func (er *FirestoreEventRepository) DeleteEvent(ctx context.Context, userEmail, eventID string) error {
	docRef := er.Client.Collection("users").Doc(userEmail).Collection("events").Doc(eventID)
	_, err := docRef.Delete(ctx)
	if err != nil {
		return fmt.Errorf("Failed to delete event: %v", err)
	}
	return nil
}

// GetAllEvents retrieves all events for a user from Firestore.
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

		// Assign the Firestore document ID to the EventID field.
		event.EventID = doc.Ref.ID
		events = append(events, event)
	}

	return events, nil
}
