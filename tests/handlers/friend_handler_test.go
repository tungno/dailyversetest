/**
 *  FriendHandler Test Suite
 *
 *  This test suite validates the behavior of the `FriendHandler`, ensuring proper functionality
 *  for friend-related operations, including sending, accepting, declining, and removing friend requests,
 *  as well as retrieving friend lists and pending requests. The tests use mocked repositories and services
 *  to isolate the handler logic from external dependencies.
 *
 *  @dependencies
 *  - mocks.MockUserRepository: Mock implementation of the UserRepository for simulating user data.
 *  - mocks.MockFriendRepository: Mock implementation of the FriendRepository for simulating friend relationships.
 *  - services.FriendService: Core service for managing friend-related operations.
 *  - handlers.FriendHandler: HTTP handler for friend endpoints.
 *  - httptest: Go's HTTP testing package for simulating HTTP requests and responses.
 *  - testify/assert: Library for making test assertions clean and readable.
 *
 *  @testcases
 *  - TestSendFriendRequestHandler: Validates the ability to send a friend request.
 *  - TestAcceptFriendRequestHandler: Verifies that pending friend requests can be accepted.
 *  - TestGetFriendsListHandler: Checks the retrieval of a user's accepted friend list.
 *  - TestRemoveFriendHandler: Ensures a user can remove an existing friend.
 *  - TestGetPendingFriendRequestsHandler: Validates retrieval of pending friend requests.
 *  - TestDeclineFriendRequestHandler: Confirms that a user can decline a pending friend request.
 *  - TestCancelFriendRequestHandler: Tests the ability to cancel a sent friend request.
 *
 *  @behaviors
 *  - Uses mock repositories to simulate user and friend data for isolated testing.
 *  - Ensures correct HTTP status codes are returned based on request outcomes.
 *  - Verifies that appropriate changes are reflected in the mock repositories.
 *  - Handles edge cases, such as missing users or invalid friend request states.
 *
 *  @example
 *  ```
 *  // Test sending a friend request
 *  mockUsers := map[string]*models.User{
 *      "user1@example.com": {Email: "user1@example.com", Username: "user1"},
 *      "user2@example.com": {Email: "user2@example.com", Username: "user2"},
 *  }
 *  userRepo := mocks.NewMockUserRepository(mockUsers)
 *  friendRepo := mocks.NewMockFriendRepository(make(map[string]*models.Friend))
 *
 *  friendService := services.NewFriendService(userRepo, friendRepo)
 *  friendHandler := handlers.NewFriendHandler(friendService)
 *
 *  req, _ := http.NewRequest("POST", "/api/friends/add", bytes.NewReader(body))
 *  ctx := context.WithValue(req.Context(), "userEmail", "user1@example.com")
 *  req = req.WithContext(ctx)
 *
 *  rr := httptest.NewRecorder()
 *  http.HandlerFunc(friendHandler.SendFriendRequest).ServeHTTP(rr, req)
 *  ```
 *
 *  @file      friend_handler_test.go
 *  @project   DailyVerse
 *  @framework Go HTTP Testing with Mock Services
 */

package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"proh2052-group6/internal/handlers"
	"proh2052-group6/internal/services"
	"proh2052-group6/pkg/models"
	"proh2052-group6/tests/mocks"
)

func TestSendFriendRequestHandler(t *testing.T) {
	mockUsers := map[string]*models.User{
		"user1@example.com": {Email: "user1@example.com", Username: "user1"},
		"user2@example.com": {Email: "user2@example.com", Username: "user2"},
	}
	userRepo := mocks.NewMockUserRepository(mockUsers)
	friendRepo := mocks.NewMockFriendRepository(make(map[string]*models.Friend))

	friendService := services.NewFriendService(userRepo, friendRepo)
	friendHandler := handlers.NewFriendHandler(friendService)

	requestData := map[string]string{
		"username": "user2",
	}
	body, _ := json.Marshal(requestData)
	req, err := http.NewRequest("POST", "/api/friends/add", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	// Mock authentication context
	ctx := context.WithValue(req.Context(), "userEmail", "user1@example.com")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(friendHandler.SendFriendRequest)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Verify response body
	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response body")
	}
	if response["message"] != "Friend request sent" {
		t.Errorf("Unexpected response message: %s", response["message"])
	}

	// Verify that the friend request was created
	friendKey := "user1@example.com_user2@example.com"
	if _, exists := friendRepo.Friends[friendKey]; !exists {
		t.Errorf("Friend request not found in mock repository")
	}
}

