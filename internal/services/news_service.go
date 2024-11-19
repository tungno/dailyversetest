// internal/services/news_service.go
package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"proh2052-group6/internal/repositories"
)

type NewsServiceInterface interface {
	FetchNews(ctx context.Context, userEmail, mode, country, query string) ([]map[string]interface{}, error)
}

type NewsService struct {
	UserRepo                  repositories.UserRepository
	HTTPClient                *http.Client
	NewsAPIURL                string
	GetCountryAndLanguageCode func(string) (string, string, error)
}

func NewNewsService(userRepo repositories.UserRepository) NewsServiceInterface {
	return &NewsService{
		UserRepo:                  userRepo,
		HTTPClient:                http.DefaultClient,
		NewsAPIURL:                "https://newsdata.io/api/1/news",
		GetCountryAndLanguageCode: GetCountryAndLanguageCode,
	}
}

var newsAPIKey = os.Getenv("NEWS_API_KEY")

func (ns *NewsService) FetchNews(ctx context.Context, userEmail, mode, country, query string) ([]map[string]interface{}, error) {
	var url string

	if mode == "local" && country == "" {
		// Fetch country from user profile
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

	if mode == "local" && country != "" {
		countryCode, languageCode, err := ns.GetCountryAndLanguageCode(country)
		if err != nil {
			return nil, fmt.Errorf("Invalid country for local news: %v", err)
		}
		url = fmt.Sprintf("%s?country=%s&language=%s&apikey=%s", ns.NewsAPIURL, countryCode, languageCode, newsAPIKey)
	} else {
		url = fmt.Sprintf("%s?language=en&apikey=%s", ns.NewsAPIURL, newsAPIKey)
	}

	if query != "" {
		url += fmt.Sprintf("&q=%s", query)
	}

	resp, err := ns.HTTPClient.Get(url)
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
