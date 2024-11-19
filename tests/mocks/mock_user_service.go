// tests/mocks/mock_user_service.go
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

// Signup mocks the Signup method.
func (m *MockUserService) Signup(ctx context.Context, user *models.User) error {
	if m.SignupFunc != nil {
		return m.SignupFunc(ctx, user)
	}
	return fmt.Errorf("SignupFunc not implemented")
}

// Login mocks the Login method.
func (m *MockUserService) Login(ctx context.Context, loginData *models.LoginRequest) (string, error) {
	if m.LoginFunc != nil {
		return m.LoginFunc(ctx, loginData)
	}
	return "", fmt.Errorf("LoginFunc not implemented")
}

// ResendOTP mocks the ResendOTP method.
func (m *MockUserService) ResendOTP(ctx context.Context, email string) error {
	if m.ResendOTPFunc != nil {
		return m.ResendOTPFunc(ctx, email)
	}
	return fmt.Errorf("ResendOTPFunc not implemented")
}

// VerifyEmail mocks the VerifyEmail method.
func (m *MockUserService) VerifyEmail(ctx context.Context, email, otp string) (string, error) {
	if m.VerifyEmailFunc != nil {
		return m.VerifyEmailFunc(ctx, email, otp)
	}
	return "", fmt.Errorf("VerifyEmailFunc not implemented")
}

// ForgotPassword mocks the ForgotPassword method.
func (m *MockUserService) ForgotPassword(ctx context.Context, email string) error {
	if m.ForgotPasswordFunc != nil {
		return m.ForgotPasswordFunc(ctx, email)
	}
	return fmt.Errorf("ForgotPasswordFunc not implemented")
}

// ResetPassword mocks the ResetPassword method.
func (m *MockUserService) ResetPassword(ctx context.Context, email, otp, newPassword string) error {
	if m.ResetPasswordFunc != nil {
		return m.ResetPasswordFunc(ctx, email, otp, newPassword)
	}
	return fmt.Errorf("ResetPasswordFunc not implemented")
}

// GetUserInfo mocks the GetUserInfo method.
func (m *MockUserService) GetUserInfo(ctx context.Context, userEmail string) (map[string]string, error) {
	if m.GetUserInfoFunc != nil {
		return m.GetUserInfoFunc(ctx, userEmail)
	}
	return nil, fmt.Errorf("GetUserInfoFunc not implemented")
}

// SearchUsersByUsername mocks the SearchUsersByUsername method.
func (m *MockUserService) SearchUsersByUsername(ctx context.Context, userEmail, query string) ([]map[string]string, error) {
	if m.SearchUsersByUsernameFunc != nil {
		return m.SearchUsersByUsernameFunc(ctx, userEmail, query)
	}
	return nil, fmt.Errorf("SearchUsersByUsernameFunc not implemented")
}
