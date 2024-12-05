/**
 *  JwtAuthMiddleware is a middleware function that validates JWT tokens for secure API endpoints.
 *  It ensures that only authenticated users can access protected resources by verifying the token
 *  provided in the "Authorization" header of incoming HTTP requests.
 *
 *  @middleware JwtAuthMiddleware
 *
 *  @behaviors
 *  - Verifies the presence and format of the Authorization header.
 *  - Parses and validates the JWT token using the secret key.
 *  - Extracts the user's email from the token claims and attaches it to the request context.
 *  - Returns a 401 Unauthorized status for invalid or missing tokens.
 *
 *  @dependencies
 *  - jwt-go: Library for working with JSON Web Tokens.
 *  - models.Claims: Struct defining the claims within the JWT token.
 *  - utils: Utility package for writing JSON responses and errors.
 *  - os.Getenv("JWT_SECRET_KEY"): Environment variable storing the JWT secret key.
 *
 *  @example
 *  ```
 *  Authorization: Bearer <valid_jwt_token>
 *
 *  Valid Request:
 *  - Header: Authorization: Bearer <jwt_token>
 *  - Claims: { "email": "user@example.com", ... }
 *  - Next handler receives the user's email in the request context.
 *
 *  Invalid Request:
 *  - Header: Authorization: Bearer <invalid_jwt_token>
 *  - Response: { "error": "Invalid or expired token" }
 *  ```
 *
 *  @file      auth.go
 *  @project   DailyVerse
 *  @framework Go HTTP Server
 *  @authors
 *      - Aayush
 *      - Tung
 *      - Boss
 *      - Majd
 */

package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"proh2052-group6/pkg/models"
	"proh2052-group6/pkg/utils"

	"github.com/dgrijalva/jwt-go"
)

// jwtSecretKey holds the JWT secret key from the environment variable.
var jwtSecretKey = os.Getenv("JWT_SECRET_KEY")

// JwtAuthMiddleware is a middleware for validating JWT tokens in incoming requests.
// It ensures that only authenticated users can access the next handler.
func JwtAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the Authorization header from the incoming request.
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.WriteJSONError(w, "Authorization token is missing", http.StatusUnauthorized)
			return
		}

		// Ensure the token format is "Bearer <token>".
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			utils.WriteJSONError(w, "Authorization token format must be 'Bearer <token>'", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		claims := &models.Claims{}

		// Parse and validate the JWT token using the secret key.
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecretKey), nil
		})

		// Handle invalid or expired tokens.
		if err != nil || !token.Valid {
			utils.WriteJSONError(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Attach the user's email to the request context.
		ctx := context.WithValue(r.Context(), "userEmail", claims.Email)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
