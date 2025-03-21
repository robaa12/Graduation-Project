package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/robaa12/gatway-service/internal/config"
	"github.com/robaa12/gatway-service/utils"
)

type AuthService struct {
	jwtService  *JWTService
	userService config.ServiceConfig
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token    string `json:"token"`
	UserID   int    `json:"user_id"`
	StoresID []int  `json:"stores_id"`
	Email    string `json:"email"`
	Image    string `json:"image"`
	Name     string `json:"name"`
}

type AuthResponse struct {
	Token       string  `json:"token"`
	UserID      int     `json:"user_id"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	Email       string  `json:"email"`
	StoresID    []int   `json:"stores_id,omitempty"`
	IsActive    bool    `json:"is_active,omitempty"`
	IsBanned    bool    `json:"is_banned,omitempty"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	Address     *string `json:"address,omitempty"`
	Image       string  `json:"image,omitempty"`
	CreatedAt   string  `json:"created_at,omitempty"`
	UpdatedAt   string  `json:"updated_at,omitempty"`
}

type UserResponse struct {
	ID          int     `json:"id"`
	FirstName   string  `json:"firstName"`
	LastName    string  `json:"lastName"`
	IsActive    bool    `json:"isActive"`
	Email       string  `json:"email"`
	IsBanned    bool    `json:"is_banned"`
	PhoneNumber *string `json:"phoneNumber"`
	StoresID    []int   `json:"stores_id"`
	Address     *string `json:"address"`
	CreateAt    string  `json:"createAt"`
	UpdateAt    string  `json:"updateAt"`
}

type APIResponse struct {
	Status  bool         `json:"status"`
	Message string       `json:"message"`
	Data    UserResponse `json:"data"`
}

type RegisterAPIResponse struct {
	Status    bool   `json:"status"`
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	StoresID  []int  `json:"stores_id"`
}

func NewAuthService(cfg *config.Config) *AuthService {
	return &AuthService{
		jwtService:  NewJWTService(cfg.Auth.JWTSecret, cfg.Auth.TokenExp),
		userService: cfg.Services["user-service"],
	}
}

// Register handles the register request
func (s *AuthService) Register(w http.ResponseWriter, r *http.Request) {
	// Make request to user service
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/user", s.userService.URL), r.Body)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		log.Println(resp.StatusCode)
		utils.ErrorJSON(w, errors.New("error registering user"), http.StatusInternalServerError)
		return
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Log the raw response for debugging
	log.Printf("Registration API Response: %s", string(body))

	// Parse the response using the struct that matches the actual format
	var registerResp RegisterAPIResponse
	err = json.Unmarshal(body, &registerResp)
	if err != nil {
		log.Printf("Failed to unmarshal API response: %v", err)
		utils.ErrorJSON(w, errors.New("Failed to parse user data"), http.StatusInternalServerError)
		return
	}

	// Generate JWT Token
	token, err := s.jwtService.GenerateToken(registerResp.ID, registerResp.StoresID)
	if err != nil {
		utils.ErrorJSON(w, errors.New("Error Generating JWT Token"), http.StatusInternalServerError)
		return
	}

	// Create the response with all available fields
	utils.WriteJSON(w, http.StatusCreated, AuthResponse{
		Token:     token,
		UserID:    registerResp.ID,
		FirstName: registerResp.FirstName,
		LastName:  registerResp.LastName,
		Email:     registerResp.Email,
		StoresID:  registerResp.StoresID,
		// Other fields might not be available in the register response
		// but can be added if they are
		Image: "",
	})
}

// Login handles the login request
func (s *AuthService) Login(w http.ResponseWriter, r *http.Request) {
	var loginReq LoginRequest
	err := utils.ReadJSON(w, r, &loginReq)
	if err != nil {
		utils.ErrorJSON(w, errors.New("invalid request body"), http.StatusBadRequest)
		return
	}

	// Call user service to authenticate user
	user, err := s.validateCredentials(loginReq)
	if err != nil {
		log.Printf("err:%s", err)
		utils.ErrorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := s.jwtService.GenerateToken(user.ID, user.StoresID)
	if err != nil {
		utils.ErrorJSON(w, errors.New("error generating token"), http.StatusInternalServerError)
		return
	}

	// Send response with all available user data
	utils.WriteJSON(w, http.StatusOK, AuthResponse{
		Token:       token,
		UserID:      user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		StoresID:    user.StoresID,
		IsActive:    user.IsActive,
		IsBanned:    user.IsBanned,
		PhoneNumber: user.PhoneNumber,
		Address:     user.Address,
		Image:       "", // Default empty image
		CreatedAt:   user.CreateAt,
		UpdatedAt:   user.UpdateAt,
	})
}

// validateCredentials validates the user credentials
func (s *AuthService) validateCredentials(login LoginRequest) (*UserResponse, error) {
	// Call user service to validate user
	reqBody, err := json.Marshal(login)
	if err != nil {
		return nil, err
	}

	// Make request to user service
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/user/login", s.userService.URL), strings.NewReader(string(reqBody)))
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
		return nil, errors.New("invalid credentials")
	}

	// Read and log the raw response for debugging
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("Login API Response: %s", string(body))

	// Parse the response
	var apiResponse APIResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, err
	}

	return &apiResponse.Data, nil
}

func (s *AuthService) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Applying middleware: AuthMiddleware")
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

func (s *AuthService) StoreOwnershipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from user context
		claims, ok := r.Context().Value("user").(*Claims)

		if !ok {
			utils.ErrorJSON(w, errors.New("Invalid Store owner"), http.StatusUnauthorized)
			return
		}

		// Get store ID from URL
		storeID, err := utils.GetID(r, "store_id")
		if err != nil {
			utils.ErrorJSON(w, errors.New("Store ID is required"), http.StatusBadRequest)
			return
		}

		// Check if store ID is in user's stores
		// Loop through user's stores array and check if storeID is in them
		for _, id := range claims.StoresID {
			if id == storeID {
				next.ServeHTTP(w, r)
				return
			}
		}
		// If store ID is not in user's stores
		utils.ErrorJSON(w, errors.New("Unauthorized Owner"), http.StatusUnauthorized)

	})
}
