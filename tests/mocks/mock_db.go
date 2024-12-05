/**
 *  MockDB simulates a database for testing purposes.
 *  This mock implementation can be used to test handlers, services, and other components
 *  without relying on an actual Firestore database.
 *
 *  @file       mock_db.go
 *  @package    mocks
 *
 *  @structs
 *  - MockDB: Simulates a database with collections for Users and Friends.
 *  - MockCollection: Simulates a Firestore collection for operations like Get, Add, Update, and Delete.
 *
 *  @methods
 *  - Collection(path): Retrieves a mock collection based on the path.
 *  - Close(): Simulates closing the database connection.
 *  - GetUserByEmail(email): Retrieves a user by email.
 *  - GetUserByUsername(username): Retrieves a user by username.
 *  - AddFriendRequest(friend): Adds a friend request.
 *  - UpdateFriendStatus(docID, status): Updates the status of a friend request.
 *  - DeleteFriend(docID): Deletes a friend relationship or request.
 *
 *  @example
 *  ```
 *  mdb := &mocks.MockDB{
 *      Users: map[string]*models.User{
 *          "test@example.com": {
 *              Email:    "test@example.com",
 *              Username: "testuser",
 *          },
 *      },
 *      Friends: make(map[string]*models.Friend),
 *  }
 *
 *  user, err := mdb.GetUserByEmail("test@example.com")
 *  if err != nil {
 *      t.Fatalf("Failed to fetch user: %v", err)
 *  }
 *  ```
 *
 *  @dependencies
 *  - pkg/models: Provides the User and Friend models for the mock database.
 *
 *  @limitations
 *  - MockDB is in-memory and does not persist data across tests.
 *  - MockDB does not implement all Firestore features, only those required for testing.
 *
 *  @authors
 *      - Aayush
 *      - Tung
 *      - Boss
 *      - Majd
 */

package mocks

import (
	"errors"
	"proh2052-group6/pkg/models"
)

// MockDB simulates a database for testing purposes.
type MockDB struct {
	Users   map[string]*models.User   // Simulated users collection.
	Friends map[string]*models.Friend // Simulated friends collection.
}

// Collection simulates retrieving a Firestore collection.
func (mdb *MockDB) Collection(path string) *MockCollection {
	return &MockCollection{
		Path:   path,
		MockDB: mdb,
	}
}

// Close simulates closing the database connection (no-op).
func (mdb *MockDB) Close() error {
	return nil
}

// MockCollection simulates a Firestore collection.
type MockCollection struct {
	Path   string
	MockDB *MockDB
}

// GetUserByEmail simulates retrieving a user by email.
func (mdb *MockDB) GetUserByEmail(email string) (*models.User, error) {
	user, exists := mdb.Users[email]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// GetUserByUsername simulates retrieving a user by username.
func (mdb *MockDB) GetUserByUsername(username string) (*models.User, error) {
	for _, user := range mdb.Users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

// AddFriendRequest simulates adding a friend request.
func (mdb *MockDB) AddFriendRequest(friend *models.Friend) error {
	docID := friend.Email + "_" + friend.FriendEmail
	mdb.Friends[docID] = friend
	return nil
}

// UpdateFriendStatus simulates updating the status of a friend request.
func (mdb *MockDB) UpdateFriendStatus(docID string, status string) error {
	friend, exists := mdb.Friends[docID]
	if !exists {
		return errors.New("friend request not found")
	}
	friend.Status = status
	return nil
}

// DeleteFriend simulates deleting a friend relationship or request.
func (mdb *MockDB) DeleteFriend(docID string) error {
	delete(mdb.Friends, docID)
	return nil
}
