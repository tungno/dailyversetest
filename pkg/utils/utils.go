// pkg/utils/utils.go
package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"time"
	"unicode"

	"github.com/dgrijalva/jwt-go"
	"math/rand"
)

// JWT Secret Key from environment variable
var jwtSecretKey = os.Getenv("JWT_SECRET_KEY")

// GenerateJWT generates a JWT token for the user
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

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// HashPassword hashes the password using SHA-256
func HashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

// IsValidPassword checks password complexity
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

// GenerateOTP generates a 6-digit OTP
func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	return randSeq(6)
}

var letters = []rune("1234567890")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// WriteJSON writes a JSON response
func WriteJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// WriteJSONError writes a JSON error response
func WriteJSONError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{
		"message": message,
	})
}
