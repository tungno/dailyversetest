/**
 *  JournalHandler Tests validate the behavior of the JournalHandler methods.
 *  They use a mock JournalService to isolate the handler logic and verify interactions with the service.
 *
 *  @file       journal_handler_test.go
 *  @package    handlers_test
 *
 *  @test_cases
 *  - TestJournalHandler_CreateJournal      - Tests creating a new journal entry.
 *  - TestJournalHandler_GetJournal         - Tests retrieving a specific journal entry.
 *  - TestJournalHandler_UpdateJournal      - Tests updating an existing journal entry.
 *  - TestJournalHandler_DeleteJournal      - Tests deleting a journal entry.
 *  - TestJournalHandler_GetAllJournals     - Tests retrieving all journal entries for a user.
 *
 *  @dependencies
 *  - mocks.NewMockJournalService: Mock implementation of JournalService for testing.
 *  - httptest: Provides utilities for testing HTTP handlers.
 *  - context.WithValue: Adds user-specific context values for testing purposes.
 *  - encoding/json: Handles JSON marshalling and unmarshalling.
 *
 *  @behaviors
 *  - Verifies HTTP status codes for each handler.
 *  - Validates request/response data consistency.
 *  - Confirms the correct service methods are called during handler execution.
 *
 *  @example
 *  ```
 *  req := httptest.NewRequest("POST", "/api/journal/save", body)
 *  rr := httptest.NewRecorder()
 *  handler := http.HandlerFunc(journalHandler.CreateJournal)
 *  handler.ServeHTTP(rr, req)
 *  ```
 *
 *  @authors
 *      - Aayush
 *      - Tung
 *      - Boss
 *      - Majd
 */

package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"proh2052-group6/internal/handlers"
	"proh2052-group6/pkg/models"
	"proh2052-group6/tests/mocks"
)

