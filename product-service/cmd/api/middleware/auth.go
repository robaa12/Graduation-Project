package middleware

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/robaa12/product-service/cmd/utils"
)

type Claims struct {
	UserID  uint   `json:"user_id"`
	StoreID uint   `json:"store_id"`
	Role    string `json:"role"`
	jwt.StandardClaims
}

func AuthenticateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.ErrorJSON(w, errors.New("authorization header is required"), http.StatusUnauthorized)
			return
		}

		// Check if the header starts with "Bearer "
		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			utils.ErrorJSON(w, errors.New("invalid authorization header format"), http.StatusUnauthorized)
			return
		}

		tokenString := bearerToken[1]

		// Parse and validate the token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Replace this with your actual JWT secret key from environment variables
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			utils.ErrorJSON(w, errors.New("invalid or expired token"), http.StatusUnauthorized)
			return
		}

		// Add claims to request context
		ctx := context.WithValue(r.Context(), "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func VerifyStoreOwnership(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the store ID from the URL
		storeID, err := utils.GetID(r, "store_id")
		if err != nil {
			utils.ErrorJSON(w, errors.New("missing store ID parameter"))
			return
		}
		// Retrieve the claims from the context
		claims, ok := r.Context().Value("claims").(*Claims)
		if !ok || claims == nil {
			utils.ErrorJSON(w, errors.New("unauthenticated or invalid token claims"), http.StatusUnauthorized)
			return
		}

		// Check if the store ID in the token matches the store ID in the URL
		if claims.StoreID != storeID {
			utils.ErrorJSON(w, errors.New("forebidden: you do not own this store"), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
