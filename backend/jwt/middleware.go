package jwt

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// No token, pass the request as unauthenticated
			next.ServeHTTP(w, r)
			return
		}

		// Split "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Malformed token", http.StatusUnauthorized)
			return
		}

		// Parse and validate the token
		token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		// If there is an error or the token is invalid, deny access
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		// Extract the claims and set user context
		if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
			ctx := context.WithValue(r.Context(), "userID", claims.UserID)
			r = r.WithContext(ctx)
		} else {
			http.Error(w, "Unauthorized: Invalid claims", http.StatusUnauthorized)
			return
		}

		// Continue with the request
		next.ServeHTTP(w, r)
	})
}
