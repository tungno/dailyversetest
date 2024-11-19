// internal/services/country_service.go
package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Country struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

func GetCountries(searchQuery string) ([]Country, error) {
	resp, err := http.Get("https://restcountries.com/v3.1/all")
	if err != nil {
		return nil, fmt.Errorf("Error fetching countries: %v", err)
	}
	defer resp.Body.Close()

	var countriesData []struct {
		Name struct {
			Common string `json:"common"`
		} `json:"name"`
		Cca2 string `json:"cca2"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&countriesData); err != nil {
		return nil, fmt.Errorf("Error decoding response: %v", err)
	}

	var countries []Country
	for _, country := range countriesData {
		countryName := strings.ToLower(country.Name.Common)
		if strings.HasPrefix(countryName, searchQuery) {
			countries = append(countries, Country{
				Name: country.Name.Common,
				Code: country.Cca2,
			})
		}
	}

	return countries, nil
}
