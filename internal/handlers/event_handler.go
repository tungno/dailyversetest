/**
 *  EventHandler handles HTTP requests related to events, including creating, retrieving,
 *  updating, and deleting events. It integrates with the EventService to perform operations
 *  and returns appropriate HTTP responses.
 *
 *  @struct   EventHandler
 *  @inherits None
 *
 *  @methods
 *  - NewEventHandler(es)         - Initializes a new EventHandler with the required EventService.
 *  - CreateEvent(w, r)           - Handles event creation requests.
 *  - GetEvent(w, r)              - Fetches a single event by its ID.
 *  - UpdateEvent(w, r)           - Updates an existing event.
 *  - DeleteEvent(w, r)           - Deletes an event by its ID.
 *  - GetAllEvents(w, r)          - Retrieves all events for the authenticated user.
 *
 *  @endpoint
 *  - /api/events/create
 *    - Method: POST
 *    - Body: Event object
 *  - /api/events/get
 *    - Method: GET
 *    - Query Parameter: eventID (string, required)
 *  - /api/events/update
 *    - Method: PUT
 *    - Query Parameter: eventID (string, required)
 *    - Body: Updated Event object
 *  - /api/events/delete
 *    - Method: DELETE
 *    - Query Parameter: eventID (string, required)
 *  - /api/events/all
 *    - Method: GET
 *
 *  @behaviors
 *  - Returns 400 Bad Request for missing or invalid inputs.
 *  - Returns 404 Not Found for non-existent event IDs.
 *  - Returns 500 Internal Server Error for service-layer failures.
 *  - On success, responds with appropriate HTTP status codes and data.
 *
 *  @dependencies
 *  - EventServiceInterface: Provides business logic for managing events.
 *  - utils.WriteJSON, utils.WriteJSONError: Utility functions for JSON responses.
 *
 *  @file      event_handler.go
 *  @project   DailyVerse
 *  @framework Go HTTP Server
 *  @authors
 *      - Aayush
 *      - Tung
 *      - Boss
 *      - Majd
 */

package handlers

import (
	"encoding/json"
	"net/http"

	"proh2052-group6/internal/services"
	"proh2052-group6/pkg/models"
	"proh2052-group6/pkg/utils"
)

// EventHandler manages HTTP requests related to event operations.
type EventHandler struct {
	EventService services.EventServiceInterface // Service for event-related operations.
}

// NewEventHandler initializes an EventHandler with the given EventService.
func NewEventHandler(es services.EventServiceInterface) *EventHandler {
	return &EventHandler{EventService: es}
}

// CreateEvent handles POST requests to create a new event.
// Body: JSON-encoded Event object.
func (eh *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var event models.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		utils.WriteJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Attach user email from context to the event.
	userEmail := r.Context().Value("userEmail").(string)
	event.Email = userEmail

	if err := eh.EventService.CreateEvent(r.Context(), &event); err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, map[string]string{
		"message": "Event created successfully",
		"eventID": event.EventID,
	})
}

// GetEvent handles GET requests to fetch a specific event by its ID.
// Query Parameter: eventID (string, required).
func (eh *EventHandler) GetEvent(w http.ResponseWriter, r *http.Request) {
	eventID := r.URL.Query().Get("eventID")
	if eventID == "" {
		utils.WriteJSONError(w, "Missing eventID parameter", http.StatusBadRequest)
		return
	}

	userEmail := r.Context().Value("userEmail").(string)
	event, err := eh.EventService.GetEvent(r.Context(), userEmail, eventID)
	if err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusNotFound)
		return
	}

	utils.WriteJSON(w, event)
}

// UpdateEvent handles PUT requests to update an existing event.
// Query Parameter: eventID (string, required).
// Body: JSON-encoded Event object with updated details.
func (eh *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	eventID := r.URL.Query().Get("eventID")
	if eventID == "" {
		utils.WriteJSONError(w, "Missing eventID parameter", http.StatusBadRequest)
		return
	}

	var event models.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		utils.WriteJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Attach user email and event ID to the event.
	userEmail := r.Context().Value("userEmail").(string)
	event.Email = userEmail
	event.EventID = eventID

	if err := eh.EventService.UpdateEvent(r.Context(), &event); err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, map[string]string{"message": "Event updated successfully"})
}

// DeleteEvent handles DELETE requests to remove an event by its ID.
// Query Parameter: eventID (string, required).
func (eh *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	eventID := r.URL.Query().Get("eventID")
	if eventID == "" {
		utils.WriteJSONError(w, "Missing eventID parameter", http.StatusBadRequest)
		return
	}

	userEmail := r.Context().Value("userEmail").(string)

	if err := eh.EventService.DeleteEvent(r.Context(), userEmail, eventID); err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, map[string]string{"message": "Event deleted successfully"})
}

// GetAllEvents handles GET requests to fetch all events for the authenticated user.
func (eh *EventHandler) GetAllEvents(w http.ResponseWriter, r *http.Request) {
	userEmail := r.Context().Value("userEmail").(string)

	events, err := eh.EventService.GetAllEvents(r.Context(), userEmail)
	if err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, events)
}
