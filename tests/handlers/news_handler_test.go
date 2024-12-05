/**
 *  TestNewsHandler_FetchNews validates the functionality of the NewsHandler's FetchNews method.
 *  It simulates fetching news articles based on user preferences and verifies the response.
 *
 *  @dependencies
 *  - mocks.NewMockUserRepository: Mock repository for simulating user data.
 *  - httptest.NewServer: Creates a mock server to simulate the external news API.
 *  - services.NewsService: News service for business logic, injected with mocks and test configurations.
 *  - handlers.NewsHandler: HTTP handler for handling news requests.
 *
 *  @behaviors
 *  - Fetches news from a mock API based on user preferences retrieved via the mock UserRepository.
 *  - Verifies that the HTTP request and response conform to expected behavior.
 *
 *  @testcases
 *  - Validates the HTTP status code (expected: 200 OK).
 *  - Ensures the response body contains the correct news data.
 *  - Simulates a real-world scenario using a mock external news API and user data.
 *
 *  @example
 *  ```
 *  // Simulate fetching local news for a user
 *  req, _ := http.NewRequest("GET", "/api/news?mode=local", nil)
 *  ctx := context.WithValue(req.Context(), "userEmail", "test@example.com")
 *  req = req.WithContext(ctx)
 *
 *  rr := httptest.NewRecorder()
 *  handler := http.HandlerFunc(newsHandler.FetchNews)
 *  handler.ServeHTTP(rr, req)
 *  ```
 */
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
	// Step 1: Setup mock UserRepository with test data
	mockUsers := map[string]*models.User{
		"test@example.com": {
			Email:   "test@example.com",
			Country: "TestCountry",
		},
	}
	mockUserRepo := mocks.NewMockUserRepository(mockUsers)

	// Step 2: Create a mock server to simulate the external news API
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mock response from the external API
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

	// Step 3: Setup the NewsService with mock repositories and test server
	newsService := &services.NewsService{
		UserRepo:   mockUserRepo,
		HTTPClient: testServer.Client(),
		NewsAPIURL: testServer.URL,
		GetCountryAndLanguageCode: func(countryName string) (string, string, error) {
			// Mock implementation to return a hardcoded country and language code
			return "testcountrycode", "en", nil
		},
	}

	// Step 4: Initialize the NewsHandler
	newsHandler := handlers.NewNewsHandler(newsService)

	// Step 5: Create a test HTTP request for the FetchNews endpoint
	req, err := http.NewRequest("GET", "/api/news?mode=local", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Set the userEmail in the request context to simulate authentication
	userEmail := "test@example.com"
	ctx := context.WithValue(req.Context(), "userEmail", userEmail)
	req = req.WithContext(ctx)

	// Step 6: Create a ResponseRecorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Step 7: Invoke the handler
	handler := http.HandlerFunc(newsHandler.FetchNews)
	handler.ServeHTTP(rr, req)

	// Step 8: Validate the HTTP status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Step 9: Parse and validate the response body
	var response []map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Errorf("Failed to decode response body: %v", err)
	}

	// Verify the number of news items
	if len(response) != 1 {
		t.Errorf("Expected 1 news item, got %d", len(response))
	}

	// Validate the content of the news item
	if response[0]["title"] != "Test News Title" {
		t.Errorf("Expected news title 'Test News Title', got '%s'", response[0]["title"])
	}
}
