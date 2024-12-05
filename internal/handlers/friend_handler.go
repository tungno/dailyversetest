// internal/handlers/friend_handler.go
package handlers

import (
	"encoding/json"
	"net/http"

	"proh2052-group6/internal/services"
	"proh2052-group6/pkg/utils"
)

type FriendHandler struct {
	FriendService services.FriendServiceInterface
}

func NewFriendHandler(fs services.FriendServiceInterface) *FriendHandler {
	return &FriendHandler{FriendService: fs}
}

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
		// Determine the type of error and return appropriate status code
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

func (fh *FriendHandler) GetFriendsList(w http.ResponseWriter, r *http.Request) {
	userEmail := r.Context().Value("userEmail").(string)

	friends, err := fh.FriendService.GetFriendsList(r.Context(), userEmail)
	if err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, friends)
}

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

func (fh *FriendHandler) GetPendingFriendRequests(w http.ResponseWriter, r *http.Request) {
	userEmail := r.Context().Value("userEmail").(string)

	requests, err := fh.FriendService.GetPendingFriendRequests(r.Context(), userEmail)
	if err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, requests)
}

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
		// Determine the type of error and return appropriate status code
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
