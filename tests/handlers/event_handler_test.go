/**
 *  EventHandler Tests validate the behavior of the EventHandler methods.
 *  They use a mock EventService to isolate the handler logic and verify interactions with the service.
 *
 *  @file       event_handler_test.go
 *  @package    handlers_test
 *
 *  @test_cases
 *  - TestEventHandler_CreateEvent      - Tests the creation of an event.
 *  - TestEventHandler_GetEvent         - Tests retrieving a specific event by ID.
 *  - TestEventHandler_UpdateEvent      - Tests updating an existing event.
 *  - TestEventHandler_DeleteEvent      - Tests deleting an event.
 *  - TestEventHandler_GetAllEvents     - Tests retrieving all events for a user.
 *
 *  @dependencies
 *  - mocks.NewMockEventService: Mock implementation of EventService for testing.
 *  - httptest: Provides utilities for testing HTTP handlers.
 *  - context.WithValue: Adds user-specific context values for testing purposes.
 *  - encoding/json: Handles JSON marshalling and unmarshalling.
 *
 *  @behaviors
 *  - Verifies HTTP status codes for each handler.
 *  - Validates request/response data consistency.
 *  - Confirms the correct service methods are called during handler execution.
 *
 *  @authors
 *      - Aayush
 *      - Tung
 *      - Boss
 *      - Majd
 */

package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"proh2052-group6/internal/handlers"
	"proh2052-group6/pkg/models"
	"proh2052-group6/tests/mocks"
)

func TestEventHandler_CreateEvent(t *testing.T) {
	// Setup mock EventService and EventHandler
	mockEventService := mocks.NewMockEventService()
	eventHandler := handlers.NewEventHandler(mockEventService)

	// Prepare request body
	event := models.Event{
		Title:       "Meeting",
		Description: "Team meeting",
		Date:        "2023-10-15",
	}
	requestBody, _ := json.Marshal(event)

	// Create HTTP request
	req, err := http.NewRequest("POST", "/api/events/create", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Inject userEmail into context
	userEmail := "test@example.com"
	ctx := context.WithValue(req.Context(), "userEmail", userEmail)
	req = req.WithContext(ctx)

	// Create ResponseRecorder to capture response
	rr := httptest.NewRecorder()

	// Invoke handler
	handler := http.HandlerFunc(eventHandler.CreateEvent)
	handler.ServeHTTP(rr, req)

	// Assert status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Parse and validate response
	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}
	expectedMessage := "Event created successfully"
	if response["message"] != expectedMessage {
		t.Errorf("Expected message '%s', got '%s'", expectedMessage, response["message"])
	}
	if response["eventID"] == "" {
		t.Errorf("Expected eventID in response")
	}

	// Verify event saved in mock service
	savedEvent, err := mockEventService.GetEvent(context.Background(), userEmail, response["eventID"])
	if err != nil {
		t.Errorf("Event was not saved in the service: %v", err)
	}
	if savedEvent.Title != event.Title {
		t.Errorf("Expected event title '%s', got '%s'", event.Title, savedEvent.Title)
	}
}

func TestEventHandler_GetEvent(t *testing.T) {
	// Setup mock EventService and EventHandler
	mockEventService := mocks.NewMockEventService()
	eventHandler := handlers.NewEventHandler(mockEventService)

	// Add event to mock service
	userEmail := "test@example.com"
	eventID := "event123"
	event := &models.Event{
		EventID:     eventID,
		Email:       userEmail,
		Title:       "Meeting",
		Description: "Team meeting",
		Date:        "2023-10-15",
	}
	mockEventService.Events[eventID] = event

	// Create HTTP request
	req, err := http.NewRequest("GET", "/api/events/get?eventID="+eventID, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Inject userEmail into context
	ctx := context.WithValue(req.Context(), "userEmail", userEmail)
	req = req.WithContext(ctx)

	// Create ResponseRecorder
	rr := httptest.NewRecorder()

	// Invoke handler
	handler := http.HandlerFunc(eventHandler.GetEvent)
	handler.ServeHTTP(rr, req)

	// Assert status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Parse and validate response
	var response models.Event
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}
	if response.EventID != eventID {
		t.Errorf("Expected eventID '%s', got '%s'", eventID, response.EventID)
	}
	if response.Title != event.Title {
		t.Errorf("Expected event title '%s', got '%s'", event.Title, response.Title)
	}
}

