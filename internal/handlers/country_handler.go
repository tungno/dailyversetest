/**
 *  CountryHandler handles HTTP requests related to retrieving country information.
 *  This handler supports filtering countries by a search query, ensuring only relevant
 *  results are returned to the client.
 *
 *  @struct   CountryHandler
 *  @inherits None
 *
 *  @methods
 *  - NewCountryHandler()         - Initializes a new CountryHandler instance.
 *  - GetCountries(w, r)          - Handles GET requests to fetch a list of countries based on a search query.
 *
 *  @endpoint
 *  - /api/countries
 *    - HTTP Method: GET
 *    - Query Parameter: `search` (optional) - A substring to filter country names (minimum length: 3 characters).
 *
 *  @behaviors
 *  - Returns an empty list if the search query is less than 3 characters.
 *  - Returns a 500 Internal Server Error if there is an issue fetching countries.
 *  - On success, returns a JSON array of countries matching the search query.
 *
 *  @example
 *  ```
 *  GET /api/countries?search=nor
 *
 *  Response:
 *  [
 *      { "name": "Norway" },
 *      { "name": "Northern Ireland" }
 *  ]
 *  ```
 *
 *  @dependencies
 *  - services.GetCountries: Fetches country data filtered by the search query.
 *
 *  @file      country_handler.go
 *  @project   DailyVerse
 *  @framework Go HTTP Server
 *  @authors
 *      - Aayush
 *      - Tung
 *      - Boss
 *      - Majd
 */

package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"proh2052-group6/internal/services"
)

// CountryHandler struct for handling country-related requests.
type CountryHandler struct{}

// NewCountryHandler initializes a new CountryHandler instance.
func NewCountryHandler() *CountryHandler {
	return &CountryHandler{}
}

// GetCountries handles GET requests to fetch a list of countries based on a search query.
// Endpoint: /api/countries
// Query Parameter:
//   - search (string, optional): Substring to filter country names. Minimum length is 3 characters.
func (ch *CountryHandler) GetCountries(w http.ResponseWriter, r *http.Request) {
	// Extract and sanitize the search query from the URL.
	searchQuery := strings.ToLower(r.URL.Query().Get("search"))

	// Return an empty list if the search query is too short.
	if len(searchQuery) < 3 {
		json.NewEncoder(w).Encode([]services.Country{})
		return
	}

	// Fetch the list of countries matching the search query.
	countries, err := services.GetCountries(searchQuery)
	if err != nil {
		// Return a 500 error if there is an issue fetching countries.
		http.Error(w, "Error fetching countries", http.StatusInternalServerError)
		return
	}

	// Encode the result as JSON and write it to the response.
	json.NewEncoder(w).Encode(countries)
}
