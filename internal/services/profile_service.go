// internal/services/profile_service.go
package services

import (
	"context"
	"fmt"

	"proh2052-group6/pkg/utils"

	"cloud.google.com/go/firestore"
)

type ProfileServiceInterface interface {
	GetProfile(ctx context.Context, userEmail string) (map[string]interface{}, error)
	UpdateProfile(ctx context.Context, userEmail string, updatedData map[string]interface{}) error
}

type ProfileService struct {
	DB DatabaseInterface
}

func NewProfileService(db DatabaseInterface) ProfileServiceInterface {
	return &ProfileService{DB: db}
}

func (ps *ProfileService) GetProfile(ctx context.Context, userEmail string) (map[string]interface{}, error) {
	doc, err := ps.DB.Collection("users").Doc(userEmail).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to get profile")
	}

	var profileData map[string]interface{}
	doc.DataTo(&profileData)
	return profileData, nil
}

func (ps *ProfileService) UpdateProfile(ctx context.Context, userEmail string, updatedData map[string]interface{}) error {
	doc, err := ps.DB.Collection("users").Doc(userEmail).Get(ctx)
	if err != nil {
		return fmt.Errorf("Failed to retrieve user data")
	}
	userData := doc.Data()
	storedHashedPassword := userData["Password"].(string)

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

	delete(updatedData, "CurrentPassword")
	delete(updatedData, "NewPassword")
	delete(updatedData, "Email") // Prevent email from being updated

	_, err = ps.DB.Collection("users").Doc(userEmail).Set(ctx, updatedData, firestore.MergeAll)
	if err != nil {
		return fmt.Errorf("Failed to update profile")
	}

	return nil
}
