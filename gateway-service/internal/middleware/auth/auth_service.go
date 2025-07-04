package auth

import (
	"bytes"
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

type Service struct {
	jwtService  *JWTService
	userService config.ServiceConfig
	client      *http.Client
}

type (
	LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	Store struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	UserData struct {
		ID          int     `json:"id"`
		FirstName   string  `json:"firstName"`
		LastName    string  `json:"lastName"`
		Email       string  `json:"email"`
		IsActive    bool    `json:"isActive"`
		IsBanned    bool    `json:"is_banned"`
		PhoneNumber *string `json:"phoneNumber,omitempty"`
		Address     *string `json:"address,omitempty"`
		CreateAt    string  `json:"createAt,omitempty"`
		UpdateAt    string  `json:"updateAt,omitempty"`
		Stores      []Store `json:"stores"`
	}

	LoginResponse struct {
		UserID int     `json:"user_id"`
		Stores []Store `json:"stores"`
		Email  string  `json:"email"`
		Name   string  `json:"name"`
		TokenResponse
	}

	RefreshTokenRequest struct {
		RefreshToken string `json:"refresh_token"`
	}

	APIResponse struct {
		ID        int     `json:"id"`
		Email     string  `json:"email"`
		FirstName string  `json:"first_name"`
		LastName  string  `json:"last_name"`
		Stores    []Store `json:"stores"`
	}
	LoginAPIResponse struct {
		Message string   `json:"message"`
		Status  bool     `json:"status"`
		Data    UserData `json:"data"`
	}
)

func NewAuthService(cfg *config.Config) *Service {
	return &Service{
		jwtService:  NewJWTService(cfg.Auth.JWTSecret, cfg.Auth.AccessTokenExp, cfg.Auth.RefreshTokenExp),
		userService: cfg.Services["user-service"],
		client:      &http.Client{Timeout: 5 * time.Second},
	}
}

// Login handles user authentication
func (s *Service) Login(w http.ResponseWriter, r *http.Request) {
	var loginReq LoginRequest
	err := utils.ReadJSON(w, r, &loginReq)
	if err != nil {
		_ = utils.ErrorJSON(w, errors.New("invalid request body"), http.StatusBadRequest)
		return
	}

	userData, err := s.authenticateUser(loginReq)
	if err != nil {
		_ = utils.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	response, err := s.generateLoginResponse(userData)
	if err != nil {
		_ = utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	// Send the response back to the client
	_ = utils.WriteJSON(w, http.StatusOK, response)
}

// Register handles user registration
func (s *Service) Register(w http.ResponseWriter, r *http.Request) {
	var registerReq struct {
		Email       string  `json:"email"`
		Password    string  `json:"password"`
		FirstName   string  `json:"firstName"`
		LastName    string  `json:"lastName"`
		PhoneNumber *string `json:"phoneNumber,omitempty"`
		Address     *string `json:"address,omitempty"`
	}

	if err := utils.ReadJSON(w, r, &registerReq); err != nil {
		_ = utils.ErrorJSON(w, fmt.Errorf("invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Log the incoming request
	log.Printf("Registration request received: %+v", registerReq)

	// Basic validation
	if registerReq.Email == "" || registerReq.Password == "" || registerReq.FirstName == "" || registerReq.LastName == "" {
		_ = utils.ErrorJSON(w, errors.New("email, password, firstName, and lastName are required"), http.StatusBadRequest)
		return
	}

	// Convert request to JSON
	jsonData, err := json.Marshal(registerReq)
	if err != nil {
		_ = utils.ErrorJSON(w, fmt.Errorf("error preparing request: %v", err), http.StatusInternalServerError)
		return
	}

	// Log the outgoing request to user service
	log.Printf("Sending to user service: %s", string(jsonData))

	// Register the user
	userData, err := s.registerUser(bytes.NewReader(jsonData))
	if err != nil {
		log.Printf("Registration error: %v", err)
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "already exists") {
			statusCode = http.StatusConflict
		} else if strings.Contains(err.Error(), "status 400") {
			statusCode = http.StatusBadRequest
		}
		_ = utils.ErrorJSON(w, err, statusCode)
		return
	}

	// Log successful registration
	log.Printf("User registered successfully: %+v", userData)

	// Generate tokens
	response, err := s.generateLoginResponse(userData)
	if err != nil {
		_ = utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Send the response
	_ = utils.WriteJSON(w, http.StatusCreated, response)
}

// RefreshToken handles token refresh requests
func (s *Service) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req RefreshTokenRequest
	if err := utils.ReadJSON(w, r, &req); err != nil {
		_ = utils.ErrorJSON(w, errors.New("invalid request body"), http.StatusBadRequest)
		return
	}

	claims, err := s.validateRefreshToken(req.RefreshToken)
	if err != nil {
		_ = utils.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	tokenResponse, err := s.generateNewTokenPair(claims.UserID, claims.StoresID)
	if err != nil {
		_ = utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	_ = utils.WriteJSON(w, http.StatusOK, tokenResponse)
}

// validateCredentials validates the user credentials
func (s *Service) authenticateUser(login LoginRequest) (*UserData, error) {
	reqBody, err := json.Marshal(login)
	if err != nil {
		return nil, err
	}

	resp, err := s.makeUserServiceRequest("POST", "/user/login", reqBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}
	log.Printf("Raw login response: %s", string(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("invalid credentials")
	}

	var loginResp LoginAPIResponse
	if err := json.Unmarshal(bodyBytes, &loginResp); err != nil {
		return nil, fmt.Errorf("error parsing login response: %v", err)
	}
	userData := loginResp.Data

	if userData.ID == 0 {
		return nil, fmt.Errorf("invalid user data received from login: %+v", userData)
	}

	return &userData, nil
}

func (s *Service) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", 1)
		if token == "" {
			_ = utils.ErrorJSON(w, errors.New("authorization header required"), http.StatusUnauthorized)
			return
		}

		claims, err := s.jwtService.ValidateToken(token)
		if err != nil {
			_ = utils.ErrorJSON(w, err, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Service) StoreOwnershipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value("user").(*Claims)
		if !ok {
			_ = utils.ErrorJSON(w, errors.New("invalid store owner"), http.StatusUnauthorized)
			return
		}

		storeID, err := utils.GetID(r, "store_id")
		if err != nil {
			_ = utils.ErrorJSON(w, err, http.StatusBadRequest)
			return
		}

		if !contains(claims.StoresID, storeID) {
			_ = utils.ErrorJSON(w, errors.New("unauthorized owner"), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Service) makeUserServiceRequest(method, path string, body any) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		switch v := body.(type) {
		case io.Reader:
			reqBody = v
		case []byte:
			reqBody = strings.NewReader(string(v))
		default:
			return nil, errors.New("invalid request body")
		}
	}
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", s.userService.URL, path), reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return s.client.Do(req)
}

// generateLoginResponse generates the login response including tokens
func (s *Service) generateLoginResponse(userData *UserData) (*LoginResponse, error) {
	if userData == nil {
		return nil, errors.New("user data is nil")
	}

	name := strings.TrimSpace(fmt.Sprintf("%s %s", userData.FirstName, userData.LastName))
	if name == "" {
		name = "Unknown"
	}
	// Get All stores ID and store them in array
	stores := []int{}
	for _, store := range userData.Stores {
		stores = append(stores, store.ID)
	}
	accessToken, refreshToken, err := s.jwtService.GenerateTokenPair(userData.ID, stores)
	if err != nil {
		return nil, err
	}

	response := &LoginResponse{
		TokenResponse: TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    int64(s.jwtService.GetAccessTokenExpiry().Seconds()),
		},
		UserID: userData.ID,
		Stores: userData.Stores,
		Email:  userData.Email,
		Name:   name,
	}

	return response, nil
}
func (s *Service) registerUser(body io.Reader) (*UserData, error) {
	resp, err := s.makeUserServiceRequest("POST", "/user", body)
	if err != nil {
		return nil, fmt.Errorf("request error: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	// Log the raw response for debugging
	log.Printf("Raw response from user service: %s", string(bodyBytes))

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("error registering user: status %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var apiResponse APIResponse
	if err := json.Unmarshal(bodyBytes, &apiResponse); err != nil {
		return nil, fmt.Errorf("error parsing response: %v, body: %s", err, string(bodyBytes))
	}

	// Log the parsed response
	log.Printf("Parsed API response: %+v", apiResponse)

	// Convert the flat APIResponse to UserData
	userData := &UserData{
		ID:        apiResponse.ID,
		FirstName: apiResponse.FirstName,
		LastName:  apiResponse.LastName,
		Email:     apiResponse.Email,
		Stores:    apiResponse.Stores,
		IsActive:  true, // Default for new users
	}

	// Validate the user data
	if userData.ID == 0 {
		return nil, fmt.Errorf("invalid user data received: %+v, raw response: %s", userData, string(bodyBytes))
	}

	return userData, nil
}

func (s *Service) validateRefreshToken(token string) (*Claims, error) {
	claims, err := s.jwtService.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != "refresh" {
		return nil, errors.New("invalid token type")
	}

	return claims, nil
}

func (s *Service) generateNewTokenPair(userID int, storesID []int) (*TokenResponse, error) {
	accessToken, refreshToken, err := s.jwtService.GenerateTokenPair(userID, storesID)
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.jwtService.GetAccessTokenExpiry().Seconds()),
	}, nil
}

func contains(slice []int, item int) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
