// internal/services/user_service.go
package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"proh2052-group6/pkg/models"
	"proh2052-group6/pkg/utils"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

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

type UserService struct {
	DB    DatabaseInterface
	Email EmailServiceInterface
}

func NewUserService(db DatabaseInterface, email EmailServiceInterface) UserServiceInterface {
	return &UserService{
		DB:    db,
		Email: email,
	}
}

func (us *UserService) Signup(ctx context.Context, user *models.User) error {
	// Validate inputs
	if user.Country == "" || user.City == "" || user.Email == "" || user.Username == "" || user.Password == "" {
		return fmt.Errorf("Country, City, Email, Username, and Password are required")
	}

	// Check if email already exists
	doc, err := us.DB.Collection("users").Doc(user.Email).Get(ctx)
	if err == nil && doc.Exists() {
		return fmt.Errorf("Email already registered")
	}

	// Validate password
	if !utils.IsValidPassword(user.Password) {
		return fmt.Errorf("Password does not meet complexity requirements")
	}

	// Hash password
	user.Password = utils.HashPassword(user.Password)
	user.IsVerified = false
	user.UsernameLower = strings.ToLower(user.Username)

	// Generate OTP
	user.OTP = utils.GenerateOTP()
	user.OTPExpiresAt = time.Now().Add(5 * time.Minute)

	// Save user to DB
	_, err = us.DB.Collection("users").Doc(user.Email).Set(ctx, user)
	if err != nil {
		return fmt.Errorf("Failed to create user: %v", err)
	}

	// Send OTP email
	subject := "Your Verification Code"
	body := fmt.Sprintf("Your OTP for email verification is: %s. It will expire in 5 minutes.", user.OTP)
	if err := us.Email.SendEmail(user.Email, subject, body); err != nil {
		return fmt.Errorf("Failed to send verification email: %v", err)
	}

	return nil
}

func (us *UserService) Login(ctx context.Context, loginData *models.LoginRequest) (string, error) {
	doc, err := us.DB.Collection("users").Doc(loginData.Email).Get(ctx)
	if err != nil || !doc.Exists() {
		return "", fmt.Errorf("Email or password is incorrect")
	}

	var user models.User
	if err := doc.DataTo(&user); err != nil {
		return "", fmt.Errorf("Failed to parse user data")
	}

	if !user.IsVerified {
		return "", fmt.Errorf("Email not verified")
	}

	if utils.HashPassword(loginData.Password) != user.Password {
		return "", fmt.Errorf("Email or password is incorrect")
	}

	// Generate JWT
	token, err := utils.GenerateJWT(user.Email)
	if err != nil {
		return "", fmt.Errorf("Failed to generate token")
	}

	return token, nil
}

func (us *UserService) ResendOTP(ctx context.Context, email string) error {
	// Fetch user data
	doc, err := us.DB.Collection("users").Doc(email).Get(ctx)
	if err != nil || !doc.Exists() {
		return fmt.Errorf("Email not registered")
	}

	var user models.User
	if err := doc.DataTo(&user); err != nil {
		return fmt.Errorf("Failed to parse user data")
	}

	if user.IsVerified {
		return fmt.Errorf("Email is already verified")
	}

	// Generate new OTP
	user.OTP = utils.GenerateOTP()
	user.OTPExpiresAt = time.Now().Add(5 * time.Minute)

	// Update the user with new OTP
	_, err = us.DB.Collection("users").Doc(email).Set(ctx, map[string]interface{}{
		"OTP":          user.OTP,
		"OTPExpiresAt": user.OTPExpiresAt,
	}, firestore.MergeAll)
	if err != nil {
		return fmt.Errorf("Failed to update OTP")
	}

	// Send OTP email
	subject := "Your New Verification Code"
	body := fmt.Sprintf("Your new OTP is: %s. It will expire in 5 minutes.", user.OTP)
	if err := us.Email.SendEmail(email, subject, body); err != nil {
		return fmt.Errorf("Failed to send OTP email")
	}

	return nil
}

func (us *UserService) VerifyEmail(ctx context.Context, email, otp string) (string, error) {
	doc, err := us.DB.Collection("users").Doc(email).Get(ctx)
	if err != nil || !doc.Exists() {
		return "", fmt.Errorf("Invalid email or OTP")
	}

	var user models.User
	if err := doc.DataTo(&user); err != nil {
		return "", fmt.Errorf("Failed to parse user data")
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

	// Update user as verified
	_, err = us.DB.Collection("users").Doc(email).Set(ctx, map[string]interface{}{
		"IsVerified": true,
		"OTP":        nil,
	}, firestore.MergeAll)
	if err != nil {
		return "", fmt.Errorf("Failed to update user verification status")
	}

	// Generate JWT
	token, err := utils.GenerateJWT(email)
	if err != nil {
		return "", fmt.Errorf("Failed to generate token")
	}

	return token, nil
}

func (us *UserService) ForgotPassword(ctx context.Context, email string) error {
	// Fetch user data
	doc, err := us.DB.Collection("users").Doc(email).Get(ctx)
	if err != nil || !doc.Exists() {
		// For security, we don't reveal whether the email exists
		return nil
	}

	var user models.User
	if err := doc.DataTo(&user); err != nil {
		return nil
	}

	// Generate OTP
	user.OTP = utils.GenerateOTP()
	user.OTPExpiresAt = time.Now().Add(5 * time.Minute)

	// Update the user with new OTP
	_, err = us.DB.Collection("users").Doc(email).Set(ctx, map[string]interface{}{
		"OTP":          user.OTP,
		"OTPExpiresAt": user.OTPExpiresAt,
	}, firestore.MergeAll)
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
	doc, err := us.DB.Collection("users").Doc(email).Get(ctx)
	if err != nil || !doc.Exists() {
		return fmt.Errorf("Invalid email or OTP")
	}

	var user models.User
	if err := doc.DataTo(&user); err != nil {
		return fmt.Errorf("Failed to parse user data")
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
	_, err = us.DB.Collection("users").Doc(email).Set(ctx, map[string]interface{}{
		"Password":     hashedPassword,
		"OTP":          nil,
		"OTPExpiresAt": nil,
	}, firestore.MergeAll)
	if err != nil {
		return fmt.Errorf("Failed to reset password")
	}

	return nil
}

func (us *UserService) GetUserInfo(ctx context.Context, userEmail string) (map[string]string, error) {
	doc, err := us.DB.Collection("users").Doc(userEmail).Get(ctx)
	if err != nil || !doc.Exists() {
		return nil, fmt.Errorf("User not found")
	}

	var user models.User
	if err := doc.DataTo(&user); err != nil {
		return nil, fmt.Errorf("Failed to parse user data")
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
	iter := us.DB.Collection("users").Where("UsernameLower", ">=", strings.ToLower(query)).Where("UsernameLower", "<=", strings.ToLower(query)+"\uf8ff").Documents(ctx)
	defer iter.Stop()

	var results []map[string]string
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Failed to search users")
		}

		var user models.User
		if err := doc.DataTo(&user); err != nil {
			continue
		}

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
