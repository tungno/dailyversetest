// tests/services_test.go
package tests

import (
	"context"
	"testing"

	"proh2052-group6/internal/services"
	"proh2052-group6/pkg/models"
	"proh2052-group6/tests/mocks"
)

func TestUserSignup(t *testing.T) {
	mockDB := &mocks.MockDB{}
	mockEmail := &mocks.MockEmailService{}

	userService := services.NewUserService(mockDB, mockEmail)

	user := &models.User{
		Email:    "test@example.com",
		Password: "Password123!",
		Country:  "Norway",
		City:     "Oslo",
		Username: "testuser",
	}

	err := userService.Signup(context.Background(), user)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(mockEmail.SentEmails) != 1 {
		t.Errorf("Expected 1 email sent, got %d", len(mockEmail.SentEmails))
	}
}

// Add more tests...
