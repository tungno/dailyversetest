// internal/handlers/profile_handler.go
package handlers

import (
	"encoding/json"
	"net/http"

	"proh2052-group6/internal/services"
	"proh2052-group6/pkg/utils"
)

type ProfileHandler struct {
	ProfileService services.ProfileServiceInterface
}

func NewProfileHandler(ps services.ProfileServiceInterface) *ProfileHandler {
	return &ProfileHandler{ProfileService: ps}
}

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

func (ph *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userEmail := r.Context().Value("userEmail").(string)

	profileData, err := ph.ProfileService.GetProfile(r.Context(), userEmail)
	if err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, profileData)
}

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
