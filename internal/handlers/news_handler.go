/**
 *  NewsHandler handles HTTP requests for fetching news articles. It integrates with the
 *  NewsService to retrieve news based on various filters like mode, country, and search query.
 *
 *  @struct   NewsHandler
 *  @inherits None
 *
 *  @methods
 *  - NewNewsHandler(ns)         - Initializes a new NewsHandler with the required NewsService.
 *  - FetchNews(w, r)            - Handles GET requests to fetch news articles based on filters.
 *
 *  @endpoint
 *  - /api/news
 *    - HTTP Method: GET
 *    - Query Parameters:
 *      - mode (string, optional): Filter for news type or category.
 *      - country (string, optional): Filter for news by country.
 *      - q (string, optional): Search query for filtering news articles.
 *
 *  @behaviors
 *  - Retrieves news articles using filters provided as query parameters.
 *  - Returns a 500 Internal Server Error for service-layer failures.
 *  - On success, responds with a JSON array of news articles.
 *
 *  @example
 *  ```
 *  GET /api/news?mode=technology&country=US&q=AI
 *
 *  Response:
 *  [
 *      {
 *          "title": "Advances in AI",
 *          "source": "TechDaily",
 *          "url": "https://example.com/ai-news"
 *      },
 *      {
 *          "title": "AI in 2024",
 *          "source": "FutureTrends",
 *          "url": "https://example.com/ai-2024"
 *      }
 *  ]
 *  ```
 *
 *  @dependencies
 *  - NewsServiceInterface: Provides the logic for fetching news articles.
 *  - utils.WriteJSON, utils.WriteJSONError: Utility functions for JSON responses.
 *
 *  @file      news_handler.go
 *  @project   DailyVerse
 *  @framework Go HTTP Server
 *  @authors
 *      - Aayush
 *      - Tung
 *      - Boss
 *      - Majd
 */

package handlers

import (
	"net/http"

	"proh2052-group6/internal/services"
	"proh2052-group6/pkg/utils"
)

// NewsHandler manages HTTP requests for fetching news articles.
type NewsHandler struct {
	NewsService services.NewsServiceInterface // Service for news-related operations.
}

// NewNewsHandler initializes a NewsHandler with the given NewsService.
func NewNewsHandler(ns services.NewsServiceInterface) *NewsHandler {
	return &NewsHandler{NewsService: ns}
}

// FetchNews handles GET requests to fetch news articles based on query parameters.
// Query Parameters:
//   - mode (string, optional): Filter for news type or category.
//   - country (string, optional): Filter for news by country.
//   - q (string, optional): Search query for filtering news articles.
func (nh *NewsHandler) FetchNews(w http.ResponseWriter, r *http.Request) {
	// Extract query parameters.
	mode := r.URL.Query().Get("mode")
	country := r.URL.Query().Get("country")
	query := r.URL.Query().Get("q")

	// Retrieve user email from the request context.
	userEmail := r.Context().Value("userEmail").(string)

	// Fetch news articles using the NewsService.
	news, err := nh.NewsService.FetchNews(r.Context(), userEmail, mode, country, query)
	if err != nil {
		// Return a 500 Internal Server Error if the news fetching fails.
		utils.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the fetched news as a JSON response.
	utils.WriteJSON(w, news)
}
