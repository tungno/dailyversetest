/**
 *  FriendService provides business logic for managing friend relationships,
 *  including sending, accepting, removing, and retrieving friend requests.
 *
 *  @file       friend_service.go
 *  @package    services
 *
 *  @interfaces
 *  - FriendServiceInterface: Defines the contract for friend-related operations.
 *
 *  @methods
 *  - NewFriendService(userRepo, friendRepo): Initializes a new FriendService instance.
 *  - SendFriendRequest(ctx, userEmail, username): Sends a friend request to another user.
 *  - AcceptFriendRequest(ctx, userEmail, username): Accepts a received friend request.
 *  - GetFriendsList(ctx, userEmail): Retrieves the list of friends for a user.
 *  - RemoveFriend(ctx, userEmail, username): Removes a friendship.
 *  - GetPendingFriendRequests(ctx, userEmail): Retrieves pending friend requests for a user.
 *  - DeclineFriendRequest(ctx, userEmail, username): Declines a received friend request.
 *  - CancelFriendRequest(ctx, userEmail, username): Cancels a sent friend request.
 *
 *  @dependencies
 *  - repositories.UserRepository: Manages user-related data.
 *  - repositories.FriendRepository: Manages friend-related data.
 *  - utils.IsValidEmail: Utility function to validate email addresses.
 *
 *  @example
 *  ```
 *  friendService := NewFriendService(userRepo, friendRepo)
 *  err := friendService.SendFriendRequest(ctx, "user@example.com", "friend@example.com")
 *  if err != nil {
 *      log.Println("Failed to send friend request:", err)
 *  }
 *  ```
 *
 *  @behaviors
 *  - Validates input, ensuring users cannot send friend requests to themselves.
 *  - Prevents duplicate friend requests or relationships.
 *  - Supports friend operations by username or email.
 *  - Fetches user summaries for pending requests, excluding sensitive information.
 *
 *  @errors
 *  - Returns errors for invalid inputs, non-existent users, or database operation failures.
 *
 *  @authors
 *      - Aayush
 *      - Tung
 *      - Boss
 *      - Majd
 */

package services

import (
	"context"
	"fmt"
	"proh2052-group6/internal/repositories"
	"proh2052-group6/pkg/models"
	"proh2052-group6/pkg/utils"
)

// FriendServiceInterface defines methods for friend-related operations.
type FriendServiceInterface interface {
	SendFriendRequest(ctx context.Context, userEmail, username string) error
	AcceptFriendRequest(ctx context.Context, userEmail, username string) error
	GetFriendsList(ctx context.Context, userEmail string) ([]models.User, error)
	RemoveFriend(ctx context.Context, userEmail, username string) error
	GetPendingFriendRequests(ctx context.Context, userEmail string) ([]models.UserSummary, error)
	DeclineFriendRequest(ctx context.Context, userEmail, username string) error
	CancelFriendRequest(ctx context.Context, userEmail, username string) error
}

// FriendService implements FriendServiceInterface.
type FriendService struct {
	UserRepo   repositories.UserRepository   // Repository for user data.
	FriendRepo repositories.FriendRepository // Repository for friend data.
}

// NewFriendService initializes a new FriendService.
func NewFriendService(userRepo repositories.UserRepository, friendRepo repositories.FriendRepository) FriendServiceInterface {
	return &FriendService{
		UserRepo:   userRepo,
		FriendRepo: friendRepo,
	}
}

