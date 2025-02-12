package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	UserID   int   `json:"user_id"`
	StoresID []int `json:"stores_id"`
	jwt.StandardClaims
}

type JWTService struct {
	secretKey []byte
	expiry    time.Duration
}

func NewJWTService(secretKey string, expiry time.Duration) *JWTService {
	return &JWTService{
		secretKey: []byte(secretKey),
		expiry:    expiry,
	}
}

func (s *JWTService) GenerateToken(userID int, storeID []int) (string, error) {
	claims := &Claims{
		UserID:   userID,
		StoresID: storeID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.expiry).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil // ✅ Explicitly return []byte
	})
	if err != nil {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid { // ✅ Check token validity before type assertion
		return nil, errors.New("invalid token")
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("token expired")
	}

	return claims, nil
}
