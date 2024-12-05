/**
 *  MockUserService provides a mock implementation of the UserServiceInterface, allowing you to define
 *  custom behavior for each user-related operation in your tests. Instead of interacting with a real
 *  database or external services, this mock service lets you simulate responses and conditions to validate
 *  how your code handles various scenarios.
 *
 *  @struct   MockUserService
 *  @inherits UserServiceInterface
 *
 *  @fields
 *  - SignupFunc (func): Customizes behavior for user signup.
 *  - LoginFunc (func): Customizes behavior for user login.
 *  - ResendOTPFunc (func): Customizes behavior for resending OTP emails.
 *  - VerifyEmailFunc (func): Customizes behavior for email verification.
 *  - ForgotPasswordFunc (func): Customizes password reset email behavior.
 *  - ResetPasswordFunc (func): Customizes behavior for resetting passwords.
 *  - GetUserInfoFunc (func): Customizes how user profile information is retrieved.
 *  - SearchUsersByUsernameFunc (func): Customizes user search results by username.
 *
 *  @behaviors
 *  - Returns errors if the corresponding function field is not set, ensuring clarity about missing
 *    mock implementations.
 *  - Allows granular control over the return values and error conditions for each method.
 *
 *  @example
 *  ```
 *  mockUserService := &MockUserService{
 *      LoginFunc: func(ctx context.Context, loginData *models.LoginRequest) (string, error) {
 *          if loginData.Email == "known@example.com" && loginData.Password == "validPass" {
 *              return "fake-jwt-token", nil
 *          }
 *          return "", fmt.Errorf("Invalid credentials")
 *      },
 *  }
 *
 *  // Use mockUserService in your tests and validate outcomes
 *  token, err := mockUserService.Login(context.Background(), &models.LoginRequest{
 *      Email: "known@example.com",
 *      Password: "validPass",
 *  })
 *  ```
 *
 *  @file      mock_user_service.go
 *  @project   DailyVerse
 *  @framework Go Testing with Mock Services
 */

package mocks

import (
	"context"
	"fmt"
	"proh2052-group6/pkg/models"
)

// MockUserService is a mock implementation of the UserServiceInterface.
type MockUserService struct {
	SignupFunc                func(ctx context.Context, user *models.User) error
	LoginFunc                 func(ctx context.Context, loginData *models.LoginRequest) (string, error)
	ResendOTPFunc             func(ctx context.Context, email string) error
	VerifyEmailFunc           func(ctx context.Context, email, otp string) (string, error)
	ForgotPasswordFunc        func(ctx context.Context, email string) error
	ResetPasswordFunc         func(ctx context.Context, email, otp, newPassword string) error
	GetUserInfoFunc           func(ctx context.Context, userEmail string) (map[string]string, error)
	SearchUsersByUsernameFunc func(ctx context.Context, userEmail, query string) ([]map[string]string, error)
}

// Signup mocks the Signup method of the UserServiceInterface.
func (m *MockUserService) Signup(ctx context.Context, user *models.User) error {
	if m.SignupFunc != nil {
		return m.SignupFunc(ctx, user)
	}
	return fmt.Errorf("SignupFunc not implemented")
}

// Login mocks the Login method, returning a token or an error.
func (m *MockUserService) Login(ctx context.Context, loginData *models.LoginRequest) (string, error) {
	if m.LoginFunc != nil {
		return m.LoginFunc(ctx, loginData)
	}
	return "", fmt.Errorf("LoginFunc not implemented")
}

// ResendOTP mocks the process of resending an OTP to the user.
func (m *MockUserService) ResendOTP(ctx context.Context, email string) error {
	if m.ResendOTPFunc != nil {
		return m.ResendOTPFunc(ctx, email)
	}
	return fmt.Errorf("ResendOTPFunc not implemented")
}

// VerifyEmail mocks the email verification process using an OTP.
func (m *MockUserService) VerifyEmail(ctx context.Context, email, otp string) (string, error) {
	if m.VerifyEmailFunc != nil {
		return m.VerifyEmailFunc(ctx, email, otp)
	}
	return "", fmt.Errorf("VerifyEmailFunc not implemented")
}

// ForgotPassword mocks sending a password reset OTP to the user’s email.
func (m *MockUserService) ForgotPassword(ctx context.Context, email string) error {
	if m.ForgotPasswordFunc != nil {
		return m.ForgotPasswordFunc(ctx, email)
	}
	return fmt.Errorf("ForgotPasswordFunc not implemented")
}

// ResetPassword mocks the password resetting process, validating the provided OTP.
func (m *MockUserService) ResetPassword(ctx context.Context, email, otp, newPassword string) error {
	if m.ResetPasswordFunc != nil {
		return m.ResetPasswordFunc(ctx, email, otp, newPassword)
	}
	return fmt.Errorf("ResetPasswordFunc not implemented")
}

// GetUserInfo mocks retrieving basic user information like email, username, country, etc.
func (m *MockUserService) GetUserInfo(ctx context.Context, userEmail string) (map[string]string, error) {
	if m.GetUserInfoFunc != nil {
		return m.GetUserInfoFunc(ctx, userEmail)
	}
	return nil, fmt.Errorf("GetUserInfoFunc not implemented")
}

// SearchUsersByUsername mocks searching for users by a query substring.
func (m *MockUserService) SearchUsersByUsername(ctx context.Context, userEmail, query string) ([]map[string]string, error) {
	if m.SearchUsersByUsernameFunc != nil {
		return m.SearchUsersByUsernameFunc(ctx, userEmail, query)
	}
	return nil, fmt.Errorf("SearchUsersByUsernameFunc not implemented")
}
