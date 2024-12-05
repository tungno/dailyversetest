// internal/services/friend_service.go
package services

import (
	"context"
	"fmt"
	"proh2052-group6/internal/repositories"
	"proh2052-group6/pkg/models"
	"proh2052-group6/pkg/utils"
)

type FriendServiceInterface interface {
	SendFriendRequest(ctx context.Context, userEmail, username string) error
	AcceptFriendRequest(ctx context.Context, userEmail, username string) error
	GetFriendsList(ctx context.Context, userEmail string) ([]models.User, error)
	RemoveFriend(ctx context.Context, userEmail, username string) error
	GetPendingFriendRequests(ctx context.Context, userEmail string) ([]models.UserSummary, error)
	DeclineFriendRequest(ctx context.Context, userEmail, username string) error
	CancelFriendRequest(ctx context.Context, userEmail, username string) error
}

type FriendService struct {
	UserRepo   repositories.UserRepository
	FriendRepo repositories.FriendRepository
}

func NewFriendService(userRepo repositories.UserRepository, friendRepo repositories.FriendRepository) FriendServiceInterface {
	return &FriendService{
		UserRepo:   userRepo,
		FriendRepo: friendRepo,
	}
}

func (fs *FriendService) SendFriendRequest(ctx context.Context, userEmail, identifier string) error {
	var friendUser *models.User
	var err error

	// Determine if identifier is an email
	if utils.IsValidEmail(identifier) {
		friendUser, err = fs.UserRepo.GetUserByEmail(ctx, identifier)
	} else {
		friendUser, err = fs.UserRepo.GetUserByUsername(ctx, identifier)
	}

	if err != nil || friendUser == nil {
		return fmt.Errorf("User not found")
	}

	friendEmail := friendUser.Email

	// Prevent sending a friend request to self
	if userEmail == friendEmail {
		return fmt.Errorf("You cannot send a friend request to yourself")
	}

	// Check if a friend request or relationship already exists
	existingRequest, err := fs.FriendRepo.GetFriendRequest(ctx, userEmail, friendEmail)
	if err == nil && existingRequest != nil {
		return fmt.Errorf("Friend request already exists or you are already friends")
	}

	// Create new friend request (pending)
	friendRequest := &models.Friend{
		Email:       userEmail,
		FriendEmail: friendEmail,
		Status:      "pending",
	}
	err = fs.FriendRepo.CreateFriendRequest(ctx, friendRequest)
	if err != nil {
		return fmt.Errorf("Failed to send friend request")
	}

	return nil
}

func (fs *FriendService) AcceptFriendRequest(ctx context.Context, userEmail, identifier string) error {
	var senderUser *models.User
	var err error

	// Try to get user by username first
	senderUser, err = fs.UserRepo.GetUserByUsername(ctx, identifier)
	if err != nil || senderUser == nil {
		// If not found, try by email
		senderUser, err = fs.UserRepo.GetUserByEmail(ctx, identifier)
		if err != nil || senderUser == nil {
			return fmt.Errorf("User not found")
		}
	}
	senderEmail := senderUser.Email

	// Find the friend request where sender is senderEmail and recipient is userEmail
	existingRequest, err := fs.FriendRepo.GetFriendRequest(ctx, senderEmail, userEmail)
	if err != nil || existingRequest == nil {
		return fmt.Errorf("Friend request not found")
	}

	// Update the status to "accepted"
	updates := map[string]interface{}{
		"Status": "accepted",
	}
	err = fs.FriendRepo.UpdateFriendRequest(ctx, senderEmail, userEmail, updates)
	if err != nil {
		return fmt.Errorf("Failed to accept friend request")
	}

	return nil
}

func (fs *FriendService) GetFriendsList(ctx context.Context, userEmail string) ([]models.User, error) {
	var friends []models.User

	// Get all accepted friend relationships involving the user
	friendRelations, err := fs.FriendRepo.GetFriends(ctx, userEmail)
	if err != nil {
		return nil, fmt.Errorf("Error fetching friends list")
	}

	for _, friendRelation := range friendRelations {
		var friendEmail string
		if friendRelation.Email == userEmail {
			friendEmail = friendRelation.FriendEmail
		} else {
			friendEmail = friendRelation.Email
		}

		// Fetch user data
		friendUser, err := fs.UserRepo.GetUserByEmail(ctx, friendEmail)
		if err != nil {
			continue
		}

		friends = append(friends, *friendUser)
	}

	return friends, nil
}

func (fs *FriendService) RemoveFriend(ctx context.Context, userEmail, username string) error {
	// Retrieve the email of the user by username
	friendUser, err := fs.UserRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return fmt.Errorf("User not found")
	}
	friendEmail := friendUser.Email

	// Remove the friendship documents in both directions
	err1 := fs.FriendRepo.DeleteFriendRequest(ctx, userEmail, friendEmail)
	err2 := fs.FriendRepo.DeleteFriendRequest(ctx, friendEmail, userEmail)

	if err1 != nil && err2 != nil {
		return fmt.Errorf("Failed to remove friend")
	}

	return nil
}

// GetPendingFriendRequests retrieves pending friend requests for a user
func (fs *FriendService) GetPendingFriendRequests(ctx context.Context, userEmail string) ([]models.UserSummary, error) {
	// Fetch pending friend requests where the recipient is the user
	friendRequests, err := fs.FriendRepo.GetPendingFriendRequests(ctx, userEmail)
	if err != nil {
		return nil, err
	}

	var pendingRequests []models.UserSummary
	for _, fr := range friendRequests {
		// Get the sender's email (the person who sent the request)
		senderEmail := fr.Email

		// Fetch user details of the sender
		user, err := fs.UserRepo.GetUserByEmail(ctx, senderEmail)
		if err != nil {
			// Skip if user not found or any error occurs
			continue
		}

		// Create a UserSummary to exclude sensitive fields
		userSummary := models.UserSummary{
			Username: user.Username,
			Email:    user.Email,
			Country:  user.Country,
			City:     user.City,
		}

		pendingRequests = append(pendingRequests, userSummary)
	}

	return pendingRequests, nil
}

func (fs *FriendService) DeclineFriendRequest(ctx context.Context, userEmail, username string) error {
	// Retrieve the email of the user by username
	senderUser, err := fs.UserRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return fmt.Errorf("User not found")
	}
	senderEmail := senderUser.Email

	// Delete the friend request where sender is senderEmail and recipient is userEmail
	err = fs.FriendRepo.DeleteFriendRequest(ctx, senderEmail, userEmail)
	if err != nil {
		return fmt.Errorf("Failed to decline friend request")
	}

	return nil
}

func (fs *FriendService) CancelFriendRequest(ctx context.Context, userEmail, username string) error {
	// Retrieve the email of the user by username
	recipientUser, err := fs.UserRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return fmt.Errorf("User not found")
	}
	recipientEmail := recipientUser.Email

	// Delete the friend request where sender is userEmail and recipient is recipientEmail
	err = fs.FriendRepo.DeleteFriendRequest(ctx, userEmail, recipientEmail)
	if err != nil {
		return fmt.Errorf("Failed to cancel friend request")
	}

	return nil
}
