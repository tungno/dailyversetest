// tests/handlers/profile_handler_test.go
package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"proh2052-group6/internal/handlers"
	"proh2052-group6/tests/mocks"
	"testing"
)

func TestProfileHandler_GetProfile(t *testing.T) {
	// Set up mock profile service
	mockProfileService := mocks.NewMockProfileService()
	userEmail := "test@example.com"
	mockProfileService.Profiles[userEmail] = map[string]interface{}{
		"Email":    userEmail,
		"Username": "testuser",
		"Country":  "TestCountry",
		"City":     "TestCity",
		"Password": "hashedpassword123", // For password verification
	}

	// Create the profile handler
	profileHandler := handlers.NewProfileHandler(mockProfileService)

	// Create a test HTTP request
	req, err := http.NewRequest("GET", "/api/profile", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Set the userEmail in the context
	ctx := context.WithValue(req.Context(), "userEmail", userEmail)
	req = req.WithContext(ctx)

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler
	handler := http.HandlerFunc(profileHandler.ProfileHandler)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Parse the response body
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	// Verify the response data
	if response["Email"] != userEmail {
		t.Errorf("Expected Email '%s', got '%s'", userEmail, response["Email"])
	}
	if response["Username"] != "testuser" {
		t.Errorf("Expected Username 'testuser', got '%s'", response["Username"])
	}
	if response["Country"] != "TestCountry" {
		t.Errorf("Expected Country 'TestCountry', got '%s'", response["Country"])
	}
	if response["City"] != "TestCity" {
		t.Errorf("Expected City 'TestCity', got '%s'", response["City"])
	}
}

func TestProfileHandler_UpdateProfile(t *testing.T) {
	// Set up mock profile service
	mockProfileService := mocks.NewMockProfileService()
	userEmail := "test@example.com"
	mockProfileService.Profiles[userEmail] = map[string]interface{}{
		"Email":    userEmail,
		"Username": "testuser",
		"Country":  "TestCountry",
		"City":     "TestCity",
		"Password": "currentpassword", // Simulated current password
	}

	// Create the profile handler
	profileHandler := handlers.NewProfileHandler(mockProfileService)

	// Prepare the updated data
	updatedData := map[string]interface{}{
		"Username":        "updateduser",
		"Country":         "UpdatedCountry",
		"CurrentPassword": "currentpassword",
		"NewPassword":     "newsecurepassword",
	}
	requestBody, _ := json.Marshal(updatedData)

	// Create a test HTTP request
	req, err := http.NewRequest("PUT", "/api/profile", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Set the userEmail in the context
	ctx := context.WithValue(req.Context(), "userEmail", userEmail)
	req = req.WithContext(ctx)

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler
	handler := http.HandlerFunc(profileHandler.ProfileHandler)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Parse the response body
	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	expectedMessage := "Successfully updated profile"
	if response["message"] != expectedMessage {
		t.Errorf("Expected message '%s', got '%s'", expectedMessage, response["message"])
	}

	// Verify that the profile was updated
	updatedProfile := mockProfileService.Profiles[userEmail]
	if updatedProfile["Username"] != "updateduser" {
		t.Errorf("Expected Username 'updateduser', got '%s'", updatedProfile["Username"])
	}
	if updatedProfile["Country"] != "UpdatedCountry" {
		t.Errorf("Expected Country 'UpdatedCountry', got '%s'", updatedProfile["Country"])
	}
	if updatedProfile["Password"] != "newsecurepassword" {
		t.Errorf("Expected Password 'newsecurepassword', got '%s'", updatedProfile["Password"])
	}
}

func TestProfileHandler_UpdateProfile_InvalidCurrentPassword(t *testing.T) {
	// Set up mock profile service
	mockProfileService := mocks.NewMockProfileService()
	userEmail := "test@example.com"
	mockProfileService.Profiles[userEmail] = map[string]interface{}{
		"Email":    userEmail,
		"Username": "testuser",
		"Country":  "TestCountry",
		"City":     "TestCity",
		"Password": "correctpassword",
	}

	// Create the profile handler
	profileHandler := handlers.NewProfileHandler(mockProfileService)

	// Prepare the updated data with incorrect current password
	updatedData := map[string]interface{}{
		"Username":        "updateduser",
		"CurrentPassword": "wrongpassword",
	}
	requestBody, _ := json.Marshal(updatedData)

	// Create a test HTTP request
	req, err := http.NewRequest("PUT", "/api/profile", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Set the userEmail in the context
	ctx := context.WithValue(req.Context(), "userEmail", userEmail)
	req = req.WithContext(ctx)

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler
	handler := http.HandlerFunc(profileHandler.ProfileHandler)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	// Verify the error message
	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	expectedError := "Invalid current password"
	if response["error"] != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, response["error"])
	}
}

func TestProfileHandler_ProfileHandler_MethodNotAllowed(t *testing.T) {
	// Set up mock profile service
	mockProfileService := mocks.NewMockProfileService()

	// Create the profile handler
	profileHandler := handlers.NewProfileHandler(mockProfileService)

	// Create a test HTTP request with unsupported method
	req, err := http.NewRequest("POST", "/api/profile", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler
	handler := http.HandlerFunc(profileHandler.ProfileHandler)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}
}