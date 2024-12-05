/**
 *  Models package defines the data structures for the application.
 *  These models are used across various layers of the application,
 *  including repositories, services, and handlers.
 *
 *  @file       models.go
 *  @package    models
 *
 *  @structs
 *  - User: Represents a user account with details like username, email, and password.
 *  - LoginRequest: Represents the request payload for user login.
 *  - Event: Represents event details for user-created events.
 *  - Journal: Represents a daily journal entry linked to a user.
 *  - Friend: Manages friendships or friend requests between users.
 *  - Claims: Represents JWT claims for authentication.
 *  - TimetableEvent: Represents events retrieved from the NTNU timetable API.
 *  - UserSummary: Provides minimal user information for frontend display.
 *
 *  @dependencies
 *  - github.com/dgrijalva/jwt-go: For handling JWT authentication claims.
 *
 *  @example
 *  ```
 *  user := models.User{
 *      Username: "JohnDoe",
 *      Email: "john@example.com",
 *      Country: "Norway",
 *      City: "Oslo",
 *  }
 *  ```
 *
 *  @authors
 *      - Aayush
 *      - Tung
 *      - Boss
 *      - Majd
 */

package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// User represents a user account with profile and authentication details.
type User struct {
	Username      string    `json:"username"`
	UsernameLower string    `json:"usernameLower"` // Lowercase version of the username for case-insensitive operations.
	Email         string    `json:"email"`
	Password      string    `json:"-"` // Stored as a hashed password.
	Country       string    `json:"country"`
	City          string    `json:"city"`
	ImageURL      string    `json:"imageUrl,omitempty"`
	FirstName     string    `json:"firstName,omitempty"`
	LastName      string    `json:"lastName,omitempty"`
	IsVerified    bool      `json:"isVerified"`
	OTP           string    `json:"-"` // One-Time Password for verification.
	OTPExpiresAt  time.Time `json:"-"` // Expiration time for the OTP.
}

// LoginRequest represents the payload for user login requests.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Event represents event details for user-created or imported events.
type Event struct {
	EventID       string `json:"eventID"`
	StreetAddress string `json:"streetAddress"`
	PostalNumber  string `json:"postalNumber"`
	Status        string `json:"status"`
	Description   string `json:"description"`
	Time          string `json:"time"`
	EventTypeID   string `json:"eventTypeID"`
	Date          string `json:"date"`
	Email         string `json:"email"` // User's email as a foreign key.
	Title         string `json:"title"`
	StartTime     string `json:"startTime"`
	EndTime       string `json:"endTime"`
}

// Journal represents a daily journal entry linked to a user.
type Journal struct {
	JournalID string `json:"journalID,omitempty"`
	Date      string `json:"date"`
	Content   string `json:"content"`
	Email     string `json:"email"` // User's email as a foreign key.
}

// Friend manages friendships or friend requests between users.
type Friend struct {
	Email       string `json:"email"`       // Email of the user who sent the request.
	FriendEmail string `json:"friendEmail"` // Email of the user who received the request.
	Status      string `json:"status"`      // "pending" or "accepted".
}

// Claims represents JWT claims for authentication and user identification.
type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// TimetableEvent represents the structure of events received from the NTNU timetable API.
type TimetableEvent struct {
	CourseCode  string `json:"courseCode"`
	CourseName  string `json:"courseName"`
	Description string `json:"description"`
	Date        string `json:"date"`      // Format: "YYYY-MM-DD".
	StartTime   string `json:"startTime"` // Format: "HH:MM".
	EndTime     string `json:"endTime"`   // Format: "HH:MM".
}

// UserSummary provides minimal user information for frontend display.
type UserSummary struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Country  string `json:"country"`
	City     string `json:"city"`
}
