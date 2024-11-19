// internal/handlers/city_handler.go
package handlers

import (
	"encoding/json"
	"net/http"

	"proh2052-group6/internal/services"
)

type CityHandler struct{}

func NewCityHandler() *CityHandler {
	return &CityHandler{}
}

func (ch *CityHandler) GetCities(w http.ResponseWriter, r *http.Request) {
	country := r.URL.Query().Get("country")
	if country == "" {
		http.Error(w, "Missing country parameter", http.StatusBadRequest)
		return
	}

	cities, err := services.GetCitiesByCountry(country)
	if err != nil {
		http.Error(w, "Error fetching cities", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(cities)
}
