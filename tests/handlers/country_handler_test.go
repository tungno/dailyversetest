// tests/handlers/country_handler_test.go
package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"proh2052-group6/internal/handlers"
	"proh2052-group6/internal/services"
)

func TestCountryHandler_GetCountries(t *testing.T) {
	// Setup a test server to mock the external API
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Respond with a mocked country list
		countriesData := []struct {
			Name struct {
				Common string `json:"common"`
			} `json:"name"`
			CCA2 string `json:"cca2"`
		}{
			{
				Name: struct {
					Common string `json:"common"`
				}{Common: "Canada"},
				CCA2: "CA",
			},
			{
				Name: struct {
					Common string `json:"common"`
				}{Common: "Cameroon"},
				CCA2: "CM",
			},
			{
				Name: struct {
					Common string `json:"common"`
				}{Common: "Cambodia"},
				CCA2: "KH",
			},
			{
				Name: struct {
					Common string `json:"common"`
				}{Common: "France"},
				CCA2: "FR",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(countriesData)
	}))
	defer testServer.Close()

	// Replace the CountriesAPIURL to point to our test server
	originalCountriesAPIURL := services.CountriesAPIURL
	services.SetCountriesAPIURL(testServer.URL)
	defer services.SetCountriesAPIURL(originalCountriesAPIURL)

	// Create the handler
	countryHandler := handlers.NewCountryHandler()

	// Create a test request with a search query
	req, err := http.NewRequest("GET", "/api/countries?search=cam", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler
	http.HandlerFunc(countryHandler.GetCountries).ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body
	var countries []services.Country
	err = json.Unmarshal(rr.Body.Bytes(), &countries)
	if err != nil {
		t.Errorf("Failed to decode response body: %v", err)
	}

	expectedCountries := []services.Country{
		{Name: "Cameroon", Code: "CM"},
		{Name: "Cambodia", Code: "KH"},
	}

	if !equalCountries(countries, expectedCountries) {
		t.Errorf("Expected countries %v, got %v", expectedCountries, countries)
	}
}

func equalCountries(a, b []services.Country) bool {
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

func TestCountryHandler_GetCountries_ShortSearch(t *testing.T) {
	// Create the handler
	countryHandler := handlers.NewCountryHandler()

	// Create a test request with a short search query
	req, err := http.NewRequest("GET", "/api/countries?search=ca", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler
	http.HandlerFunc(countryHandler.GetCountries).ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body
	var countries []services.Country
	err = json.Unmarshal(rr.Body.Bytes(), &countries)
	if err != nil {
		t.Errorf("Failed to decode response body: %v", err)
	}

	if len(countries) != 0 {
		t.Errorf("Expected 0 countries, got %d", len(countries))
	}
}

func TestCountryHandler_GetCountries_ExternalAPIError(t *testing.T) {
	// Setup a test server to simulate an error response
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}))
	defer testServer.Close()

	// Replace the CountriesAPIURL to point to our test server
	originalCountriesAPIURL := services.CountriesAPIURL
	services.SetCountriesAPIURL(testServer.URL)
	defer services.SetCountriesAPIURL(originalCountriesAPIURL)

	// Create the handler
	countryHandler := handlers.NewCountryHandler()

	// Create a test request with a valid search query
	req, err := http.NewRequest("GET", "/api/countries?search=can", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler
	http.HandlerFunc(countryHandler.GetCountries).ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	// Check the response body
	expectedError := "Error fetching countries\n"
	if rr.Body.String() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, rr.Body.String())
	}
}
