/**
 *  ProfileService provides business logic for managing user profiles.
 *  This service enables retrieval and updates of user profiles with validation and security measures.
 *
 *  @interface ProfileServiceInterface
 *  @methods
 *  - GetProfile(ctx, userEmail)                 - Retrieves the profile data for the specified user.
 *  - UpdateProfile(ctx, userEmail, updatedData) - Updates the profile data for the specified user.
 *
 *  @struct   ProfileService
 *  @inherits ProfileServiceInterface
 *
 *  @methods
 *  - NewProfileService(userRepo)               - Creates a new ProfileService instance with a user repository.
 *  - GetProfile(ctx, userEmail)                - Implementation for retrieving user profile data.
 *  - UpdateProfile(ctx, userEmail, updatedData)- Implementation for updating user profile data.
 *
 *  @behaviors
 *  - Ensures that user data is validated before updating the profile.
 *  - Validates the current password for sensitive updates, such as password changes.
 *  - Prevents updating protected fields like the email address.
 *  - Converts user data from struct to a map for JSON compatibility.
 *
 *  @dependencies
 *  - repositories.UserRepository: Repository for interacting with the Firestore user data.
 *  - utils: Utility package for password hashing, validation, and security checks.
 *
 *  @example
 *  ```
 *  // Retrieve profile data
 *  profile, err := profileService.GetProfile(ctx, "user@example.com")
 *
 *  // Update profile data
 *  updates := map[string]interface{}{
 *      "City": "Oslo",
 *      "CurrentPassword": "old_password",
 *      "NewPassword": "new_password123",
 *  }
 *  err := profileService.UpdateProfile(ctx, "user@example.com", updates)
 *  ```
 *
 *  @file      profile_service.go
 *  @project   DailyVerse
 *  @framework Go HTTP Server & Firestore API
 */

package services

import (
	"context"
	"fmt"

	"proh2052-group6/internal/repositories"
	"proh2052-group6/pkg/utils"
)

// ProfileServiceInterface defines the methods for managing user profiles.
type ProfileServiceInterface interface {
	GetProfile(ctx context.Context, userEmail string) (map[string]interface{}, error)
	UpdateProfile(ctx context.Context, userEmail string, updatedData map[string]interface{}) error
}

// ProfileService provides implementations for ProfileServiceInterface methods.
type ProfileService struct {
	UserRepo repositories.UserRepository
}

// NewProfileService initializes a new ProfileService with the given UserRepository.
func NewProfileService(userRepo repositories.UserRepository) ProfileServiceInterface {
	return &ProfileService{UserRepo: userRepo}
}

// GetProfile retrieves the profile data for the specified user.
func (ps *ProfileService) GetProfile(ctx context.Context, userEmail string) (map[string]interface{}, error) {
	// Fetch user data from the repository.
	user, err := ps.UserRepo.GetUserByEmail(ctx, userEmail)
	if err != nil {
		return nil, fmt.Errorf("Failed to get profile")
	}

	// Convert user struct to a map[string]interface{} for JSON compatibility.
	profileData := map[string]interface{}{
		"Email":    user.Email,
		"Username": user.Username,
		"Country":  user.Country,
		"City":     user.City,
		// Add other fields as required.
	}

	return profileData, nil
}

// UpdateProfile updates the profile data for the specified user with validation.
func (ps *ProfileService) UpdateProfile(ctx context.Context, userEmail string, updatedData map[string]interface{}) error {
	// Retrieve the current user data.
	user, err := ps.UserRepo.GetUserByEmail(ctx, userEmail)
	if err != nil {
		return fmt.Errorf("Failed to retrieve user data")
	}
	storedHashedPassword := user.Password

	// Validate the current password.
	currentPassword, ok := updatedData["CurrentPassword"].(string)
	if !ok || !utils.CheckPasswordHash(currentPassword, storedHashedPassword) {
		return fmt.Errorf("Invalid current password")
	}

	// Validate and update the password if a new password is provided.
	if newPassword, ok := updatedData["NewPassword"].(string); ok && newPassword != "" {
		if !utils.IsValidPassword(newPassword) {
			return fmt.Errorf("Password does not meet complexity requirements")
		}
		updatedData["Password"] = utils.HashPassword(newPassword)
	}

	// Remove fields that should not be updated directly.
	delete(updatedData, "CurrentPassword")
	delete(updatedData, "NewPassword")
	delete(updatedData, "Email") // Prevent updating the email address.

	// Update the user data in the repository.
	err = ps.UserRepo.UpdateUser(ctx, userEmail, updatedData)
	if err != nil {
		return fmt.Errorf("Failed to update profile")
	}

	return nil
}
