// tests/mocks/mock_city_service.go
package mocks

import (
	"fmt"
)

// MockCityService is a mock implementation of the CityServiceInterface.
// It allows you to define custom behavior for the GetCitiesByCountry method.
type MockCityService struct {
	GetCitiesByCountryFunc func(country string) ([]string, error)
}

// GetCitiesByCountry calls the mocked GetCitiesByCountryFunc if it's set.
// Otherwise, it returns nil or a default value.
func (m *MockCityService) GetCitiesByCountry(country string) ([]string, error) {
	if m.GetCitiesByCountryFunc != nil {
		return m.GetCitiesByCountryFunc(country)
	}
	return nil, fmt.Errorf("GetCitiesByCountryFunc not implemented")
}
