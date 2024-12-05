/**
 *  NewsService provides business logic for fetching news articles based on user-specific
 *  preferences or general search criteria. It integrates with an external news API and
 *  uses the UserRepository to fetch user details when needed.
 *
 *  @interface NewsServiceInterface
 *  @inherits None
 *
 *  @methods
 *  - FetchNews(ctx, userEmail, mode, country, query) - Fetches news articles from the news API based on the input parameters.
 *
 *  @dependencies
 *  - repositories.UserRepository: Fetches user details to determine local news preferences.
 *  - newsdata.io: External news API for fetching articles.
 *
 *  @example
 *  ```
 *  // Fetch general news
 *  articles, err := newsService.FetchNews(ctx, "", "general", "", "technology")
 *
 *  // Fetch local news based on user profile
 *  articles, err := newsService.FetchNews(ctx, "user@example.com", "local", "", "")
 *  ```
 *
 *  @file      news_service.go
 *  @project   DailyVerse
 *  @framework Go HTTP Client with JSON Integration
 */

package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"proh2052-group6/internal/repositories"
)

// NewsServiceInterface defines the contract for fetching news articles.
type NewsServiceInterface interface {
	// FetchNews retrieves news articles based on user and query parameters.
	FetchNews(ctx context.Context, userEmail, mode, country, query string) ([]map[string]interface{}, error)
}

// NewsService implements the NewsServiceInterface and interacts with the external news API.
type NewsService struct {
	UserRepo                  repositories.UserRepository          // Repository for fetching user data.
	HTTPClient                *http.Client                         // HTTP client for making API requests.
	NewsAPIURL                string                               // Base URL of the news API.
	GetCountryAndLanguageCode func(string) (string, string, error) // Helper function to map country names to codes.
}

// NewNewsService initializes a NewsService instance with default values.
func NewNewsService(userRepo repositories.UserRepository) NewsServiceInterface {
	return &NewsService{
		UserRepo:                  userRepo,
		HTTPClient:                http.DefaultClient,
		NewsAPIURL:                "https://newsdata.io/api/1/news",
		GetCountryAndLanguageCode: GetCountryAndLanguageCode,
	}
}

// Global variable for the news API key, sourced from environment variables.
var newsAPIKey = os.Getenv("NEWS_API_KEY")

// FetchNews fetches news articles based on the input parameters.
// Parameters:
// - ctx: Request context for handling deadlines and cancellations.
// - userEmail: The email of the user requesting news (used for local news preferences).
// - mode: Specifies the type of news (e.g., "local").
// - country: The country for which news is requested.
// - query: Search query for filtering news articles.
func (ns *NewsService) FetchNews(ctx context.Context, userEmail, mode, country, query string) ([]map[string]interface{}, error) {
	var url string

	// Handle "local" mode by fetching the user's country if not provided.
	if mode == "local" && country == "" {
		user, err := ns.UserRepo.GetUserByEmail(ctx, userEmail)
		if err != nil || user == nil {
			return nil, fmt.Errorf("Failed to fetch user profile")
		}

		if user.Country != "" {
			country = user.Country
		} else {
			return nil, fmt.Errorf("Country not found in user profile")
		}
	}

	// Construct the API URL for local or general news.
	if mode == "local" && country != "" {
		countryCode, languageCode, err := ns.GetCountryAndLanguageCode(country)
		if err != nil {
			return nil, fmt.Errorf("Invalid country for local news: %v", err)
		}
		url = fmt.Sprintf("%s?country=%s&language=%s&apikey=%s", ns.NewsAPIURL, countryCode, languageCode, newsAPIKey)
	} else {
		url = fmt.Sprintf("%s?language=en&apikey=%s", ns.NewsAPIURL, newsAPIKey)
	}

	// Append query parameter if a search term is provided.
	if query != "" {
		url += fmt.Sprintf("&q=%s", query)
	}

	// Send the HTTP GET request to the news API.
	resp, err := ns.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch news")
	}
	defer resp.Body.Close()

	// Parse the JSON response from the news API.
	var result struct {
		Status       string                   `json:"status"`
		TotalResults int                      `json:"totalResults"`
		Results      []map[string]interface{} `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("Failed to parse news data")
	}

	return result.Results, nil
}
