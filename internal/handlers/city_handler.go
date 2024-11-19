// internal/handlers/city_handler.go
package handlers

import (
	"net/http"

	"proh2052-group6/internal/services"
	"proh2052-group6/pkg/utils"
)

type CityHandler struct {
	CityService services.CityServiceInterface
	UserService services.UserServiceInterface
}

// NewCityHandler initializes a new CityHandler with injected services
func NewCityHandler(cs services.CityServiceInterface, us services.UserServiceInterface) *CityHandler {
	return &CityHandler{
		CityService: cs,
		UserService: us,
	}
}

func (ch *CityHandler) GetCities(w http.ResponseWriter, r *http.Request) {
	country := r.URL.Query().Get("country")
	if country == "" {
		http.Error(w, "Missing country parameter", http.StatusBadRequest)
		return
	}

	cities, err := ch.CityService.GetCitiesByCountry(country)
	if err != nil {
		http.Error(w, "Error fetching cities", http.StatusInternalServerError)
		return
	}

	// Wrap the cities in a 'data' field
	response := map[string]interface{}{
		"data": cities,
	}
	utils.WriteJSON(w, response)
}
