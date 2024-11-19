// internal/handlers/news_handler.go
package handlers

import (
	"net/http"

	"proh2052-group6/internal/services"
	"proh2052-group6/pkg/utils"
)

type NewsHandler struct {
	NewsService services.NewsServiceInterface
}

func NewNewsHandler(ns services.NewsServiceInterface) *NewsHandler {
	return &NewsHandler{NewsService: ns}
}

func (nh *NewsHandler) FetchNews(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("mode")
	country := r.URL.Query().Get("country")
	query := r.URL.Query().Get("q")

	userEmail := r.Context().Value("userEmail").(string)

	news, err := nh.NewsService.FetchNews(r.Context(), userEmail, mode, country, query)
	if err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, news)
}