func TestAcceptFriendRequestHandler(t *testing.T) {
	mockUsers := map[string]*models.User{
		"user1@example.com": {Email: "user1@example.com", Username: "user1"},
		"user2@example.com": {Email: "user2@example.com", Username: "user2"},
	}
	userRepo := mocks.NewMockUserRepository(mockUsers)
	friendRepo := mocks.NewMockFriendRepository(map[string]*models.Friend{
		"user2@example.com_user1@example.com": {
			Email:       "user2@example.com",
			FriendEmail: "user1@example.com",
			Status:      "pending",
		},
	})

	friendService := services.NewFriendService(userRepo, friendRepo)
	friendHandler := handlers.NewFriendHandler(friendService)

	requestData := map[string]string{
		"username": "user2",
	}
	body, _ := json.Marshal(requestData)
	req, err := http.NewRequest("POST", "/api/friends/accept", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	// Mock authentication context
	ctx := context.WithValue(req.Context(), "userEmail", "user1@example.com")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(friendHandler.AcceptFriendRequest)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Verify response body
	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response body")
	}
	if response["message"] != "Friend request accepted" {
		t.Errorf("Unexpected response message: %s", response["message"])
	}

	// Verify that the friend request status has been updated
	friendKey := "user2@example.com_user1@example.com"
	friend, exists := friendRepo.Friends[friendKey]
	if !exists {
		t.Errorf("Friend request not found in mock repository")
	} else if friend.Status != "accepted" {
		t.Errorf("Friend request status not updated to 'accepted'")
	}
}