func TestEventHandler_UpdateEvent(t *testing.T) {
	// Create a mock event service
	mockEventService := mocks.NewMockEventService()
	eventHandler := handlers.NewEventHandler(mockEventService)

	// Add an event to the mock service
	userEmail := "test@example.com"
	eventID := "event123"
	event := &models.Event{
		EventID:     eventID,
		Email:       userEmail,
		Title:       "Meeting",
		Description: "Team meeting",
		Date:        "2023-10-15",
	}
	mockEventService.Events[eventID] = event

	// Prepare the updated event data
	updatedEvent := models.Event{
		Title:       "Updated Meeting",
		Description: "Updated team meeting",
		Date:        "2023-10-16",
	}
	requestBody, _ := json.Marshal(updatedEvent)

	// Create a new HTTP request
	req, err := http.NewRequest("PUT", "/api/events/update?eventID="+eventID, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Set the userEmail in the context
	ctx := context.WithValue(req.Context(), "userEmail", userEmail)
	req = req.WithContext(ctx)

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler
	handler := http.HandlerFunc(eventHandler.UpdateEvent)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Verify that the event was updated in the mock service
	updatedEventInService, err := mockEventService.GetEvent(context.Background(), userEmail, eventID)
	if err != nil {
		t.Errorf("Event was not found in the service: %v", err)
	}

	if updatedEventInService.Title != updatedEvent.Title {
		t.Errorf("Expected event title '%s', got '%s'", updatedEvent.Title, updatedEventInService.Title)
	}
}

func TestEventHandler_DeleteEvent(t *testing.T) {
	// Create a mock event service
	mockEventService := mocks.NewMockEventService()
	eventHandler := handlers.NewEventHandler(mockEventService)

	// Add an event to the mock service
	userEmail := "test@example.com"
	eventID := "event123"
	event := &models.Event{
		EventID:     eventID,
		Email:       userEmail,
		Title:       "Meeting",
		Description: "Team meeting",
		Date:        "2023-10-15",
	}
	mockEventService.Events[eventID] = event

	// Create a new HTTP request
	req, err := http.NewRequest("DELETE", "/api/events/delete?eventID="+eventID, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Set the userEmail in the context
	ctx := context.WithValue(req.Context(), "userEmail", userEmail)
	req = req.WithContext(ctx)

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler
	handler := http.HandlerFunc(eventHandler.DeleteEvent)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Verify that the event was deleted from the mock service
	_, err = mockEventService.GetEvent(context.Background(), userEmail, eventID)
	if err == nil {
		t.Errorf("Expected event to be deleted, but it still exists")
	}
}

func TestEventHandler_GetAllEvents(t *testing.T) {
	// Create a mock event service
	mockEventService := mocks.NewMockEventService()
	eventHandler := handlers.NewEventHandler(mockEventService)

	// Add events to the mock service
	userEmail := "test@example.com"
	event1 := &models.Event{
		EventID:     "event1",
		Email:       userEmail,
		Title:       "Meeting",
		Description: "Team meeting",
		Date:        "2023-10-15",
	}
	event2 := &models.Event{
		EventID:     "event2",
		Email:       userEmail,
		Title:       "Conference",
		Description: "Annual conference",
		Date:        "2023-11-20",
	}
	mockEventService.Events[event1.EventID] = event1
	mockEventService.Events[event2.EventID] = event2

	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/api/events/all", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Set the userEmail in the context
	ctx := context.WithValue(req.Context(), "userEmail", userEmail)
	req = req.WithContext(ctx)

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler
	handler := http.HandlerFunc(eventHandler.GetAllEvents)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body
	var response []models.Event
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	if len(response) != 2 {
		t.Errorf("Expected 2 events, got %d", len(response))
	}
}
