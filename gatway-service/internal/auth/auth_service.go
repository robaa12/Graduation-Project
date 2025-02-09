package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/robaa12/gatway-service/config"
	"github.com/robaa12/gatway-service/internal/utils"
)

type AuthService struct {
	jwtService     *JWTService
	userServiceURL string
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	UserID       int    `json:"user_id"`
	StoresID     []int  `json:"stores_id"`
}

type UserResponse struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	StoresID  []int  `json:"stores_id"`
}

type APIResponse struct {
	Status  bool         `json:"status"`
	Message string       `json:"message"`
	Data    UserResponse `json:"data"`
}

func NewAuthService(cfg *config.Config) *AuthService {
	return &AuthService{
		jwtService:     NewJWTService(cfg.Auth.JWTSecret, cfg.Auth.TokenExp),
		userServiceURL: cfg.Services.UserService.URL,
	}
}

// Login handles the login request
func (s *AuthService) Login(w http.ResponseWriter, r *http.Request) {
	var loginReq LoginRequest
	err := utils.ReadJSON(w, r, &loginReq)
	if err != nil {
		utils.ErrorJSON(w, errors.New("Invalid request body"), http.StatusBadRequest)
		return
	}

	// Call user service to authenticate user
	user, err := s.validateCredentials(loginReq)
	if err != nil {
		log.Printf("err:%s", err)
		utils.ErrorJSON(w, errors.New("Invalid credentials"), http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := s.jwtService.GenerateToken(user.ID, user.StoresID)
	if err != nil {
		utils.ErrorJSON(w, errors.New("Error generating token"), http.StatusInternalServerError)
		return
	}

	// Send response
	utils.WriteJSON(w, http.StatusOK, LoginResponse{
		Token:    token,
		UserID:   user.ID,
		StoresID: user.StoresID,
	})
}

// validateCredentials validates the user credentials
func (s *AuthService) validateCredentials(login LoginRequest) (*UserResponse, error) {
	// Call user service to validate user
	reqBody, err := json.Marshal(login)
	if err != nil {
		return nil, err
	}

	// Make request to user service (Solve this when we know the endpoint path)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/user/login", s.userServiceURL), strings.NewReader(string(reqBody)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Invalid credentials")
	}
	// Store the response data in a struct

	var apiResponse APIResponse

	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		return nil, err
	}

	user := apiResponse.Data
	return &user, nil
}

func (s *AuthService) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.ErrorJSON(w, errors.New("Authorization header is required"), http.StatusUnauthorized)
			return
		}

		// Remove "Bearer " prefix
		token := strings.Replace(authHeader, "Bearer ", "", 1)

		// Validate token
		claims, err := s.jwtService.ValidateToken(token)
		if err != nil {
			utils.ErrorJSON(w, errors.New("Invalid token"), http.StatusUnauthorized)
			return
		}

		// Add claims to context
		ctx := context.WithValue(r.Context(), "user", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
