package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"proh2052-group6/internal/handlers"
	"proh2052-group6/internal/services"
	"proh2052-group6/pkg/models"
	"proh2052-group6/tests/mocks"
)

func TestNewsHandler_FetchNews(t *testing.T) {
	// Create a mock user repository
	mockUsers := map[string]*models.User{
		"test@example.com": {
			Email:   "test@example.com",
			Country: "TestCountry",
		},
	}
	mockUserRepo := mocks.NewMockUserRepository(mockUsers)

	// Setup a test server to mock the external news API
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// we can verify the request parameters if needed
		// For example, check the API key or query parameters

		// Respond with a mocked news response
		newsResponse := struct {
			Status       string                   `json:"status"`
			TotalResults int                      `json:"totalResults"`
			Results      []map[string]interface{} `json:"results"`
		}{
			Status:       "success",
			TotalResults: 1,
			Results: []map[string]interface{}{
				{
					"title":       "Test News Title",
					"description": "Test News Description",
					"link":        "https://example.com/news",
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(newsResponse)
	}))
	defer testServer.Close()

	// Create the news service with the mock user repo and custom HTTP client
	newsService := &services.NewsService{
		UserRepo:   mockUserRepo,
		HTTPClient: testServer.Client(),
		NewsAPIURL: testServer.URL,
		GetCountryAndLanguageCode: func(countryName string) (string, string, error) {
			// Mock implementation
			return "testcountrycode", "en", nil
		},
	}

	// Create the news handler
	newsHandler := handlers.NewNewsHandler(newsService)

	// Create a test HTTP request
	req, err := http.NewRequest("GET", "/api/news?mode=local", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Set the userEmail in the context
	userEmail := "test@example.com"
	ctx := context.WithValue(req.Context(), "userEmail", userEmail)
	req = req.WithContext(ctx)

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler
	handler := http.HandlerFunc(newsHandler.FetchNews)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body
	var response []map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Errorf("Failed to decode response body: %v", err)
	}

	if len(response) != 1 {
		t.Errorf("Expected 1 news item, got %d", len(response))
	}

	if response[0]["title"] != "Test News Title" {
		t.Errorf("Expected news title 'Test News Title', got '%s'", response[0]["title"])
	}
}
