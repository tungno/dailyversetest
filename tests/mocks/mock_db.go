// tests/mocks/mock_db.go
package mocks

import (
	"errors"
	"proh2052-group6/pkg/models"
)

type MockDB struct {
	Users   map[string]*models.User
	Friends map[string]*models.Friend
}

func (mdb *MockDB) Collection(path string) *MockCollection {
	return &MockCollection{
		Path:   path,
		MockDB: mdb,
	}
}
func (mdb *MockDB) Close() error {
	// Mock Close method; no action needed
	return nil
}

// MockCollection simulates a Firestore Collection
type MockCollection struct {
	Path   string
	MockDB *MockDB
}

// Implement methods to simulate database operations
func (mdb *MockDB) GetUserByEmail(email string) (*models.User, error) {
	user, exists := mdb.Users[email]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (mdb *MockDB) GetUserByUsername(username string) (*models.User, error) {
	for _, user := range mdb.Users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (mdb *MockDB) AddFriendRequest(friend *models.Friend) error {
	docID := friend.Email + "_" + friend.FriendEmail
	mdb.Friends[docID] = friend
	return nil
}

func (mdb *MockDB) UpdateFriendStatus(docID string, status string) error {
	friend, exists := mdb.Friends[docID]
	if !exists {
		return errors.New("friend request not found")
	}
	friend.Status = status
	return nil
}

func (mdb *MockDB) DeleteFriend(docID string) error {
	delete(mdb.Friends, docID)
	return nil
}

// Add more methods as needed for testing
