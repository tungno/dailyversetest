/**
 *  MockProfileService simulates a profile service for testing purposes.
 *  It provides in-memory operations to manage user profiles, including retrieval
 *  and updates. This mock service allows testing handlers and services without
 *  relying on an actual database or profile service.
 *
 *  @file       mock_profile_service.go
 *  @package    mocks
 *
 *  @structs
 *  - MockProfileService: Simulates a profile service with an in-memory store for user profiles.
 *
 *  @methods
 *  - NewMockProfileService: Initializes a new instance of MockProfileService.
 *  - GetProfile(ctx, userEmail): Simulates retrieving a user profile by email.
 *  - UpdateProfile(ctx, userEmail, updatedData): Simulates updating a user's profile.
 *
 *  @example
 *  ```
 *  mockService := mocks.NewMockProfileService()
 *  profile := map[string]interface{}{
 *      "Email":    "user@example.com",
 *      "Username": "testuser",
 *      "Country":  "TestCountry",
 *      "City":     "TestCity",
 *      "Password": "hashedpassword",
 *  }
 *  mockService.Profiles["user@example.com"] = profile
 *
 *  updatedData := map[string]interface{}{
 *      "Username":        "updateduser",
 *      "CurrentPassword": "hashedpassword",
 *      "NewPassword":     "newhashedpassword",
 *  }
 *
 *  err := mockService.UpdateProfile(context.Background(), "user@example.com", updatedData)
 *  if err != nil {
 *      t.Fatalf("Failed to update profile: %v", err)
 *  }
 *  ```
 *
 *  @dependencies
 *  - context: Context for managing request lifecycle in services.
 *
 *  @limitations
 *  - MockProfileService is in-memory and does not persist data across tests.
 *  - Password hashing is simulated without using secure hashing mechanisms.
 *
 *  @authors
 *      - Aayush
 *      - Tung
 *      - Boss
 *      - Majd
 */

package mocks

import (
	"context"
	"errors"
)

// MockProfileService simulates a profile service for testing.
type MockProfileService struct {
	Profiles map[string]map[string]interface{} // In-memory store for profiles.
	Users    map[string]map[string]interface{} // In-memory store for users.
}

// NewMockProfileService initializes a new instance of MockProfileService.
func NewMockProfileService() *MockProfileService {
	return &MockProfileService{
		Profiles: make(map[string]map[string]interface{}),
		Users:    make(map[string]map[string]interface{}),
	}
}

// GetProfile simulates retrieving a user profile by email.
func (mps *MockProfileService) GetProfile(ctx context.Context, userEmail string) (map[string]interface{}, error) {
	profile, exists := mps.Profiles[userEmail]
	if !exists {
		return nil, errors.New("profile not found")
	}
	return profile, nil
}

// UpdateProfile simulates updating a user's profile.
func (mps *MockProfileService) UpdateProfile(ctx context.Context, userEmail string, updatedData map[string]interface{}) error {
	profile, exists := mps.Profiles[userEmail]
	if !exists {
		return errors.New("profile not found")
	}

	// Simulate password validation.
	currentPassword, ok := updatedData["CurrentPassword"].(string)
	if !ok || currentPassword != profile["Password"] {
		return errors.New("invalid current password")
	}

	// Update the profile with new data.
	for key, value := range updatedData {
		switch key {
		case "CurrentPassword":
			// Skip updating current password.
		case "NewPassword":
			// Simulate updating the password (no actual hashing in mock).
			if newPassword, ok := value.(string); ok && newPassword != "" {
				profile["Password"] = newPassword
			}
		default:
			// Update other fields.
			profile[key] = value
		}
	}

	return nil
}
