// tests/handlers_test.go
package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"proh2052-group6/internal/handlers"
	"proh2052-group6/internal/services"
	"proh2052-group6/pkg/models"
	"proh2052-group6/tests/mocks"
)

func TestUserSignupHandler(t *testing.T) {
	mockDB := &mocks.MockDB{}
	mockEmail := &mocks.MockEmailService{}
	userService := services.NewUserService(mockDB, mockEmail)
	userHandler := handlers.NewUserHandler(userService)

	user := models.User{
		Email:    "test@example.com",
		Password: "Password123!",
		Country:  "Norway",
		City:     "Oslo",
		Username: "testuser",
	}

	body, _ := json.Marshal(user)
	req, err := http.NewRequest("POST", "/api/signup", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(userHandler.Signup)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

// Add more tests...
