/**
 *  MockCityService provides a mock implementation of the CityServiceInterface for testing purposes.
 *  This mock allows you to define custom behavior for the `GetCitiesByCountry` method, enabling
 *  controlled testing of components that depend on CityService without using the actual implementation.
 *
 *  @struct   MockCityService
 *  @inherits CityServiceInterface
 *
 *  @fields
 *  - GetCitiesByCountryFunc (func): A customizable function that simulates the behavior of
 *    `GetCitiesByCountry` for specific test cases.
 *
 *  @methods
 *  - GetCitiesByCountry(country) ([]string, error): Calls the mock function to simulate fetching cities
 *    by country. If the mock function is not defined, it returns a default error.
 *
 *  @example
 *  ```
 *  // Define mock behavior
 *  mockCityService := &MockCityService{
 *      GetCitiesByCountryFunc: func(country string) ([]string, error) {
 *          if country == "TestCountry" {
 *              return []string{"City1", "City2"}, nil
 *          }
 *          return nil, fmt.Errorf("Country not found")
 *      },
 *  }
 *
 *  // Call the mocked method
 *  cities, err := mockCityService.GetCitiesByCountry("TestCountry")
 *  fmt.Println(cities) // Output: [City1 City2]
 *  ```
 *
 *  @file      mock_city_service.go
 *  @project   DailyVerse
 *  @framework Go Testing with Mock Services
 */

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
// Otherwise, it returns nil or a default error.
func (m *MockCityService) GetCitiesByCountry(country string) ([]string, error) {
	if m.GetCitiesByCountryFunc != nil {
		return m.GetCitiesByCountryFunc(country)
	}
	return nil, fmt.Errorf("GetCitiesByCountryFunc not implemented")
}
