// tests/handlers/city_handler_test.go
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
	// Setup mock CityService
	mockCityService := &mocks.MockCityService{
		GetCitiesByCountryFunc: func(country string) ([]string, error) {
			if country == "TestCountry" {
				return []string{"City1", "City2", "City3"}, nil
			}
			return nil, fmt.Errorf("invalid country")
		},
	}

	// Setup mock UserService (not used since country is provided via query param)
	mockUserService := &mocks.MockUserService{
		// No methods need to be set for this test
	}

	// Initialize handler with mocks
	cityHandler := handlers.NewCityHandler(mockCityService, mockUserService)

	// Create test request with 'country' parameter
	req, err := http.NewRequest("GET", "/api/cities?country=TestCountry", nil)
	assert.NoError(t, err, "Failed to create request")

	// Create ResponseRecorder
	rr := httptest.NewRecorder()

	// Call handler
	http.HandlerFunc(cityHandler.GetCities).ServeHTTP(rr, req)

	// Check status code
	assert.Equal(t, http.StatusOK, rr.Code, "Handler should return status 200 OK")

	// Check response body
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
	// Setup mock CityService (not used since country is provided)
	mockCityService := &mocks.MockCityService{
		GetCitiesByCountryFunc: func(country string) ([]string, error) {
			// This should not be called in this test since 'country' is provided
			return []string{}, nil
		},
	}

	// Setup mock UserService (not used since 'country' is provided)
	mockUserService := &mocks.MockUserService{}

	// Initialize handler with mocks
	cityHandler := handlers.NewCityHandler(mockCityService, mockUserService)

	// Create test request without 'country' parameter
	req, err := http.NewRequest("GET", "/api/cities", nil)
	assert.NoError(t, err, "Failed to create request")

	// Create ResponseRecorder
	rr := httptest.NewRecorder()

	// Call handler
	http.HandlerFunc(cityHandler.GetCities).ServeHTTP(rr, req)

	// Check status code
	assert.Equal(t, http.StatusBadRequest, rr.Code, "Handler should return status 400 Bad Request")

	// Check response body
	expectedError := "Missing country parameter\n"
	assert.Equal(t, expectedError, rr.Body.String(), "Error message should match")
}

func TestCityHandler_GetCities_ExternalAPIError(t *testing.T) {
	// Setup mock CityService to return an error
	mockCityService := &mocks.MockCityService{
		GetCitiesByCountryFunc: func(country string) ([]string, error) {
			return nil, fmt.Errorf("error fetching cities: country not found")
		},
	}

	// Setup mock UserService (not used since country is provided via query param)
	mockUserService := &mocks.MockUserService{
		// No methods need to be set for this test
	}

	// Initialize handler with mocks
	cityHandler := handlers.NewCityHandler(mockCityService, mockUserService)

	// Create test request with invalid 'country' parameter
	req, err := http.NewRequest("GET", "/api/cities?country=UnknownCountry", nil)
	assert.NoError(t, err, "Failed to create request")

	// Create ResponseRecorder
	rr := httptest.NewRecorder()

	// Call handler
	http.HandlerFunc(cityHandler.GetCities).ServeHTTP(rr, req)

	// Check status code
	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Handler should return status 500 Internal Server Error")

	// Check response body
	expectedError := "Error fetching cities\n"
	assert.Equal(t, expectedError, rr.Body.String(), "Error message should match")
}
