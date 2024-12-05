/**
 *  TimetableHandler is responsible for handling HTTP requests related to timetable operations,
 *  including importing timetables from ICS content. This handler integrates with the
 *  TimetableService to provide the necessary functionality.
 *
 *  @struct   TimetableHandler
 *  @inherits None
 *
 *  @properties
 *  - TimetableService - A service interface for managing timetable-related operations.
 *
 *  @methods
 *  - NewTimetableHandler(ts)               - Initializes a new TimetableHandler with the required service.
 *  - ImportTimetable(w, r)                 - Handles POST requests to import timetables from ICS content.
 *
 *  @endpoints
 *  - /api/timetables/import (POST)
 *    - HTTP Method: POST
 *    - Request Body: JSON object containing ICS content.
 *    - Behavior: Imports a timetable for the authenticated user based on the provided ICS content.
 *
 *  @behaviors
 *  - Validates the presence of required parameters (e.g., `icsContent`) and request body fields.
 *  - Returns a 400 Bad Request error if parameters or body content are invalid or missing.
 *  - Returns a 401 Unauthorized error if the user is not authenticated.
 *  - Returns a 500 Internal Server Error if an error occurs during processing.
 *  - On success, returns a JSON object containing a success message.
 *
 *  @examples
 *  Import Timetable:
 *  ```
 *  POST /api/timetables/import
 *  Body: {
 *      "icsContent": "BEGIN:VCALENDAR\nVERSION:2.0\n..."
 *  }
 *
 *  Response:
 *  {
 *      "message": "Timetable imported successfully"
 *  }
 *  ```
 *
 *  @dependencies
 *  - TimetableServiceInterface: Provides methods for timetable import and management.
 *  - utils.WriteJSON: Utility function to write JSON responses.
 *  - utils.WriteJSONError: Utility function to write error responses in JSON format.
 *
 *  @file      timetable_handler.go
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
	"proh2052-group6/pkg/utils"
)

// TimetableHandler struct handles requests related to timetable operations.
type TimetableHandler struct {
	TimetableService services.TimetableServiceInterface // Service for managing timetable-related logic.
}

// NewTimetableHandler initializes a new TimetableHandler with the necessary dependencies.
func NewTimetableHandler(ts services.TimetableServiceInterface) *TimetableHandler {
	return &TimetableHandler{TimetableService: ts}
}

// ImportTimetable handles POST requests to import a timetable using ICS content.
// Endpoint: /api/timetables/import
func (th *TimetableHandler) ImportTimetable(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		ICSContent string `json:"icsContent"` // The ICS content of the timetable to import.
	}

	// Decode the request body into the requestData struct.
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		utils.WriteJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate that ICSContent is not empty.
	if requestData.ICSContent == "" {
		utils.WriteJSONError(w, "ICS content is required", http.StatusBadRequest)
		return
	}

	// Retrieve the authenticated user's email from the request context.
	userEmail, ok := r.Context().Value("userEmail").(string)
	if !ok {
		utils.WriteJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Attempt to import the timetable using the service.
	err := th.TimetableService.ImportTimetable(r.Context(), userEmail, requestData.ICSContent)
	if err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with a success message.
	utils.WriteJSON(w, map[string]string{
		"message": "Timetable imported successfully",
	})
}
