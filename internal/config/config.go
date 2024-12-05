/**
 *  Configuration file for defining global constants and external API URLs
 *  used in the DailyVerse application. This file centralizes configuration
 *  settings for better maintainability and scalability.
 *
 *  @file      config.go
 *  @project   DailyVerse
 *  @framework Go HTTP Server
 *  @purpose   Provides global configuration variables for external API integrations.
 *  @authors
 *      - Aayush
 *      - Tung
 *      - Boss
 *      - Majd
 */

package config

var (
	// CountriesAPIURL defines the endpoint for retrieving country data.
	CountriesAPIURL = "https://restcountries.com/v3.1/all"

	// CitiesAPIURL defines the endpoint for retrieving cities based on countries.
	CitiesAPIURL = "https://countriesnow.space/api/v0.1/countries/cities"
)
