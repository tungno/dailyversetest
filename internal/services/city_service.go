// internal/services/city_service.go
package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	cityHTTPClient = http.DefaultClient
	CitiesAPIURL   = "https://countriesnow.space/api/v0.1/countries/cities"
)

func SetCityHTTPClient(client *http.Client) {
	cityHTTPClient = client
}

func SetCitiesAPIURL(url string) {
	CitiesAPIURL = url
}

func GetCitiesByCountry(country string) ([]string, error) {
	// Create the request body for the external API
	requestBody, err := json.Marshal(map[string]string{"country": country})
	if err != nil {
		return nil, fmt.Errorf("Failed to create request body: %v", err)
	}

	// Make a POST request to the external API
	resp, err := cityHTTPClient.Post(CitiesAPIURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("Error fetching cities: %v", err)
	}
	defer resp.Body.Close()

	// Read and parse the response body
	var cityResponse struct {
		Error bool     `json:"error"`
		Msg   string   `json:"msg"`
		Data  []string `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&cityResponse); err != nil {
		return nil, fmt.Errorf("Error decoding cities response: %v", err)
	}

	if cityResponse.Error {
		return nil, fmt.Errorf("Error fetching cities: %s", cityResponse.Msg)
	}

	return cityResponse.Data, nil
}
