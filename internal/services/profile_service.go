// internal/services/profile_service.go
package services

import (
	"context"
	"fmt"

	"proh2052-group6/internal/repositories"
	"proh2052-group6/pkg/utils"
)

type ProfileServiceInterface interface {
	GetProfile(ctx context.Context, userEmail string) (map[string]interface{}, error)
	UpdateProfile(ctx context.Context, userEmail string, updatedData map[string]interface{}) error
}

type ProfileService struct {
	UserRepo repositories.UserRepository
}

func NewProfileService(userRepo repositories.UserRepository) ProfileServiceInterface {
	return &ProfileService{UserRepo: userRepo}
}

func (ps *ProfileService) GetProfile(ctx context.Context, userEmail string) (map[string]interface{}, error) {
	user, err := ps.UserRepo.GetUserByEmail(ctx, userEmail)
	if err != nil {
		return nil, fmt.Errorf("Failed to get profile")
	}

	// Convert user struct to map[string]interface{}
	profileData := map[string]interface{}{
		"Email":    user.Email,
		"Username": user.Username,
		"Country":  user.Country,
		"City":     user.City,
		// Include other fields as needed
	}

	return profileData, nil
}

func (ps *ProfileService) UpdateProfile(ctx context.Context, userEmail string, updatedData map[string]interface{}) error {
	user, err := ps.UserRepo.GetUserByEmail(ctx, userEmail)
	if err != nil {
		return fmt.Errorf("Failed to retrieve user data")
	}
	storedHashedPassword := user.Password

	// Ensure the current password is provided for all updates
	currentPassword, ok := updatedData["CurrentPassword"].(string)
	if !ok || utils.HashPassword(currentPassword) != storedHashedPassword {
		return fmt.Errorf("Invalid current password")
	}

	// Update password if new password is provided
	if newPassword, ok := updatedData["NewPassword"].(string); ok && newPassword != "" {
		if !utils.IsValidPassword(newPassword) {
			return fmt.Errorf("Password does not meet complexity requirements")
		}
		updatedData["Password"] = utils.HashPassword(newPassword)
	}

	// Remove fields that should not be updated directly
	delete(updatedData, "CurrentPassword")
	delete(updatedData, "NewPassword")
	delete(updatedData, "Email") // Prevent email from being updated

	err = ps.UserRepo.UpdateUser(ctx, userEmail, updatedData)
	if err != nil {
		return fmt.Errorf("Failed to update profile")
	}

	return nil
}
