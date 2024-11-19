// internal/handlers/country_handler.go
package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"proh2052-group6/internal/services"
)

type CountryHandler struct {
}

func NewCountryHandler() *CountryHandler {
	return &CountryHandler{}
}

func (ch *CountryHandler) GetCountries(w http.ResponseWriter, r *http.Request) {
	searchQuery := strings.ToLower(r.URL.Query().Get("search"))

	if len(searchQuery) < 3 {
		json.NewEncoder(w).Encode([]services.Country{})
		return
	}

	countries, err := services.GetCountries(searchQuery)
	if err != nil {
		http.Error(w, "Error fetching countries", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(countries)
}
