// internal/services/news_service.go
package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type NewsServiceInterface interface {
	FetchNews(ctx context.Context, userEmail, mode, country, query string) ([]map[string]interface{}, error)
}

type NewsService struct {
	DB DatabaseInterface
}

func NewNewsService() NewsServiceInterface {
	return &NewsService{}
}

var newsAPIKey = os.Getenv("NEWS_API_KEY")

func (ns *NewsService) FetchNews(ctx context.Context, userEmail, mode, country, query string) ([]map[string]interface{}, error) {
	if mode == "local" && country == "" {
		// Fetch country from user profile
		doc, err := ns.DB.Collection("users").Doc(userEmail).Get(ctx)
		if err != nil {
			return nil, fmt.Errorf("Failed to fetch user profile")
		}

		if profileCountry, ok := doc.Data()["Country"].(string); ok && profileCountry != "" {
			country = profileCountry
		} else {
			return nil, fmt.Errorf("Country not found in user profile")
		}
	}

	var url string
	if mode == "local" && country != "" {
		countryCode, languageCode, err := getCountryAndLanguageCode(country)
		if err != nil {
			return nil, fmt.Errorf("Invalid country for local news")
		}
		url = fmt.Sprintf("https://newsdata.io/api/1/news?country=%s&language=%s&apikey=%s", countryCode, languageCode, newsAPIKey)
	} else {
		url = fmt.Sprintf("https://newsdata.io/api/1/news?language=en&apikey=%s", newsAPIKey)
	}

	if query != "" {
		url += fmt.Sprintf("&q=%s", query)
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch news")
	}
	defer resp.Body.Close()

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

// Helper functions...

func getCountryAndLanguageCode(countryName string) (string, string, error) {
	// Implement mapping
	return "", "", nil
}
