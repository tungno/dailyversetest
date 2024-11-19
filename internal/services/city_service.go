// internal/services/city_service.go
package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"proh2052-group6/internal/config"
)

// CityServiceInterface defines the methods for CityService
type CityServiceInterface interface {
	GetCitiesByCountry(country string) ([]string, error)
}

// CityService implements CityServiceInterface
type CityService struct {
	HTTPClient   *http.Client
	CitiesAPIURL string
}

// NewCityService initializes a new CityService
func NewCityService() CityServiceInterface {
	return &CityService{
		HTTPClient:   http.DefaultClient,
		CitiesAPIURL: config.CitiesAPIURL,
	}
}

// GetCitiesByCountry fetches cities for a given country
func (cs *CityService) GetCitiesByCountry(country string) ([]string, error) {
	// Create the request body for the external API
	requestBody, err := json.Marshal(map[string]string{"country": country})
	if err != nil {
		return nil, fmt.Errorf("failed to create request body: %v", err)
	}

	// Make a POST request to the external API
	resp, err := cs.HTTPClient.Post(cs.CitiesAPIURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("error fetching cities: %v", err)
	}
	defer resp.Body.Close()

	// Read and parse the response body
	var cityResponse struct {
		Error bool     `json:"error"`
		Msg   string   `json:"msg"`
		Data  []string `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&cityResponse); err != nil {
		return nil, fmt.Errorf("error decoding cities response: %v", err)
	}

	if cityResponse.Error {
		return nil, fmt.Errorf("error fetching cities: %s", cityResponse.Msg)
	}

	return cityResponse.Data, nil
}
