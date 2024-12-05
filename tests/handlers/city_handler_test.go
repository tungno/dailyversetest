/**
 *  CityHandler Test Suite
 *
 *  This test suite validates the functionality of the CityHandler, ensuring that it:
 *  - Correctly fetches cities when a valid 'country' parameter is provided.
 *  - Returns an error when the 'country' parameter is missing.
 *  - Handles errors from the CityService gracefully and returns appropriate status codes.
 *
 *  @dependencies
 *  - mocks.MockCityService: Mock implementation of the CityService for testing.
 *  - mocks.MockUserService: Mock implementation of the UserService for dependency injection.
 *  - testify/assert: Library for test assertions.
 *
 *  @file      city_handler_test.go
 *  @project   DailyVerse
 *  @framework Go HTTP Testing with Testify
 */

package handlers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"proh2052-group6/internal/handlers"
	"proh2052-group6/tests/mocks"

	"github.com/stretchr/testify/assert"
)

func TestCityHandler_GetCities_WithCountryParam(t *testing.T) {
	// Test Case: Fetch cities successfully with a valid 'country' parameter.

	// Setup mock CityService with expected behavior.
	mockCityService := &mocks.MockCityService{
		GetCitiesByCountryFunc: func(country string) ([]string, error) {
			if country == "TestCountry" {
				return []string{"City1", "City2", "City3"}, nil
			}
			return nil, fmt.Errorf("invalid country")
		},
	}

	// Mock UserService (not used in this test).
	mockUserService := &mocks.MockUserService{}

	// Initialize CityHandler with mocks.
	cityHandler := handlers.NewCityHandler(mockCityService, mockUserService)

	// Create a test HTTP request with the 'country' parameter.
	req, err := http.NewRequest("GET", "/api/cities?country=TestCountry", nil)
	assert.NoError(t, err, "Failed to create request")

	// Create a ResponseRecorder to capture the response.
	rr := httptest.NewRecorder()

	// Invoke the handler.
	http.HandlerFunc(cityHandler.GetCities).ServeHTTP(rr, req)

	// Validate the response.
	assert.Equal(t, http.StatusOK, rr.Code, "Handler should return status 200 OK")

	// Parse and verify the JSON response.
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err, "Response should be valid JSON")

	cities, ok := response["data"].([]interface{})
	assert.True(t, ok, "Expected 'data' to be an array")

	expectedCities := []string{"City1", "City2", "City3"}
	assert.Equal(t, len(expectedCities), len(cities), "Number of cities mismatch")

	for i, city := range expectedCities {
		assert.Equal(t, city, cities[i].(string), fmt.Sprintf("Expected city '%s', got '%s'", city, cities[i].(string)))
	}
}

func TestCityHandler_GetCities_WithoutCountryParam(t *testing.T) {
	// Test Case: Return error when the 'country' parameter is missing.

	// Setup mock services (not used in this test).
	mockCityService := &mocks.MockCityService{}
	mockUserService := &mocks.MockUserService{}

	// Initialize CityHandler with mocks.
	cityHandler := handlers.NewCityHandler(mockCityService, mockUserService)

	// Create a test HTTP request without the 'country' parameter.
	req, err := http.NewRequest("GET", "/api/cities", nil)
	assert.NoError(t, err, "Failed to create request")

	// Create a ResponseRecorder to capture the response.
	rr := httptest.NewRecorder()

	// Invoke the handler.
	http.HandlerFunc(cityHandler.GetCities).ServeHTTP(rr, req)

	// Validate the response.
	assert.Equal(t, http.StatusBadRequest, rr.Code, "Handler should return status 400 Bad Request")

	// Validate the error message.
	expectedError := "Missing country parameter\n"
	assert.Equal(t, expectedError, rr.Body.String(), "Error message should match")
}

func TestCityHandler_GetCities_ExternalAPIError(t *testing.T) {
	// Test Case: Handle errors from the CityService gracefully.

	// Setup mock CityService to return an error.
	mockCityService := &mocks.MockCityService{
		GetCitiesByCountryFunc: func(country string) ([]string, error) {
			return nil, fmt.Errorf("error fetching cities: country not found")
		},
	}

	// Mock UserService (not used in this test).
	mockUserService := &mocks.MockUserService{}

	// Initialize CityHandler with mocks.
	cityHandler := handlers.NewCityHandler(mockCityService, mockUserService)

	// Create a test HTTP request with an invalid 'country' parameter.
	req, err := http.NewRequest("GET", "/api/cities?country=UnknownCountry", nil)
	assert.NoError(t, err, "Failed to create request")

	// Create a ResponseRecorder to capture the response.
	rr := httptest.NewRecorder()

	// Invoke the handler.
	http.HandlerFunc(cityHandler.GetCities).ServeHTTP(rr, req)

	// Validate the response.
	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Handler should return status 500 Internal Server Error")

	// Validate the error message.
	expectedError := "Error fetching cities\n"
	assert.Equal(t, expectedError, rr.Body.String(), "Error message should match")
}
