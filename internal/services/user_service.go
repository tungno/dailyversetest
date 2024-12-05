/**
 *  UserService provides business logic for managing user accounts, including authentication,
 *  password recovery, email verification, and user search functionality. It integrates with
 *  the UserRepository and EmailService to perform operations.
 *
 *  @interface UserServiceInterface
 *  @inherits None
 *
 *  @methods
 *  - Signup(ctx, user)                      - Handles user registration with validation and email verification.
 *  - Login(ctx, loginData)                  - Authenticates a user and generates a JWT token.
 *  - ResendOTP(ctx, email)                  - Resends the OTP for email verification.
 *  - VerifyEmail(ctx, email, otp)           - Verifies a user's email using an OTP.
 *  - ForgotPassword(ctx, email)             - Sends an OTP to reset the user's password.
 *  - ResetPassword(ctx, email, otp, newPwd) - Resets the user's password using an OTP.
 *  - GetUserInfo(ctx, userEmail)            - Fetches the user's profile information.
 *  - SearchUsersByUsername(ctx, userEmail, query) - Searches for users by username.
 *
 *  @dependencies
 *  - repositories.UserRepository: Repository for interacting with user data in the database.
 *  - EmailServiceInterface: Service for sending emails to users.
 *  - utils: Utility package for password hashing, OTP generation, and JWT token handling.
 *
 *  @behaviors
 *  - Ensures secure handling of user data, including password hashing and OTP validation.
 *  - Provides detailed error messages for user-related operations.
 *  - Prevents unauthorized access by validating user inputs and tokens.
 *
 *  @example
 *  ```
 *  // Register a new user
 *  user := &models.User{
 *      Email: "user@example.com",
 *      Username: "JohnDoe",
 *      Country: "Norway",
 *      City: "Oslo",
 *      Password: "SecurePassword123",
 *  }
 *  err := userService.Signup(ctx, user)
 *
 *  // Login an existing user
 *  token, err := userService.Login(ctx, &models.LoginRequest{
 *      Email: "user@example.com",
 *      Password: "SecurePassword123",
 *  })
 *  ```
 *
 *  @file      user_service.go
 *  @project   DailyVerse
 *  @framework Go HTTP Server with Email Integration
 */

package services

import (
	"context"
	"fmt"
	"proh2052-group6/internal/repositories"
	"strings"
	"time"

	"proh2052-group6/pkg/models"
	"proh2052-group6/pkg/utils"
)

// UserServiceInterface defines the contract for user management operations.
type UserServiceInterface interface {
	Signup(ctx context.Context, user *models.User) error
	Login(ctx context.Context, loginData *models.LoginRequest) (string, error)
	ResendOTP(ctx context.Context, email string) error
	VerifyEmail(ctx context.Context, email, otp string) (string, error)
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, email, otp, newPassword string) error
	GetUserInfo(ctx context.Context, userEmail string) (map[string]string, error)
	SearchUsersByUsername(ctx context.Context, userEmail, query string) ([]map[string]string, error)
}

// UserService implements UserServiceInterface and interacts with repositories and email services.
type UserService struct {
	UserRepo repositories.UserRepository // Repository for user-related database operations.
	Email    EmailServiceInterface       // Email service for sending OTPs and notifications.
}

// NewUserService initializes a new UserService with a UserRepository and EmailService.
func NewUserService(userRepo repositories.UserRepository, emailService EmailServiceInterface) UserServiceInterface {
	return &UserService{
		UserRepo: userRepo,
		Email:    emailService,
	}
}

// Signup registers a new user with validation, OTP generation, and email verification.
func (us *UserService) Signup(ctx context.Context, user *models.User) error {
	if user.Country == "" || user.City == "" || user.Email == "" || user.Username == "" || user.Password == "" {
		return fmt.Errorf("Country, City, Email, Username, and Password are required")
	}

	existingUser, err := us.UserRepo.GetUserByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		return fmt.Errorf("Email already registered")
	}

	if !utils.IsValidPassword(user.Password) {
		return fmt.Errorf("Password does not meet complexity requirements")
	}

	user.Password = utils.HashPassword(user.Password)
	user.IsVerified = false
	user.UsernameLower = strings.ToLower(user.Username)
	user.OTP = utils.GenerateOTP()
	user.OTPExpiresAt = time.Now().Add(5 * time.Minute)

	if err := us.UserRepo.CreateUser(ctx, user); err != nil {
		return fmt.Errorf("Failed to create user: %v", err)
	}

	subject := "Your Verification Code"
	body := fmt.Sprintf("Your OTP for email verification is: %s. It will expire in 5 minutes.", user.OTP)
	if err := us.Email.SendEmail(user.Email, subject, body); err != nil {
		return fmt.Errorf("Failed to send verification email: %v", err)
	}

	return nil
}

// Login authenticates a user and returns a JWT token if successful.
func (us *UserService) Login(ctx context.Context, loginData *models.LoginRequest) (string, error) {
	user, err := us.UserRepo.GetUserByEmail(ctx, loginData.Email)
	if err != nil || user == nil {
		return "", fmt.Errorf("Email or password is incorrect")
	}

	if !user.IsVerified {
		return "", fmt.Errorf("Email not verified")
	}

	if utils.HashPassword(loginData.Password) != user.Password {
		return "", fmt.Errorf("Email or password is incorrect")
	}

	token, err := utils.GenerateJWT(user.Email)
	if err != nil {
		return "", fmt.Errorf("Failed to generate token")
	}

	return token, nil
}

