// pkg/models/models.go
package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// User model for signup and login
type User struct {
	Username      string    `json:"username"`
	UsernameLower string    `json:"usernameLower"`
	Email         string    `json:"email"`
	Password      string    `json:"password"`
	Country       string    `json:"country"`
	City          string    `json:"city"`
	ImageURL      string    `json:"imageUrl,omitempty"`
	FirstName     string    `json:"firstName,omitempty"`
	LastName      string    `json:"lastName,omitempty"`
	IsVerified    bool      `json:"isVerified"`
	OTP           string    `json:"-"`
	OTPExpiresAt  time.Time `json:"-"`
}

// LoginRequest model
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Event model representing event details
type Event struct {
	EventID       string `json:"eventID"`
	StreetAddress string `json:"streetAddress"`
	PostalNumber  string `json:"postalNumber"`
	Status        string `json:"status"`
	Description   string `json:"description"`
	Time          string `json:"time"`
	EventTypeID   string `json:"eventTypeID"`
	Date          string `json:"date"`
	Email         string `json:"email"` // User's email as foreign key
	Title         string `json:"title"`
	StartTime     string `json:"startTime"`
	EndTime       string `json:"endTime"`
}

// Journal model for daily journal entries
type Journal struct {
	JournalID string `json:"journalID,omitempty"`
	Date      string `json:"date"`
	Content   string `json:"content"`
	Email     string `json:"email"` // User's email as foreign key
}

// Friend model to manage friendships between users
type Friend struct {
	Email       string `json:"email"`       // Email of the user who sent the request
	FriendEmail string `json:"friendEmail"` // Email of the user who received the request
	Status      string `json:"status"`      // "pending" or "accepted"
}

// JWT Claims structure
type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}
