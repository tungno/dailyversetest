/**
 *  CountryService provides functionality to retrieve country data from an external API
 *  and filter the results based on a search query. This service integrates with a RESTful
 *  countries API to fetch information about all countries.
 *
 *  @struct   Country
 *  @inherits None
 *
 *  @methods
 *  - SetCountryHTTPClient(client)       - Sets a custom HTTP client for API requests (useful for testing).
 *  - SetCountriesAPIURL(url)            - Sets the API endpoint for fetching country data.
 *  - GetCountries(searchQuery)          - Fetches and filters country data based on the search query.
 *
 *  @behaviors
 *  - Retrieves data from the countries API, defined in `config.CountriesAPIURL`.
 *  - Filters countries by name, matching the search query with a case-insensitive prefix.
 *  - Ensures graceful handling of errors during API calls or JSON decoding.
 *
 *  @dependencies
 *  - config.CountriesAPIURL: Configuration variable for the countries API endpoint.
 *  - http.Client: HTTP client used for API requests.
 *  - json: Used for decoding JSON responses from the API.
 *
 *  @example
 *  ```
 *  // Fetch countries starting with "nor"
 *  SetCountriesAPIURL("https://restcountries.com/v3.1/all")
 *  countries, err := GetCountries("nor")
 *
 *  Response:
 *  [
 *      { "name": "Norway", "code": "NO" },
 *      { "name": "Northern Ireland", "code": "GB" }
 *  ]
 *  ```
 *
 *  @file      country_service.go
 *  @project   DailyVerse
 *  @framework Go HTTP Client & REST API Integration
 */

package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"proh2052-group6/internal/config"
	"strings"
)

// Country represents a country entity with its name and code.
type Country struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

var (
	countryHTTPClient = http.DefaultClient // Default HTTP client for making API calls.
)

// SetCountryHTTPClient allows setting a custom HTTP client for testing or customization.
func SetCountryHTTPClient(client *http.Client) {
	countryHTTPClient = client
}

// SetCountriesAPIURL sets the API endpoint for fetching country data.
func SetCountriesAPIURL(url string) {
	config.CountriesAPIURL = url
}

// GetCountries fetches and filters country data based on a search query.
// Returns a list of countries whose names start with the given query.
func GetCountries(searchQuery string) ([]Country, error) {
	// Fetch data from the countries API.
	resp, err := countryHTTPClient.Get(config.CountriesAPIURL)
	if err != nil {
		return nil, fmt.Errorf("Error fetching countries: %v", err)
	}
	defer resp.Body.Close()

	// Decode the API response into a temporary structure.
	var countriesData []struct {
		Name struct {
			Common string `json:"common"`
		} `json:"name"`
		CCA2 string `json:"cca2"` // Country code.
	}

	if err := json.NewDecoder(resp.Body).Decode(&countriesData); err != nil {
		return nil, fmt.Errorf("Error decoding response: %v", err)
	}

	// Filter countries by the search query (case-insensitive prefix match).
	var countries []Country
	for _, country := range countriesData {
		countryName := strings.ToLower(country.Name.Common)
		if strings.HasPrefix(countryName, searchQuery) {
			countries = append(countries, Country{
				Name: country.Name.Common,
				Code: country.CCA2,
			})
		}
	}

	return countries, nil
}
