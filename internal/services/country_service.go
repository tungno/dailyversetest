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

var (
	countryHTTPClient = http.DefaultClient
	CountriesAPIURL   = "https://restcountries.com/v3.1/all"
)

func SetCountryHTTPClient(client *http.Client) {
	countryHTTPClient = client
}

func SetCountriesAPIURL(url string) {
	CountriesAPIURL = url
}

func GetCountries(searchQuery string) ([]Country, error) {
	resp, err := countryHTTPClient.Get(CountriesAPIURL)
	if err != nil {
		return nil, fmt.Errorf("Error fetching countries: %v", err)
	}
	defer resp.Body.Close()

	var countriesData []struct {
		Name struct {
			Common string `json:"common"`
		} `json:"name"`
		CCA2 string `json:"cca2"`
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
				Code: country.CCA2,
			})
		}
	}

	return countries, nil
}
