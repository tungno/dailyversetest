/**
 *  CityService provides functionality to fetch cities based on the specified country.
 *  It communicates with an external API to retrieve city data.
 *
 *  @file       city_service.go
 *  @package    services
 *
 *  @interfaces
 *  - CityServiceInterface - Defines the contract for city-related operations.
 *
 *  @methods
 *  - NewCityService()                                - Initializes a new instance of CityService.
 *  - GetCitiesByCountry(country) ([]string, error)   - Fetches a list of cities for the specified country.
 *
 *  @dependencies
 *  - config.CitiesAPIURL: Configuration value containing the external API endpoint.
 *  - http.Client: Used for making HTTP requests.
 *
 *  @behaviors
 *  - Sends a POST request to the external API with the country name as the request payload.
 *  - Parses the JSON response and returns the list of cities on success.
 *  - Handles errors gracefully, including API errors, decoding errors, and connection issues.
 *
 *  @example
 *  ```
 *  cityService := NewCityService()
 *  cities, err := cityService.GetCitiesByCountry("Norway")
 *  if err != nil {
 *      log.Fatal("Failed to fetch cities:", err)
 *  }
 *  fmt.Println("Cities in Norway:", cities)
 *  ```
 *
 *  @errors
 *  - Returns an error if the request body cannot be created.
 *  - Returns an error if the HTTP request fails.
 *  - Returns an error if the API response indicates a failure.
 *  - Returns an error if the JSON response cannot be decoded.
 *
 *  @authors
 *      - Aayush
 *      - Tung
 *      - Boss
 *      - Majd
 */

package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"proh2052-group6/internal/config"
)

// CityServiceInterface defines the methods for CityService.
type CityServiceInterface interface {
	// GetCitiesByCountry fetches cities for a given country.
	GetCitiesByCountry(country string) ([]string, error)
}

// CityService implements CityServiceInterface.
type CityService struct {
	HTTPClient   *http.Client // HTTP client for making API requests.
	CitiesAPIURL string       // URL of the external cities API.
}

// NewCityService initializes a new CityService.
func NewCityService() CityServiceInterface {
	return &CityService{
		HTTPClient:   http.DefaultClient,
		CitiesAPIURL: config.CitiesAPIURL,
	}
}

// GetCitiesByCountry fetches cities for a given country by calling an external API.
func (cs *CityService) GetCitiesByCountry(country string) ([]string, error) {
	// Create the request body for the external API.
	requestBody, err := json.Marshal(map[string]string{"country": country})
	if err != nil {
		return nil, fmt.Errorf("failed to create request body: %v", err)
	}

	// Make a POST request to the external API.
	resp, err := cs.HTTPClient.Post(cs.CitiesAPIURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("error fetching cities: %v", err)
	}
	defer resp.Body.Close()

	// Read and parse the response body.
	var cityResponse struct {
		Error bool     `json:"error"` // Indicates if there was an error in the API response.
		Msg   string   `json:"msg"`   // Error message or additional information from the API.
		Data  []string `json:"data"`  // List of cities returned by the API.
	}

	if err := json.NewDecoder(resp.Body).Decode(&cityResponse); err != nil {
		return nil, fmt.Errorf("error decoding cities response: %v", err)
	}

	// Check if the API response contains an error.
	if cityResponse.Error {
		return nil, fmt.Errorf("error fetching cities: %s", cityResponse.Msg)
	}

	// Return the list of cities on success.
	return cityResponse.Data, nil
}
