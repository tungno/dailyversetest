// internal/middleware/auth.go
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

var jwtSecretKey = os.Getenv("JWT_SECRET_KEY")

func JwtAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.WriteJSONError(w, "Authorization token is missing", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			utils.WriteJSONError(w, "Authorization token format must be 'Bearer <token>'", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		claims := &models.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecretKey), nil
		})

		if err != nil || !token.Valid {
			utils.WriteJSONError(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Pass the user's email to the next handler using context
		ctx := context.WithValue(r.Context(), "userEmail", claims.Email)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