// ResendOTP sends a new OTP to the user's email for verification.
func (us *UserService) ResendOTP(ctx context.Context, email string) error {
	user, err := us.UserRepo.GetUserByEmail(ctx, email)
	if err != nil || user == nil {
		return fmt.Errorf("Email not registered")
	}

	if user.IsVerified {
		return fmt.Errorf("Email is already verified")
	}

	user.OTP = utils.GenerateOTP()
	user.OTPExpiresAt = time.Now().Add(5 * time.Minute)

	updates := map[string]interface{}{
		"OTP":          user.OTP,
		"OTPExpiresAt": user.OTPExpiresAt,
	}
	if err := us.UserRepo.UpdateUser(ctx, email, updates); err != nil {
		return fmt.Errorf("Failed to update OTP")
	}

	subject := "Your New Verification Code"
	body := fmt.Sprintf("Your new OTP is: %s. It will expire in 5 minutes.", user.OTP)
	if err := us.Email.SendEmail(email, subject, body); err != nil {
		return fmt.Errorf("Failed to send OTP email")
	}

	return nil
}

// VerifyEmail verifies the user's email using the provided OTP and updates their status.
func (us *UserService) VerifyEmail(ctx context.Context, email, otp string) (string, error) {
	user, err := us.UserRepo.GetUserByEmail(ctx, email)
	if err != nil || user == nil {
		return "", fmt.Errorf("Invalid email or OTP")
	}

	if user.IsVerified {
		return "", fmt.Errorf("Email is already verified")
	}

	if user.OTP != otp {
		return "", fmt.Errorf("Invalid OTP")
	}

	if time.Now().After(user.OTPExpiresAt) {
		return "", fmt.Errorf("OTP has expired")
	}

	updates := map[string]interface{}{
		"IsVerified":   true,
		"OTP":          nil,
		"OTPExpiresAt": nil,
	}
	if err := us.UserRepo.UpdateUser(ctx, email, updates); err != nil {
		return "", fmt.Errorf("Failed to update user verification status")
	}

	token, err := utils.GenerateJWT(email)
	if err != nil {
		return "", fmt.Errorf("Failed to generate token")
	}

	return token, nil
}

func (us *UserService) ForgotPassword(ctx context.Context, email string) error {
	// Fetch user data
	user, err := us.UserRepo.GetUserByEmail(ctx, email)
	if err != nil || user == nil {
		// For security, we don't reveal whether the email exists
		return nil
	}

	// Generate OTP
	user.OTP = utils.GenerateOTP()
	user.OTPExpiresAt = time.Now().Add(5 * time.Minute)

	// Update the user with new OTP
	updates := map[string]interface{}{
		"OTP":          user.OTP,
		"OTPExpiresAt": user.OTPExpiresAt,
	}
	err = us.UserRepo.UpdateUser(ctx, email, updates)
	if err != nil {
		return fmt.Errorf("Failed to update OTP")
	}

	// Send OTP email
	subject := "Password Reset Request"
	body := fmt.Sprintf("Your OTP for password reset is: %s. It will expire in 5 minutes.", user.OTP)
	if err := us.Email.SendEmail(email, subject, body); err != nil {
		return fmt.Errorf("Failed to send OTP email")
	}

	return nil
}

func (us *UserService) ResetPassword(ctx context.Context, email, otp, newPassword string) error {
	user, err := us.UserRepo.GetUserByEmail(ctx, email)
	if err != nil || user == nil {
		return fmt.Errorf("Invalid email or OTP")
	}

	if user.OTP != otp {
		return fmt.Errorf("Invalid OTP")
	}

	if time.Now().After(user.OTPExpiresAt) {
		return fmt.Errorf("OTP has expired")
	}

	if !utils.IsValidPassword(newPassword) {
		return fmt.Errorf("Password does not meet complexity requirements")
	}

	hashedPassword := utils.HashPassword(newPassword)

	// Update the user's password and clear OTP
	updates := map[string]interface{}{
		"Password":     hashedPassword,
		"OTP":          nil,
		"OTPExpiresAt": nil,
	}
	err = us.UserRepo.UpdateUser(ctx, email, updates)
	if err != nil {
		return fmt.Errorf("Failed to reset password")
	}

	return nil
}

func (us *UserService) GetUserInfo(ctx context.Context, userEmail string) (map[string]string, error) {
	user, err := us.UserRepo.GetUserByEmail(ctx, userEmail)
	if err != nil || user == nil {
		return nil, fmt.Errorf("User not found")
	}

	userInfo := map[string]string{
		"email":    user.Email,
		"username": user.Username,
		"country":  user.Country,
		"city":     user.City,
	}

	return userInfo, nil
}

func (us *UserService) SearchUsersByUsername(ctx context.Context, userEmail, query string) ([]map[string]string, error) {
	users, err := us.UserRepo.SearchUsersByUsername(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("Failed to search users")
	}

	var results []map[string]string
	for _, user := range users {
		// Exclude the requesting user from the results
		if user.Email == userEmail {
			continue
		}

		results = append(results, map[string]string{
			"username": user.Username,
			"email":    user.Email,
		})
	}

	return results, nil
}