func TestJournalHandler_CreateJournal(t *testing.T) {
	// Setup mock JournalService and JournalHandler
	mockJournalService := mocks.NewMockJournalService()
	journalHandler := handlers.NewJournalHandler(mockJournalService)

	// Prepare request body
	journal := models.Journal{
		Date:    "2023-10-15",
		Content: "Today was a good day.",
	}
	requestBody, _ := json.Marshal(journal)

	// Create HTTP request
	req, err := http.NewRequest("POST", "/api/journal/save", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Inject userEmail into context
	userEmail := "test@example.com"
	ctx := context.WithValue(req.Context(), "userEmail", userEmail)
	req = req.WithContext(ctx)

	// Create ResponseRecorder to capture response
	rr := httptest.NewRecorder()

	// Invoke handler
	handler := http.HandlerFunc(journalHandler.CreateJournal)
	handler.ServeHTTP(rr, req)

	// Assert status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Parse and validate response
	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	expectedMessage := "Journal created successfully"
	if response["message"] != expectedMessage {
		t.Errorf("Expected message '%s', got '%s'", expectedMessage, response["message"])
	}
	if response["journalID"] == "" {
		t.Errorf("Expected journalID in response")
	}

	// Verify journal saved in mock service
	savedJournal, err := mockJournalService.GetJournal(context.Background(), userEmail, response["journalID"])
	if err != nil {
		t.Errorf("Journal was not saved in the service: %v", err)
	}
	if savedJournal.Content != journal.Content {
		t.Errorf("Expected journal content '%s', got '%s'", journal.Content, savedJournal.Content)
	}
}

func TestJournalHandler_GetJournal(t *testing.T) {
	// Setup mock JournalService and JournalHandler
	mockJournalService := mocks.NewMockJournalService()
	journalHandler := handlers.NewJournalHandler(mockJournalService)

	// Add a journal to the mock service
	userEmail := "test@example.com"
	journalID := "journal123"
	journal := &models.Journal{
		JournalID: journalID,
		Email:     userEmail,
		Date:      "2023-10-15",
		Content:   "Today was a good day.",
	}
	mockJournalService.Journals[journalID] = journal

	// Create HTTP request
	req, err := http.NewRequest("GET", "/api/journal?journalID="+journalID, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Inject userEmail into context
	ctx := context.WithValue(req.Context(), "userEmail", userEmail)
	req = req.WithContext(ctx)

	// Create ResponseRecorder
	rr := httptest.NewRecorder()

	// Invoke handler
	handler := http.HandlerFunc(journalHandler.GetJournal)
	handler.ServeHTTP(rr, req)

	// Assert status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Parse and validate response
	var response models.Journal
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	if response.JournalID != journalID {
		t.Errorf("Expected journalID '%s', got '%s'", journalID, response.JournalID)
	}
	if response.Content != journal.Content {
		t.Errorf("Expected journal content '%s', got '%s'", journal.Content, response.Content)
	}
}

func TestJournalHandler_UpdateJournal(t *testing.T) {
	// Create a mock journal service
	mockJournalService := mocks.NewMockJournalService()
	journalHandler := handlers.NewJournalHandler(mockJournalService)

	// Add a journal to the mock service
	userEmail := "test@example.com"
	journalID := "journal123"
	journal := &models.Journal{
		JournalID: journalID,
		Email:     userEmail,
		Date:      "2023-10-15",
		Content:   "Today was a good day.",
	}
	mockJournalService.Journals[journalID] = journal

	// Prepare the updated journal data
	updatedJournal := models.Journal{
		Content: "Updated journal content.",
	}
	requestBody, _ := json.Marshal(updatedJournal)

	// Create a new HTTP request
	req, err := http.NewRequest("PUT", "/api/journal/update?journalID="+journalID, bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Set the userEmail in the context
	ctx := context.WithValue(req.Context(), "userEmail", userEmail)
	req = req.WithContext(ctx)

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler
	handler := http.HandlerFunc(journalHandler.UpdateJournal)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Verify that the journal was updated in the mock service
	updatedJournalInService, err := mockJournalService.GetJournal(context.Background(), userEmail, journalID)
	if err != nil {
		t.Errorf("Journal was not found in the service: %v", err)
	}

	if updatedJournalInService.Content != updatedJournal.Content {
		t.Errorf("Expected journal content '%s', got '%s'", updatedJournal.Content, updatedJournalInService.Content)
	}
}

func TestJournalHandler_DeleteJournal(t *testing.T) {
	// Create a mock journal service
	mockJournalService := mocks.NewMockJournalService()
	journalHandler := handlers.NewJournalHandler(mockJournalService)

	// Add a journal to the mock service
	userEmail := "test@example.com"
	journalID := "journal123"
	journal := &models.Journal{
		JournalID: journalID,
		Email:     userEmail,
		Date:      "2023-10-15",
		Content:   "Today was a good day.",
	}
	mockJournalService.Journals[journalID] = journal

	// Create a new HTTP request
	req, err := http.NewRequest("DELETE", "/api/journal/delete?journalID="+journalID, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Set the userEmail in the context
	ctx := context.WithValue(req.Context(), "userEmail", userEmail)
	req = req.WithContext(ctx)

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler
	handler := http.HandlerFunc(journalHandler.DeleteJournal)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Verify that the journal was deleted from the mock service
	_, err = mockJournalService.GetJournal(context.Background(), userEmail, journalID)
	if err == nil {
		t.Errorf("Expected journal to be deleted, but it still exists")
	}
}

func TestJournalHandler_GetAllJournals(t *testing.T) {
	// Create a mock journal service
	mockJournalService := mocks.NewMockJournalService()
	journalHandler := handlers.NewJournalHandler(mockJournalService)

	// Add journals to the mock service
	userEmail := "test@example.com"
	journal1 := &models.Journal{
		JournalID: "journal1",
		Email:     userEmail,
		Date:      "2023-10-15",
		Content:   "First journal entry.",
	}
	journal2 := &models.Journal{
		JournalID: "journal2",
		Email:     userEmail,
		Date:      "2023-10-16",
		Content:   "Second journal entry.",
	}
	mockJournalService.Journals[journal1.JournalID] = journal1
	mockJournalService.Journals[journal2.JournalID] = journal2

	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/api/journals", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Set the userEmail in the context
	ctx := context.WithValue(req.Context(), "userEmail", userEmail)
	req = req.WithContext(ctx)

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	// Call the handler
	handler := http.HandlerFunc(journalHandler.GetAllJournals)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body
	var response []models.Journal
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}
	if len(response) != 2 {
		t.Errorf("Expected 2 journals, got %d", len(response))
	}
}
