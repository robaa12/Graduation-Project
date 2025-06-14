package auth

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

// TokenResponse defines the structure for authentication token responses
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type Claims struct {
	UserID    int    `json:"user_id"`
	StoresID  []int  `json:"stores_id"`
	TokenType string `json:"token_type"`
	jwt.StandardClaims
}

type JWTService struct {
	secretKey          []byte
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

// GetAccessTokenExpiry returns the configured access token expiry duration
func (s *JWTService) GetAccessTokenExpiry() time.Duration {
	return s.accessTokenExpiry
}

func NewJWTService(secretKey string, accessExpiry, refreshExpiry time.Duration) *JWTService {
	return &JWTService{
		secretKey:          []byte(secretKey),
		accessTokenExpiry:  accessExpiry,
		refreshTokenExpiry: refreshExpiry,
	}
}

func (s *JWTService) GenerateTokenPair(userID int, storeID []int) (accessToken, refreshToken string, err error) {
	accessToken, err = s.GenerateToken(userID, storeID, "access", s.accessTokenExpiry)
	if err != nil {
		return "", "", err
	}
	refreshToken, err = s.GenerateToken(userID, storeID, "refresh", s.refreshTokenExpiry)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *JWTService) GenerateToken(userID int, storeID []int, tokenType string, expiry time.Duration) (string, error) {
	claims := &Claims{
		UserID:    userID,
		StoresID:  storeID,
		TokenType: tokenType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expiry).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return s.secretKey, nil // ✅ Explicitly return []byte
	})
	if err != nil {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("token expired")
	}

	return claims, nil
}

// GenerateUpdatedTokenResponse generates a new TokenResponse if the new store ID is not already present.
func (s *JWTService) GenerateUpdatedTokenResponse(userID int, currentStoreIDs []int, newStoreID uint) (*TokenResponse, error) {
	for _, sid := range currentStoreIDs {
		if sid == int(newStoreID) {
			return nil, nil // No update needed
		}
	}
	updatedStoreIDs := append(currentStoreIDs, int(newStoreID))
	accessToken, refreshToken, err := s.GenerateTokenPair(userID, updatedStoreIDs)
	if err != nil {
		return nil, err
	}
	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.GetAccessTokenExpiry().Seconds()),
	}, nil
}

func (s *JWTService) GetUserIDFromJWT(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		return "", errors.New("missing or invalid authorization header")
	}
	tokenString := authHeader[7:]

	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return "", errors.New("Invalid credentials.")
	}

	return strconv.Itoa(claims.UserID), nil
}
