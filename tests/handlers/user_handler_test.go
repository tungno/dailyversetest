// tests/handlers/user_handler_test.go
package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"proh2052-group6/internal/handlers"
	"proh2052-group6/internal/services"
	"proh2052-group6/pkg/models"
	"proh2052-group6/pkg/utils"
	"proh2052-group6/tests/mocks"
)

func TestUserHandler_Signup(t *testing.T) {
	// Create mocks
	mockUserRepo := mocks.NewMockUserRepository(make(map[string]*models.User))
	mockEmailService := &mocks.MockEmailService{}
	userService := services.NewUserService(mockUserRepo, mockEmailService)
	userHandler := handlers.NewUserHandler(userService)

	// Create a test HTTP request
	user := models.User{
		Email:    "test@example.com",
		Username: "testuser",
		Password: "Password123!",
		Country:  "TestCountry",
		City:     "TestCity",
	}
	requestBody, _ := json.Marshal(user)
	req, err := http.NewRequest("POST", "/api/signup", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler
	handler := http.HandlerFunc(userHandler.Signup)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body
	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	expectedMessage := "Signup successful. Please verify your email."
	if response["message"] != expectedMessage {
		t.Errorf("Expected message '%s', got '%s'", expectedMessage, response["message"])
	}

	// Check that the user was created in the mock repository
	savedUser, err := mockUserRepo.GetUserByEmail(context.Background(), user.Email)
	if err != nil {
		t.Errorf("User was not saved in repository")
	} else {
		// Use 'savedUser' in assertions
		if savedUser.Email != user.Email {
			t.Errorf("Expected saved user email to be '%s', got '%s'", user.Email, savedUser.Email)
		}
		if savedUser.Username != user.Username {
			t.Errorf("Expected saved user username to be '%s', got '%s'", user.Username, savedUser.Username)
		}
	}

	// Check that an email was sent
	if len(mockEmailService.SentEmails) != 1 {
		t.Errorf("Expected 1 email to be sent, got %d", len(mockEmailService.SentEmails))
	}
}

func TestUserHandler_Login(t *testing.T) {
	// Create mocks
	mockUserRepo := mocks.NewMockUserRepository(make(map[string]*models.User))
	mockEmailService := &mocks.MockEmailService{}
	userService := services.NewUserService(mockUserRepo, mockEmailService)
	userHandler := handlers.NewUserHandler(userService)

	// Add a user to the mock repository
	user := &models.User{
		Email:      "test@example.com",
		Username:   "testuser",
		Password:   utils.HashPassword("Password123!"),
		Country:    "TestCountry",
		City:       "TestCity",
		IsVerified: true,
	}
	mockUserRepo.CreateUser(context.Background(), user)

	// Create a test HTTP request
	loginData := models.LoginRequest{
		Email:    "test@example.com",
		Password: "Password123!",
	}
	requestBody, _ := json.Marshal(loginData)
	req, err := http.NewRequest("POST", "/api/login", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler
	handler := http.HandlerFunc(userHandler.Login)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body
	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	if response["token"] == "" {
		t.Errorf("Expected a token in response")
	}
}

func TestUserHandler_ResendOTP(t *testing.T) {
	// Create mocks
	mockUserRepo := mocks.NewMockUserRepository(make(map[string]*models.User))
	mockEmailService := &mocks.MockEmailService{}
	userService := services.NewUserService(mockUserRepo, mockEmailService)
	userHandler := handlers.NewUserHandler(userService)

	// Add an unverified user to the mock repository
	user := &models.User{
		Email:      "test@example.com",
		Username:   "testuser",
		Password:   utils.HashPassword("Password123!"),
		Country:    "TestCountry",
		City:       "TestCity",
		IsVerified: false,
	}
	mockUserRepo.CreateUser(context.Background(), user)

	// Create a test HTTP request
	requestData := map[string]string{
		"email": "test@example.com",
	}
	requestBody, _ := json.Marshal(requestData)
	req, err := http.NewRequest("POST", "/api/resend-otp", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler
	handler := http.HandlerFunc(userHandler.ResendOTP)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body
	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	expectedMessage := "A new OTP has been sent to your email address."
	if response["message"] != expectedMessage {
		t.Errorf("Expected message '%s', got '%s'", expectedMessage, response["message"])
	}

	// Check that an email was sent
	if len(mockEmailService.SentEmails) != 1 {
		t.Errorf("Expected 1 email to be sent, got %d", len(mockEmailService.SentEmails))
	}
}

func TestUserHandler_VerifyEmail(t *testing.T) {
	// Create mocks
	mockUserRepo := mocks.NewMockUserRepository(make(map[string]*models.User))
	mockEmailService := &mocks.MockEmailService{}
	userService := services.NewUserService(mockUserRepo, mockEmailService)
	userHandler := handlers.NewUserHandler(userService)

	// Add an unverified user with an OTP
	user := &models.User{
		Email:        "test@example.com",
		Username:     "testuser",
		Password:     utils.HashPassword("Password123!"),
		Country:      "TestCountry",
		City:         "TestCity",
		IsVerified:   false,
		OTP:          "123456",
		OTPExpiresAt: time.Now().Add(5 * time.Minute),
	}
	mockUserRepo.CreateUser(context.Background(), user)

	// Create a test HTTP request
	requestData := map[string]string{
		"email": "test@example.com",
		"otp":   "123456",
	}
	requestBody, _ := json.Marshal(requestData)
	req, err := http.NewRequest("POST", "/api/verify-email", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler
	handler := http.HandlerFunc(userHandler.VerifyEmail)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body
	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	expectedMessage := "Email verified successfully"
	if response["message"] != expectedMessage {
		t.Errorf("Expected message '%s', got '%s'", expectedMessage, response["message"])
	}

	if response["token"] == "" {
		t.Errorf("Expected a token in response")
	}
}

func TestUserHandler_GetUserInfo(t *testing.T) {
	// Create mocks
	mockUserRepo := mocks.NewMockUserRepository(make(map[string]*models.User))
	mockEmailService := &mocks.MockEmailService{}
	userService := services.NewUserService(mockUserRepo, mockEmailService)
	userHandler := handlers.NewUserHandler(userService)

	// Add a verified user to the mock repository
	user := &models.User{
		Email:      "test@example.com",
		Username:   "testuser",
		Password:   utils.HashPassword("Password123!"),
		Country:    "TestCountry",
		City:       "TestCity",
		IsVerified: true,
	}
	mockUserRepo.CreateUser(context.Background(), user)

	// Create a test HTTP request
	req, err := http.NewRequest("GET", "/api/user-info", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add the userEmail to the request context
	ctx := context.WithValue(req.Context(), "userEmail", user.Email)
	req = req.WithContext(ctx)

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler
	handler := http.HandlerFunc(userHandler.GetUserInfo)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body
	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	if response["email"] != user.Email {
		t.Errorf("Expected email '%s', got '%s'", user.Email, response["email"])
	}
	if response["username"] != user.Username {
		t.Errorf("Expected username '%s', got '%s'", user.Username, response["username"])
	}
}