func TestGetFriendsListHandler(t *testing.T) {
	mockUsers := map[string]*models.User{
		"user1@example.com": {Email: "user1@example.com", Username: "user1"},
		"user2@example.com": {Email: "user2@example.com", Username: "user2"},
		"user3@example.com": {Email: "user3@example.com", Username: "user3"},
	}
	userRepo := mocks.NewMockUserRepository(mockUsers)
	friendRepo := mocks.NewMockFriendRepository(map[string]*models.Friend{
		"user1@example.com_user2@example.com": {
			Email:       "user1@example.com",
			FriendEmail: "user2@example.com",
			Status:      "accepted",
		},
		"user3@example.com_user1@example.com": {
			Email:       "user3@example.com",
			FriendEmail: "user1@example.com",
			Status:      "accepted",
		},
	})
	friendService := services.NewFriendService(userRepo, friendRepo)
	friendHandler := handlers.NewFriendHandler(friendService)

	req, err := http.NewRequest("GET", "/api/friends/list", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Mock authentication context
	ctx := context.WithValue(req.Context(), "userEmail", "user1@example.com")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(friendHandler.GetFriendsList)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Verify response body
	var friends []models.User
	err = json.Unmarshal(rr.Body.Bytes(), &friends)
	if err != nil {
		t.Errorf("Failed to parse response body")
	}

	if len(friends) != 2 {
		t.Errorf("Expected 2 friends, got %d", len(friends))
	}

	// Check that the friends are the correct users
	friendEmails := map[string]bool{
		"user2@example.com": false,
		"user3@example.com": false,
	}
	for _, friend := range friends {
		if _, exists := friendEmails[friend.Email]; exists {
			friendEmails[friend.Email] = true
		} else {
			t.Errorf("Unexpected friend email: %s", friend.Email)
		}
	}
	for email, found := range friendEmails {
		if !found {
			t.Errorf("Friend %s not found in response", email)
		}
	}
}

func TestRemoveFriendHandler(t *testing.T) {
	mockUsers := map[string]*models.User{
		"user1@example.com": {Email: "user1@example.com", Username: "user1"},
		"user2@example.com": {Email: "user2@example.com", Username: "user2"},
	}
	userRepo := mocks.NewMockUserRepository(mockUsers)
	friendRepo := mocks.NewMockFriendRepository(map[string]*models.Friend{
		"user1@example.com_user2@example.com": {
			Email:       "user1@example.com",
			FriendEmail: "user2@example.com",
			Status:      "accepted",
		},
	})
	friendService := services.NewFriendService(userRepo, friendRepo)
	friendHandler := handlers.NewFriendHandler(friendService)

	requestData := map[string]string{
		"username": "user2",
	}
	body, _ := json.Marshal(requestData)
	req, err := http.NewRequest("POST", "/api/friends/remove", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	// Mock authentication context
	ctx := context.WithValue(req.Context(), "userEmail", "user1@example.com")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(friendHandler.RemoveFriend)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Verify response body
	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response body")
	}
	if response["message"] != "Friend removed" {
		t.Errorf("Unexpected response message: %s", response["message"])
	}

	// Verify that the friend relationship has been removed
	friendKey1 := "user1@example.com_user2@example.com"
	friendKey2 := "user2@example.com_user1@example.com"
	if _, exists := friendRepo.Friends[friendKey1]; exists {
		t.Errorf("Friend relationship not removed from mock repository (key: %s)", friendKey1)
	}
	if _, exists := friendRepo.Friends[friendKey2]; exists {
		t.Errorf("Friend relationship not removed from mock repository (key: %s)", friendKey2)
	}
}

func TestGetPendingFriendRequestsHandler(t *testing.T) {
	mockUsers := map[string]*models.User{
		"user1@example.com": {Email: "user1@example.com", Username: "user1"},
		"user2@example.com": {Email: "user2@example.com", Username: "user2"},
		"user3@example.com": {Email: "user3@example.com", Username: "user3"},
	}
	userRepo := mocks.NewMockUserRepository(mockUsers)
	friendRepo := mocks.NewMockFriendRepository(map[string]*models.Friend{
		"user2@example.com_user1@example.com": {
			Email:       "user2@example.com",
			FriendEmail: "user1@example.com",
			Status:      "pending",
		},
		"user3@example.com_user1@example.com": {
			Email:       "user3@example.com",
			FriendEmail: "user1@example.com",
			Status:      "pending",
		},
	})
	friendService := services.NewFriendService(userRepo, friendRepo)
	friendHandler := handlers.NewFriendHandler(friendService)

	req, err := http.NewRequest("GET", "/api/friends/requests", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Mock authentication context
	ctx := context.WithValue(req.Context(), "userEmail", "user1@example.com")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(friendHandler.GetPendingFriendRequests)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Verify response body
	var requests []models.User
	err = json.Unmarshal(rr.Body.Bytes(), &requests)
	if err != nil {
		t.Errorf("Failed to parse response body")
	}

	if len(requests) != 2 {
		t.Errorf("Expected 2 pending requests, got %d", len(requests))
	}

	// Check that the requests are from the correct users
	requestEmails := map[string]bool{
		"user2@example.com": false,
		"user3@example.com": false,
	}
	for _, user := range requests {
		if _, exists := requestEmails[user.Email]; exists {
			requestEmails[user.Email] = true
		} else {
			t.Errorf("Unexpected request from email: %s", user.Email)
		}
	}
	for email, found := range requestEmails {
		if !found {
			t.Errorf("Pending request from %s not found in response", email)
		}
	}
}

func TestDeclineFriendRequestHandler(t *testing.T) {
	mockUsers := map[string]*models.User{
		"user1@example.com": {Email: "user1@example.com", Username: "user1"},
		"user2@example.com": {Email: "user2@example.com", Username: "user2"},
	}
	userRepo := mocks.NewMockUserRepository(mockUsers)
	friendRepo := mocks.NewMockFriendRepository(map[string]*models.Friend{
		"user2@example.com_user1@example.com": {
			Email:       "user2@example.com",
			FriendEmail: "user1@example.com",
			Status:      "pending",
		},
	})
	friendService := services.NewFriendService(userRepo, friendRepo)
	friendHandler := handlers.NewFriendHandler(friendService)

	requestData := map[string]string{
		"username": "user2",
	}
	body, _ := json.Marshal(requestData)
	req, err := http.NewRequest("POST", "/api/friends/decline", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	// Mock authentication context
	ctx := context.WithValue(req.Context(), "userEmail", "user1@example.com")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(friendHandler.DeclineFriendRequest)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Verify response body
	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response body")
	}
	if response["message"] != "Friend request declined" {
		t.Errorf("Unexpected response message: %s", response["message"])
	}

	// Verify that the friend request has been removed
	friendKey := "user2@example.com_user1@example.com"
	if _, exists := friendRepo.Friends[friendKey]; exists {
		t.Errorf("Friend request not removed from mock repository")
	}
}

func TestCancelFriendRequestHandler(t *testing.T) {
	mockUsers := map[string]*models.User{
		"user1@example.com": {Email: "user1@example.com", Username: "user1"},
		"user2@example.com": {Email: "user2@example.com", Username: "user2"},
	}
	userRepo := mocks.NewMockUserRepository(mockUsers)
	friendRepo := mocks.NewMockFriendRepository(map[string]*models.Friend{
		"user1@example.com_user2@example.com": {
			Email:       "user1@example.com",
			FriendEmail: "user2@example.com",
			Status:      "pending",
		},
	})
	friendService := services.NewFriendService(userRepo, friendRepo)
	friendHandler := handlers.NewFriendHandler(friendService)

	requestData := map[string]string{
		"username": "user2",
	}
	body, _ := json.Marshal(requestData)
	req, err := http.NewRequest("POST", "/api/friends/cancel", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	// Mock authentication context
	ctx := context.WithValue(req.Context(), "userEmail", "user1@example.com")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(friendHandler.CancelFriendRequest)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Verify response body
	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response body")
	}
	if response["message"] != "Friend request canceled" {
		t.Errorf("Unexpected response message: %s", response["message"])
	}

	// Verify that the friend request has been removed
	friendKey := "user1@example.com_user2@example.com"
	if _, exists := friendRepo.Friends[friendKey]; exists {
		t.Errorf("Friend request not removed from mock repository")
	}
}
