/**
 *  ProfileHandler handles HTTP requests related to user profile management.
 *  This handler allows users to view and update their profile information.
 *
 *  @struct   ProfileHandler
 *  @inherits None
 *
 *  @methods
 *  - NewProfileHandler(ps)           - Initializes a new ProfileHandler instance with a ProfileService interface.
 *  - ProfileHandler(w, r)            - Routes HTTP requests based on the HTTP method.
 *  - GetProfile(w, r)                - Handles GET requests to fetch the authenticated user's profile.
 *  - UpdateProfile(w, r)             - Handles PUT requests to update the authenticated user's profile.
 *
 *  @endpoints
 *  - /api/profile
 *    - HTTP Method: GET
 *      - Fetches the profile information of the authenticated user.
 *    - HTTP Method: PUT
 *      - Body: `{ "field1": "value1", "field2": "value2", ... }`
 *      - Updates the profile information of the authenticated user with the provided data.
 *
 *  @behaviors
 *  - Ensures user authentication by retrieving `userEmail` from the request context.
 *  - Returns meaningful status codes based on the success or failure of operations.
 *  - Validates request payloads for PUT requests.
 *
 *  @example
 *  ```
 *  GET /api/profile
 *
 *  Response:
 *  {
 *      "name": "John Doe",
 *      "email": "john.doe@example.com",
 *      "age": 30,
 *      "location": "Norway"
 *  }
 *
 *  PUT /api/profile
 *  Body: { "location": "Canada" }
 *
 *  Response:
 *  { "message": "Successfully updated profile" }
 *  ```
 *
 *  @dependencies
 *  - services.ProfileServiceInterface: Interface for managing profile-related operations.
 *  - utils: Utility package for writing JSON responses and errors.
 *
 *  @file      profile_handler.go
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

// ProfileHandler struct for handling profile-related requests.
type ProfileHandler struct {
	ProfileService services.ProfileServiceInterface
}

// NewProfileHandler initializes a new ProfileHandler instance.
func NewProfileHandler(ps services.ProfileServiceInterface) *ProfileHandler {
	return &ProfileHandler{ProfileService: ps}
}

// ProfileHandler routes HTTP requests based on the HTTP method.
// Supported Methods:
//   - GET: Fetches the user's profile.
//   - PUT: Updates the user's profile.
func (ph *ProfileHandler) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		ph.GetProfile(w, r)
	case "PUT":
		ph.UpdateProfile(w, r)
	default:
		utils.WriteJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetProfile handles GET requests to fetch the authenticated user's profile.
func (ph *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userEmail := r.Context().Value("userEmail").(string)

	profileData, err := ph.ProfileService.GetProfile(r.Context(), userEmail)
	if err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, profileData)
}

// UpdateProfile handles PUT requests to update the authenticated user's profile.
func (ph *ProfileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userEmail := r.Context().Value("userEmail").(string)

	var updatedData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		utils.WriteJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := ph.ProfileService.UpdateProfile(r.Context(), userEmail, updatedData); err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, map[string]string{"message": "Successfully updated profile"})
}
