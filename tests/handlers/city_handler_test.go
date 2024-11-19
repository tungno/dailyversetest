// tests/handlers/city_handler_test.go
package handlers_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"proh2052-group6/internal/handlers"
	"proh2052-group6/internal/services"
)

func TestCityHandler_GetCities(t *testing.T) {
	// Setup a test server to mock the external API
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify that the request method is POST
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		// Read the request body
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Failed to read request body: %v", err)
		}
		defer r.Body.Close()
		// Parse the request body
		var requestBody map[string]string
		err = json.Unmarshal(bodyBytes, &requestBody)
		if err != nil {
			t.Fatalf("Failed to unmarshal request body: %v", err)
		}
		// Check that the country parameter is correct
		if requestBody["country"] != "TestCountry" {
			t.Errorf("Expected country 'TestCountry', got '%s'", requestBody["country"])
		}
		// Respond with a mocked city list
		cityResponse := struct {
			Error bool     `json:"error"`
			Msg   string   `json:"msg"`
			Data  []string `json:"data"`
		}{
			Error: false,
			Msg:   "Success",
			Data:  []string{"City1", "City2", "City3"},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cityResponse)
	}))
	defer testServer.Close()

	// Replace the CitiesAPIURL to point to our test server
	originalCitiesAPIURL := services.CitiesAPIURL
	services.CitiesAPIURL = testServer.URL
	defer func() {
		services.CitiesAPIURL = originalCitiesAPIURL
	}()

	// Create the handler
	cityHandler := handlers.NewCityHandler()

	// Create a test request
	req, err := http.NewRequest("GET", "/api/cities?country=TestCountry", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a ResponseRecorder
	rr := httptest.NewRecorder()

	// Call the handler
	http.HandlerFunc(cityHandler.GetCities).ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body
	var cities []string
	err = json.NewDecoder(rr.Body).Decode(&cities)
	if err != nil {
		t.Errorf("Failed to decode response body: %v", err)
	}

	expectedCities := []string{"City1", "City2", "City3"}
	if !equalStringSlices(cities, expectedCities) {
		t.Errorf("Expected cities %v, got %v", expectedCities, cities)
	}
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestCityHandler_GetCities_MissingCountry(t *testing.T) {
	// Create the handler
	cityHandler := handlers.NewCityHandler()

	// Create a test request without the country parameter
	req, err := http.NewRequest("GET", "/api/cities", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a ResponseRecorder
	rr := httptest.NewRecorder()

	// Call the handler
	http.HandlerFunc(cityHandler.GetCities).ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	// Check the response body
	expectedError := "Missing country parameter\n"
	if rr.Body.String() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, rr.Body.String())
	}
}

func TestCityHandler_GetCities_ExternalAPIError(t *testing.T) {
	// Setup a test server to mock the external API
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Respond with an error
		cityResponse := struct {
			Error bool     `json:"error"`
			Msg   string   `json:"msg"`
			Data  []string `json:"data"`
		}{
			Error: true,
			Msg:   "Country not found",
			Data:  nil,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cityResponse)
	}))
	defer testServer.Close()

	// Replace the CitiesAPIURL to point to our test server
	originalCitiesAPIURL := services.CitiesAPIURL
	services.CitiesAPIURL = testServer.URL
	defer func() {
		services.CitiesAPIURL = originalCitiesAPIURL
	}()

	// Create the handler
	cityHandler := handlers.NewCityHandler()

	// Create a test request
	req, err := http.NewRequest("GET", "/api/cities?country=UnknownCountry", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a ResponseRecorder
	rr := httptest.NewRecorder()

	// Call the handler
	http.HandlerFunc(cityHandler.GetCities).ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	// Check the response body
	expectedError := "Error fetching cities\n"
	if rr.Body.String() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, rr.Body.String())
	}
}