// SendFriendRequest sends a friend request to another user.
func (fs *FriendService) SendFriendRequest(ctx context.Context, userEmail, identifier string) error {
	var friendUser *models.User
	var err error

	// Determine if identifier is an email.
	if utils.IsValidEmail(identifier) {
		friendUser, err = fs.UserRepo.GetUserByEmail(ctx, identifier)
	} else {
		friendUser, err = fs.UserRepo.GetUserByUsername(ctx, identifier)
	}

	if err != nil || friendUser == nil {
		return fmt.Errorf("User not found")
	}

	friendEmail := friendUser.Email

	// Prevent sending a friend request to self.
	if userEmail == friendEmail {
		return fmt.Errorf("You cannot send a friend request to yourself")
	}

	// Check for existing friend requests or relationships.
	existingRequest, err := fs.FriendRepo.GetFriendRequest(ctx, userEmail, friendEmail)
	if err == nil && existingRequest != nil {
		return fmt.Errorf("Friend request already exists or you are already friends")
	}

	// Create a new friend request with "pending" status.
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

// AcceptFriendRequest accepts a pending friend request.
func (fs *FriendService) AcceptFriendRequest(ctx context.Context, userEmail, identifier string) error {
	var senderUser *models.User
	var err error

	// Get the sender user by username or email.
	senderUser, err = fs.UserRepo.GetUserByUsername(ctx, identifier)
	if err != nil || senderUser == nil {
		senderUser, err = fs.UserRepo.GetUserByEmail(ctx, identifier)
		if err != nil || senderUser == nil {
			return fmt.Errorf("User not found")
		}
	}
	senderEmail := senderUser.Email

	// Find the friend request sent by senderEmail to userEmail.
	existingRequest, err := fs.FriendRepo.GetFriendRequest(ctx, senderEmail, userEmail)
	if err != nil || existingRequest == nil {
		return fmt.Errorf("Friend request not found")
	}

	// Update the status of the request to "accepted".
	updates := map[string]interface{}{
		"Status": "accepted",
	}
	err = fs.FriendRepo.UpdateFriendRequest(ctx, senderEmail, userEmail, updates)
	if err != nil {
		return fmt.Errorf("Failed to accept friend request")
	}

	return nil
}

// GetFriendsList retrieves the list of friends for a user.
func (fs *FriendService) GetFriendsList(ctx context.Context, userEmail string) ([]models.User, error) {
	var friends []models.User

	// Fetch all accepted friend relationships.
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

		// Fetch user details of the friend.
		friendUser, err := fs.UserRepo.GetUserByEmail(ctx, friendEmail)
		if err != nil {
			continue
		}

		friends = append(friends, *friendUser)
	}

	return friends, nil
}

// RemoveFriend removes a friendship.
func (fs *FriendService) RemoveFriend(ctx context.Context, userEmail, username string) error {
	// Retrieve the friend's email.
	friendUser, err := fs.UserRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return fmt.Errorf("User not found")
	}
	friendEmail := friendUser.Email

	// Remove the friendship in both directions.
	err1 := fs.FriendRepo.DeleteFriendRequest(ctx, userEmail, friendEmail)
	err2 := fs.FriendRepo.DeleteFriendRequest(ctx, friendEmail, userEmail)

	if err1 != nil && err2 != nil {
		return fmt.Errorf("Failed to remove friend")
	}

	return nil
}

// GetPendingFriendRequests retrieves pending friend requests for a user.
func (fs *FriendService) GetPendingFriendRequests(ctx context.Context, userEmail string) ([]models.UserSummary, error) {
	friendRequests, err := fs.FriendRepo.GetPendingFriendRequests(ctx, userEmail)
	if err != nil {
		return nil, err
	}

	var pendingRequests []models.UserSummary
	for _, fr := range friendRequests {
		senderEmail := fr.Email

		// Fetch user details of the sender.
		user, err := fs.UserRepo.GetUserByEmail(ctx, senderEmail)
		if err != nil {
			continue
		}

		// Create a UserSummary for the pending request.
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

// DeclineFriendRequest declines a received friend request.
func (fs *FriendService) DeclineFriendRequest(ctx context.Context, userEmail, username string) error {
	senderUser, err := fs.UserRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return fmt.Errorf("User not found")
	}
	senderEmail := senderUser.Email

	// Delete the friend request.
	err = fs.FriendRepo.DeleteFriendRequest(ctx, senderEmail, userEmail)
	if err != nil {
		return fmt.Errorf("Failed to decline friend request")
	}

	return nil
}

// CancelFriendRequest cancels a sent friend request.
func (fs *FriendService) CancelFriendRequest(ctx context.Context, userEmail, username string) error {
	recipientUser, err := fs.UserRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return fmt.Errorf("User not found")
	}
	recipientEmail := recipientUser.Email

	// Delete the friend request.
	err = fs.FriendRepo.DeleteFriendRequest(ctx, userEmail, recipientEmail)
	if err != nil {
		return fmt.Errorf("Failed to cancel friend request")
	}

	return nil
}
