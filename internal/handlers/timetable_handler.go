// internal/handlers/timetable_handler.go
package handlers

import (
	"encoding/json"
	"net/http"

	"proh2052-group6/internal/services"
	"proh2052-group6/pkg/utils"
)

type TimetableHandler struct {
	TimetableService services.TimetableServiceInterface
}

func NewTimetableHandler(ts services.TimetableServiceInterface) *TimetableHandler {
	return &TimetableHandler{TimetableService: ts}
}

func (th *TimetableHandler) ImportTimetable(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		ICSContent string `json:"icsContent"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		utils.WriteJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if requestData.ICSContent == "" {
		utils.WriteJSONError(w, "ICS content is required", http.StatusBadRequest)
		return
	}

	userEmail, ok := r.Context().Value("userEmail").(string)
	if !ok {
		utils.WriteJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := th.TimetableService.ImportTimetable(r.Context(), userEmail, requestData.ICSContent)
	if err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, map[string]string{
		"message": "Timetable imported successfully",
	})
}
