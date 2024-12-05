/**
 *  UserHandler manages HTTP requests related to user authentication and account operations,
 *  including signup, login, email verification, password reset, and user search functionality.
 *  It integrates with the UserService to handle business logic and returns appropriate HTTP responses.
 *
 *  @struct   UserHandler
 *  @inherits None
 *
 *  @methods
 *  - NewUserHandler(us)                  - Initializes a new UserHandler with the required UserService.
 *  - Signup(w, r)                        - Handles user signup requests.
 *  - Login(w, r)                         - Handles user login requests.
 *  - ResendOTP(w, r)                     - Resends an OTP for email verification.
 *  - VerifyEmail(w, r)                   - Verifies a user's email with an OTP.
 *  - ForgotPassword(w, r)                - Initiates a password reset by sending an OTP to the user's email.
 *  - ResetPassword(w, r)                 - Resets the user's password using an OTP.
 *  - GetUserInfo(w, r)                   - Fetches the authenticated user's information.
 *  - SearchUsersByUsername(w, r)         - Searches for users by username.
 *
 *  @endpoint
 *  - /api/signup                         - POST request to register a new user.
 *  - /api/login                          - POST request to log in an existing user.
 *  - /api/resend-otp                     - POST request to resend an OTP for email verification.
 *  - /api/verify-email                   - POST request to verify a user's email with an OTP.
 *  - /api/forgot-password                - POST request to initiate a password reset.
 *  - /api/reset-password                 - POST request to reset a user's password.
 *  - /api/me                             - GET request to fetch the authenticated user's information.
 *  - /api/users/search                   - GET request to search for users by username.
 *
 *  @behaviors
 *  - Validates incoming request data and handles errors appropriately.
 *  - Communicates with the UserService to perform user-related operations.
 *  - Returns JSON responses with appropriate HTTP status codes.
 *
 *  @example
 *  ```
 *  POST /api/signup
 *  Body: {
 *      "email": "user@example.com",
 *      "password": "securePassword123",
 *      "username": "user123"
 *  }
 *
 *  Response: {
 *      "message": "Signup successful. Please verify your email."
 *  }
 *  ```
 *
 *  @dependencies
 *  - UserServiceInterface: Provides business logic for user operations.
 *  - utils.WriteJSON, utils.WriteJSONError: Utility functions for JSON responses.
 *
 *  @file      user_handler.go
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
	"encoding/json"
	"net/http"

	"proh2052-group6/internal/services"
	"proh2052-group6/pkg/models"
	"proh2052-group6/pkg/utils"
)

// UserHandler handles user-related HTTP requests.
type UserHandler struct {
	UserService services.UserServiceInterface // Service for user-related business logic.
}

// NewUserHandler initializes a UserHandler with the given UserService.
func NewUserHandler(us services.UserServiceInterface) *UserHandler {
	return &UserHandler{UserService: us}
}

// Signup handles POST requests for user registration.
func (uh *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.WriteJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := uh.UserService.Signup(r.Context(), &user); err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, map[string]string{"message": "Signup successful. Please verify your email."})
}

// Login handles POST requests for user login.
func (uh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginData models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		utils.WriteJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := uh.UserService.Login(r.Context(), &loginData)
	if err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	utils.WriteJSON(w, map[string]string{"token": token})
}

// ResendOTP handles POST requests to resend an OTP for email verification.
func (uh *UserHandler) ResendOTP(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		utils.WriteJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := uh.UserService.ResendOTP(r.Context(), requestData.Email); err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, map[string]string{"message": "A new OTP has been sent to your email address."})
}

// VerifyEmail handles POST requests to verify a user's email using an OTP.
func (uh *UserHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		utils.WriteJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := uh.UserService.VerifyEmail(r.Context(), requestData.Email, requestData.OTP)
	if err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.WriteJSON(w, map[string]string{"message": "Email verified successfully", "token": token})
}

// ForgotPassword handles POST requests to initiate a password reset.
func (uh *UserHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		utils.WriteJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := uh.UserService.ForgotPassword(r.Context(), requestData.Email); err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, map[string]string{"message": "If the email exists, an OTP has been sent."})
}

// ResetPassword handles POST requests to reset a user's password using an OTP.
func (uh *UserHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Email       string `json:"email"`
		OTP         string `json:"otp"`
		NewPassword string `json:"newPassword"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		utils.WriteJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := uh.UserService.ResetPassword(r.Context(), requestData.Email, requestData.OTP, requestData.NewPassword); err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.WriteJSON(w, map[string]string{"message": "Password has been reset successfully."})
}

// GetUserInfo handles GET requests to fetch the authenticated user's information.
func (uh *UserHandler) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	userEmail := r.Context().Value("userEmail").(string)

	userInfo, err := uh.UserService.GetUserInfo(r.Context(), userEmail)
	if err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	utils.WriteJSON(w, userInfo)
}

// SearchUsersByUsername handles GET requests to search for users by username.
func (uh *UserHandler) SearchUsersByUsername(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		utils.WriteJSONError(w, "Missing search query", http.StatusBadRequest)
		return
	}

	userEmail := r.Context().Value("userEmail").(string)

	results, err := uh.UserService.SearchUsersByUsername(r.Context(), userEmail, query)
	if err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusNotFound)
		return
	}

	utils.WriteJSON(w, results)
}
