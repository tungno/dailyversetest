// internal/handlers/event_handler.go
package handlers

import (
	"encoding/json"
	"net/http"

	"proh2052-group6/internal/services"
	"proh2052-group6/pkg/models"
	"proh2052-group6/pkg/utils"
)

type EventHandler struct {
	EventService services.EventServiceInterface
}

func NewEventHandler(es services.EventServiceInterface) *EventHandler {
	return &EventHandler{EventService: es}
}

func (eh *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var event models.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		utils.WriteJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

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

	userEmail := r.Context().Value("userEmail").(string)
	event.Email = userEmail
	event.EventID = eventID

	if err := eh.EventService.UpdateEvent(r.Context(), &event); err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, map[string]string{"message": "Event updated successfully"})
}

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

func (eh *EventHandler) GetAllEvents(w http.ResponseWriter, r *http.Request) {
	userEmail := r.Context().Value("userEmail").(string)

	events, err := eh.EventService.GetAllEvents(r.Context(), userEmail)
	if err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, events)
}
