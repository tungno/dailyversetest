// tests/mocks/mock_profile_service.go
package mocks

import (
	"context"
	"errors"
)

type MockProfileService struct {
	Profiles map[string]map[string]interface{}
	Users    map[string]map[string]interface{}
}

func NewMockProfileService() *MockProfileService {
	return &MockProfileService{
		Profiles: make(map[string]map[string]interface{}),
		Users:    make(map[string]map[string]interface{}),
	}
}

func (mps *MockProfileService) GetProfile(ctx context.Context, userEmail string) (map[string]interface{}, error) {
	profile, exists := mps.Profiles[userEmail]
	if !exists {
		return nil, errors.New("profile not found")
	}
	return profile, nil
}

func (mps *MockProfileService) UpdateProfile(ctx context.Context, userEmail string, updatedData map[string]interface{}) error {
	profile, exists := mps.Profiles[userEmail]
	if !exists {
		return errors.New("profile not found")
	}

	// Simulate password hashing
	currentPassword, ok := updatedData["CurrentPassword"].(string)
	if !ok || currentPassword != profile["Password"] {
		return errors.New("invalid current password")
	}

	// Update the profile with new data
	for key, value := range updatedData {
		switch key {
		case "CurrentPassword":
			// Skip
		case "NewPassword":
			if newPassword, ok := value.(string); ok && newPassword != "" {
				profile["Password"] = newPassword // Simulate hashing
			}
		default:
			profile[key] = value
		}
	}

	return nil
}
