/**
 *  FriendHandler handles HTTP requests related to friend interactions in the application.
 *  This handler supports actions such as sending, accepting, and declining friend requests,
 *  as well as retrieving friend lists and managing friendships.
 *
 *  @struct   FriendHandler
 *  @inherits None
 *
 *  @methods
 *  - NewFriendHandler(fs)               - Initializes a new FriendHandler instance with a FriendService interface.
 *  - SendFriendRequest(w, r)           - Handles POST requests to send a friend request to a user.
 *  - AcceptFriendRequest(w, r)         - Handles POST requests to accept a friend request.
 *  - GetFriendsList(w, r)              - Handles GET requests to fetch a user's list of friends.
 *  - RemoveFriend(w, r)                - Handles DELETE requests to remove a friend from a user's friend list.
 *  - GetPendingFriendRequests(w, r)    - Handles GET requests to fetch pending friend requests for a user.
 *  - DeclineFriendRequest(w, r)        - Handles POST requests to decline a friend request.
 *  - CancelFriendRequest(w, r)         - Handles DELETE requests to cancel a sent friend request.
 *
 *  @endpoints
 *  - /api/friends/send
 *    - HTTP Method: POST
 *    - Body: `{ "usernameOrEmail": "string" }`
 *    - Sends a friend request to the specified user by username or email.
 *
 *  - /api/friends/accept
 *    - HTTP Method: POST
 *    - Body: `{ "usernameOrEmail": "string" }`
 *    - Accepts a friend request from the specified user by username or email.
 *
 *  - /api/friends/list
 *    - HTTP Method: GET
 *    - Fetches the list of friends for the authenticated user.
 *
 *  - /api/friends/remove
 *    - HTTP Method: DELETE
 *    - Body: `{ "username": "string" }`
 *    - Removes the specified user from the authenticated user's friend list.
 *
 *  - /api/friends/pending
 *    - HTTP Method: GET
 *    - Fetches the pending friend requests for the authenticated user.
 *
 *  - /api/friends/decline
 *    - HTTP Method: POST
 *    - Body: `{ "usernameOrEmail": "string" }`
 *    - Declines a friend request from the specified user by username or email.
 *
 *  - /api/friends/cancel
 *    - HTTP Method: DELETE
 *    - Body: `{ "username": "string" }`
 *    - Cancels a sent friend request to the specified user.
 *
 *  @behaviors
 *  - Validates request payloads and responds with appropriate error messages for invalid inputs.
 *  - Ensures user authentication via `userEmail` in the request context.
 *  - Returns meaningful status codes based on the success or failure of operations.
 *
 *  @example
 *  ```
 *  POST /api/friends/send
 *  Body: { "usernameOrEmail": "john.doe@example.com" }
 *
 *  Response:
 *  { "message": "Friend request sent" }
 *
 *  GET /api/friends/list
 *  Response:
 *  [
 *      { "username": "john_doe", "email": "john.doe@example.com" },
 *      { "username": "jane_doe", "email": "jane.doe@example.com" }
 *  ]
 *  ```
 *
 *  @dependencies
 *  - services.FriendServiceInterface: Interface for managing friend-related operations.
 *  - utils: Utility package for writing JSON responses and errors.
 *
 *  @file      friend_handler.go
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

// FriendHandler struct for handling friend-related requests.
type FriendHandler struct {
	FriendService services.FriendServiceInterface
}

// NewFriendHandler initializes a new FriendHandler instance.
func NewFriendHandler(fs services.FriendServiceInterface) *FriendHandler {
	return &FriendHandler{FriendService: fs}
}

// SendFriendRequest handles POST requests to send a friend request to a user.
func (fh *FriendHandler) SendFriendRequest(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		UsernameOrEmail string `json:"usernameOrEmail"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		utils.WriteJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if requestData.UsernameOrEmail == "" {
		utils.WriteJSONError(w, "Username or Email is required", http.StatusBadRequest)
		return
	}

	userEmail, ok := r.Context().Value("userEmail").(string)
	if !ok {
		utils.WriteJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := fh.FriendService.SendFriendRequest(r.Context(), userEmail, requestData.UsernameOrEmail)
	if err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.WriteJSON(w, map[string]string{"message": "Friend request sent"})
}

// AcceptFriendRequest handles POST requests to accept a friend request.
func (fh *FriendHandler) AcceptFriendRequest(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		UsernameOrEmail string `json:"usernameOrEmail"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		utils.WriteJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if requestData.UsernameOrEmail == "" {
		utils.WriteJSONError(w, "Username or Email is required", http.StatusBadRequest)
		return
	}

	userEmail, ok := r.Context().Value("userEmail").(string)
	if !ok {
		utils.WriteJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := fh.FriendService.AcceptFriendRequest(r.Context(), userEmail, requestData.UsernameOrEmail)
	if err != nil {
		switch err.Error() {
		case "User not found", "Friend request not found":
			utils.WriteJSONError(w, err.Error(), http.StatusNotFound)
		default:
			utils.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	utils.WriteJSON(w, map[string]string{"message": "Friend request accepted"})
}

// GetFriendsList handles GET requests to fetch the authenticated user's friends list.
func (fh *FriendHandler) GetFriendsList(w http.ResponseWriter, r *http.Request) {
	userEmail := r.Context().Value("userEmail").(string)

	friends, err := fh.FriendService.GetFriendsList(r.Context(), userEmail)
	if err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, friends)
}

// RemoveFriend handles DELETE requests to remove a friend from the user's friend list.
func (fh *FriendHandler) RemoveFriend(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		utils.WriteJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userEmail := r.Context().Value("userEmail").(string)

	if err := fh.FriendService.RemoveFriend(r.Context(), userEmail, requestData.Username); err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, map[string]string{"message": "Friend removed"})
}

// GetPendingFriendRequests handles GET requests to fetch pending friend requests for the user.
func (fh *FriendHandler) GetPendingFriendRequests(w http.ResponseWriter, r *http.Request) {
	userEmail := r.Context().Value("userEmail").(string)

	requests, err := fh.FriendService.GetPendingFriendRequests(r.Context(), userEmail)
	if err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, requests)
}

// DeclineFriendRequest handles POST requests to decline a friend request.
func (fh *FriendHandler) DeclineFriendRequest(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		UsernameOrEmail string `json:"usernameOrEmail"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		utils.WriteJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if requestData.UsernameOrEmail == "" {
		utils.WriteJSONError(w, "Username or Email is required", http.StatusBadRequest)
		return
	}

	userEmail, ok := r.Context().Value("userEmail").(string)
	if !ok {
		utils.WriteJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := fh.FriendService.DeclineFriendRequest(r.Context(), userEmail, requestData.UsernameOrEmail)
	if err != nil {
		switch err.Error() {
		case "User not found", "Friend request not found":
			utils.WriteJSONError(w, err.Error(), http.StatusNotFound)
		default:
			utils.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	utils.WriteJSON(w, map[string]string{"message": "Friend request declined"})
}

// CancelFriendRequest handles DELETE requests to cancel a sent friend request.
func (fh *FriendHandler) CancelFriendRequest(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		utils.WriteJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userEmail := r.Context().Value("userEmail").(string)

	if err := fh.FriendService.CancelFriendRequest(r.Context(), userEmail, requestData.Username); err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, map[string]string{"message": "Friend request canceled"})
}
