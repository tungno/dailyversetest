/**
 *  JournalHandler is responsible for handling HTTP requests related to journal operations,
 *  including creating, retrieving, updating, and deleting journals. This handler integrates
 *  with the JournalService to provide the necessary functionality.
 *
 *  @struct   JournalHandler
 *  @inherits None
 *
 *  @properties
 *  - JournalService - A service interface for managing journal-related operations.
 *
 *  @methods
 *  - NewJournalHandler(js)                - Initializes a new JournalHandler with the required service.
 *  - CreateJournal(w, r)                  - Handles POST requests to create a new journal.
 *  - GetJournal(w, r)                     - Handles GET requests to fetch a specific journal by its ID.
 *  - UpdateJournal(w, r)                  - Handles PUT requests to update an existing journal by its ID.
 *  - DeleteJournal(w, r)                  - Handles DELETE requests to delete a specific journal by its ID.
 *  - GetAllJournals(w, r)                 - Handles GET requests to fetch all journals for the logged-in user.
 *
 *  @endpoints
 *  - /api/journals (POST)
 *    - HTTP Method: POST
 *    - Request Body: JSON object representing a journal.
 *    - Behavior: Creates a new journal for the authenticated user.
 *
 *  - /api/journals/{journalID} (GET)
 *    - HTTP Method: GET
 *    - Query Parameter: `journalID` (required) - The ID of the journal to retrieve.
 *    - Behavior: Fetches a specific journal by ID for the authenticated user.
 *
 *  - /api/journals/{journalID} (PUT)
 *    - HTTP Method: PUT
 *    - Query Parameter: `journalID` (required) - The ID of the journal to update.
 *    - Request Body: JSON object representing updated journal data.
 *    - Behavior: Updates the specified journal for the authenticated user.
 *
 *  - /api/journals/{journalID} (DELETE)
 *    - HTTP Method: DELETE
 *    - Query Parameter: `journalID` (required) - The ID of the journal to delete.
 *    - Behavior: Deletes the specified journal for the authenticated user.
 *
 *  - /api/journals (GET)
 *    - HTTP Method: GET
 *    - Behavior: Fetches all journals for the authenticated user.
 *
 *  @behaviors
 *  - Validates the presence of required parameters (e.g., `journalID`) and request body fields.
 *  - Returns a 400 Bad Request error if parameters or body content are invalid or missing.
 *  - Returns a 404 Not Found error if the specified journal does not exist.
 *  - Returns a 500 Internal Server Error if an error occurs during processing.
 *  - On success, returns a JSON object containing the journal data or a success message.
 *
 *  @examples
 *  Create Journal:
 *  ```
 *  POST /api/journals
 *  Body: {
 *      "title": "My Journal",
 *      "content": "Today was a good day."
 *  }
 *
 *  Response:
 *  {
 *      "message": "Journal created successfully",
 *      "journalID": "12345"
 *  }
 *  ```
 *
 *  Get Journal:
 *  ```
 *  GET /api/journals?journalID=12345
 *
 *  Response:
 *  {
 *      "journalID": "12345",
 *      "title": "My Journal",
 *      "content": "Today was a good day.",
 *      "email": "user@example.com"
 *  }
 *  ```
 *
 *  @dependencies
 *  - JournalServiceInterface: Provides methods for journal management (CRUD operations).
 *  - utils.WriteJSON: Utility function to write JSON responses.
 *  - utils.WriteJSONError: Utility function to write error responses in JSON format.
 *
 *  @file      journal_handler.go
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

// JournalHandler struct handles requests related to journal operations.
type JournalHandler struct {
	JournalService services.JournalServiceInterface // Service for managing journal-related logic.
}

// NewJournalHandler initializes a new JournalHandler with the necessary dependencies.
func NewJournalHandler(js services.JournalServiceInterface) *JournalHandler {
	return &JournalHandler{JournalService: js}
}

// CreateJournal handles POST requests to create a new journal.
// Endpoint: /api/journals
func (jh *JournalHandler) CreateJournal(w http.ResponseWriter, r *http.Request) {
	var journal models.Journal
	if err := json.NewDecoder(r.Body).Decode(&journal); err != nil {
		utils.WriteJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userEmail := r.Context().Value("userEmail").(string)
	journal.Email = userEmail

	if err := jh.JournalService.CreateJournal(r.Context(), &journal); err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, map[string]string{
		"message":   "Journal created successfully",
		"journalID": journal.JournalID,
	})
}

// GetJournal handles GET requests to retrieve a specific journal by ID.
// Endpoint: /api/journals/{journalID}
func (jh *JournalHandler) GetJournal(w http.ResponseWriter, r *http.Request) {
	journalID := r.URL.Query().Get("journalID")
	if journalID == "" {
		utils.WriteJSONError(w, "Missing journalID parameter", http.StatusBadRequest)
		return
	}

	userEmail := r.Context().Value("userEmail").(string)
	journal, err := jh.JournalService.GetJournal(r.Context(), userEmail, journalID)
	if err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusNotFound)
		return
	}

	utils.WriteJSON(w, journal)
}

// UpdateJournal handles PUT requests to update an existing journal by ID.
// Endpoint: /api/journals/{journalID}
func (jh *JournalHandler) UpdateJournal(w http.ResponseWriter, r *http.Request) {
	journalID := r.URL.Query().Get("journalID")
	if journalID == "" {
		utils.WriteJSONError(w, "Missing journalID parameter", http.StatusBadRequest)
		return
	}

	var journal models.Journal
	if err := json.NewDecoder(r.Body).Decode(&journal); err != nil {
		utils.WriteJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userEmail := r.Context().Value("userEmail").(string)
	journal.Email = userEmail
	journal.JournalID = journalID

	if err := jh.JournalService.UpdateJournal(r.Context(), &journal); err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, map[string]string{"message": "Journal updated successfully"})
}

// DeleteJournal handles DELETE requests to delete a specific journal by ID.
// Endpoint: /api/journals/{journalID}
func (jh *JournalHandler) DeleteJournal(w http.ResponseWriter, r *http.Request) {
	journalID := r.URL.Query().Get("journalID")
	if journalID == "" {
		utils.WriteJSONError(w, "Missing journalID parameter", http.StatusBadRequest)
		return
	}

	userEmail := r.Context().Value("userEmail").(string)

	if err := jh.JournalService.DeleteJournal(r.Context(), userEmail, journalID); err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, map[string]string{"message": "Journal deleted successfully"})
}

// GetAllJournals handles GET requests to fetch all journals for the logged-in user.
// Endpoint: /api/journals
func (jh *JournalHandler) GetAllJournals(w http.ResponseWriter, r *http.Request) {
	userEmail := r.Context().Value("userEmail").(string)

	journals, err := jh.JournalService.GetAllJournals(r.Context(), userEmail)
	if err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, journals)
}
