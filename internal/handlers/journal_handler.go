// internal/handlers/journal_handler.go
package handlers

import (
	"encoding/json"
	"net/http"

	"proh2052-group6/internal/services"
	"proh2052-group6/pkg/models"
	"proh2052-group6/pkg/utils"
)

type JournalHandler struct {
	JournalService services.JournalServiceInterface
}

func NewJournalHandler(js services.JournalServiceInterface) *JournalHandler {
	return &JournalHandler{JournalService: js}
}

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

func (jh *JournalHandler) GetAllJournals(w http.ResponseWriter, r *http.Request) {
	userEmail := r.Context().Value("userEmail").(string)

	journals, err := jh.JournalService.GetAllJournals(r.Context(), userEmail)
	if err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, journals)
}
