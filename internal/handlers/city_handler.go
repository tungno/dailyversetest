/**
 *  CityHandler is responsible for handling HTTP requests related to city operations,
 *  including retrieving cities based on a specific country. This handler integrates with
 *  the CityService and UserService to provide the necessary functionality.
 *
 *  @struct   CityHandler
 *  @inherits None
 *
 *  @properties
 *  - CityService - A service interface for fetching city data based on specified parameters.
 *  - UserService - A service interface for user-related operations (currently unused but available for future enhancements).
 *
 *  @methods
 *  - NewCityHandler(cs, us)  -  Initializes a new CityHandler with the required services.
 *  - GetCities(w, r)         -  Handles GET requests to fetch cities for a specified country.
 *
 *  @endpoint
 *  - /api/cities
 *    - HTTP Method: GET
 *    - Query Parameter: `country` (required) - The name of the country to filter the cities.
 *
 *  @behaviors
 *  - Returns a 400 Bad Request error if the 'country' parameter is missing.
 *  - Returns a 500 Internal Server Error if an error occurs while fetching cities.
 *  - On success, returns a JSON object with a `data` field containing the list of cities.
 *
 *  @example
 *  ```
 *  GET /api/cities?country=Norway
 *
 *  Response:
 *  {
 *      "data": [
 *          "Oslo",
 *          "Bergen",
 *          "Trondheim"
 *      ]
 *  }
 *  ```
 *
 *  @dependencies
 *  - CityServiceInterface: Provides methods to retrieve city data.
 *  - UserServiceInterface: Placeholder for user-related operations.
 *  - utils.WriteJSON: Utility function to write JSON responses.
 *
 *  @file      city_handler.go
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
	"net/http"

	"proh2052-group6/internal/services"
	"proh2052-group6/pkg/utils"
)

// CityHandler struct handles requests related to city operations.
type CityHandler struct {
	CityService services.CityServiceInterface // Service for managing city-related logic.
	UserService services.UserServiceInterface // Service for managing user-related logic.
}

// NewCityHandler initializes a new CityHandler with the necessary dependencies.
func NewCityHandler(cs services.CityServiceInterface, us services.UserServiceInterface) *CityHandler {
	return &CityHandler{
		CityService: cs,
		UserService: us,
	}
}

// GetCities handles GET requests to retrieve a list of cities based on the provided country parameter.
// Endpoint: /api/cities
// Query Parameter:
//   - country (string): The name of the country to filter cities.
func (ch *CityHandler) GetCities(w http.ResponseWriter, r *http.Request) {
	// Extract the 'country' query parameter from the request URL.
	country := r.URL.Query().Get("country")
	if country == "" {
		// Return 400 Bad Request if 'country' parameter is missing.
		http.Error(w, "Missing country parameter", http.StatusBadRequest)
		return
	}

	// Fetch the list of cities for the given country.
	cities, err := ch.CityService.GetCitiesByCountry(country)
	if err != nil {
		// Return 500 Internal Server Error if fetching cities fails.
		http.Error(w, "Error fetching cities", http.StatusInternalServerError)
		return
	}

	// Wrap the fetched cities in a JSON response with a 'data' field.
	response := map[string]interface{}{
		"data": cities,
	}

	// Write the JSON response.
	utils.WriteJSON(w, response)
}
