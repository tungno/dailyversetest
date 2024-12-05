/**
 *  Utils Package provides utility functions for common operations such as password hashing,
 *  JSON response handling, JWT generation, email validation, and OTP generation. These
 *  functions serve as reusable components throughout the project.
 *
 *  @file      utils.go
 *  @package   utils
 *  @purpose   Utility functions for authentication, validation, and response handling.
 *
 *  @methods
 *  - GenerateJWT(email)                   - Generates a JWT token for the given email.
 *  - HashPassword(password)               - Hashes a password using SHA-256.
 *  - IsValidPassword(password)            - Validates password complexity requirements.
 *  - GenerateOTP()                        - Generates a random 6-digit OTP.
 *  - WriteJSON(w, data)                   - Writes a JSON response to the HTTP response writer.
 *  - WriteJSONError(w, message, code)     - Writes an error message as a JSON response.
 *  - CheckPasswordHash(password, hash)    - Compares a plain password with its hashed version.
 *  - IsValidEmail(email)                  - Validates if a string is a properly formatted email.
 *
 *  @dependencies
 *  - golang.org/x/crypto/bcrypt: Used for secure password hashing and comparison.
 *  - github.com/dgrijalva/jwt-go: Used for generating and validating JWT tokens.
 *  - crypto/sha256: Provides hashing capabilities.
 *
 *  @example
 *  ```
 *  hashedPassword := HashPassword("Secure@123")
 *  isValid := IsValidPassword("Secure@123")
 *  ```
 *
 *  @environment_variables
 *  - JWT_SECRET_KEY: Secret key used for signing JWT tokens.
 *
 *  @authors
 *      - Aayush
 *      - Tung
 *      - Boss
 *      - Majd
 */

package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"regexp"
	"time"
	"unicode"

	"github.com/dgrijalva/jwt-go"
	"math/rand"
)

// JWT Secret Key from environment variables
var jwtSecretKey = os.Getenv("JWT_SECRET_KEY")

// Claims defines the JWT token structure.
type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// GenerateJWT generates a JWT token for a given email.
// Parameters:
//   - email: The email address to associate with the token.
//
// Returns:
//   - string: A signed JWT token.
//   - error: Returns an error if token signing fails.
func GenerateJWT(email string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecretKey))
}

// HashPassword hashes a given password using SHA-256.
// Parameters:
//   - password: The plain text password to hash.
//
// Returns:
//   - string: The hashed password as a hexadecimal string.
func HashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

// IsValidPassword checks if a password meets complexity requirements.
// Requirements:
//   - At least 8 characters.
//   - Contains an uppercase letter, a number, and a special character.
//
// Parameters:
//   - password: The password to validate.
//
// Returns:
//   - bool: True if the password meets the requirements, false otherwise.
func IsValidPassword(password string) bool {
	var hasMinLen, hasUpper, hasNumber, hasSpecial bool
	if len(password) >= 8 {
		hasMinLen = true
	}
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasNumber && hasSpecial
}

// GenerateOTP generates a random 6-digit OTP.
// Returns:
//   - string: A 6-digit OTP as a string.
func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	return randSeq(6)
}

var letters = []rune("1234567890")

// randSeq generates a random string of n digits.
// Parameters:
//   - n: The length of the random string.
//
// Returns:
//   - string: A random string of digits.
func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// WriteJSON writes a JSON response to the HTTP response writer.
// Parameters:
//   - w: The HTTP response writer.
//   - data: The data to encode as JSON.
func WriteJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// WriteJSONError writes an error message as a JSON response with a specific status code.
// Parameters:
//   - w: The HTTP response writer.
//   - message: The error message.
//   - code: The HTTP status code.
func WriteJSONError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{
		"message": message,
	})
}

// CheckPasswordHash compares a plain password with a hashed password.
// Parameters:
//   - password: The plain text password.
//   - hash: The hashed password to compare.
//
// Returns:
//   - bool: True if the passwords match, false otherwise.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// IsValidEmail validates if a string is a properly formatted email address.
// Parameters:
//   - email: The email address to validate.
//
// Returns:
//   - bool: True if the email is valid, false otherwise.
func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*$`)
	return emailRegex.MatchString(email)
}
