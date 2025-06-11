package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/robaa12/gatway-service/internal/config"
	"github.com/robaa12/gatway-service/internal/middleware/auth"
	"github.com/robaa12/gatway-service/utils"
)

// UserHandler handles user-related operations
type UserHandler struct {
	cfg        *config.Config
	jwtService *auth.JWTService
}

// NewUserHandler creates a new user handler
func NewUserHandler(cfg *config.Config, jwtService *auth.JWTService) *UserHandler {
	return &UserHandler{
		cfg:        cfg,
		jwtService: jwtService,
	}
}

// GetUser handles requests to get user information
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	// Get user ID from JWT token
	userID, err := h.jwtService.GetUserIDFromJWT(r)
	if err != nil {
		utils.ErrorJSON(w, fmt.Errorf("unauthorized: %v", err), http.StatusUnauthorized)
		return
	}

	// Get user information from user-service
	userInfo, err := h.getUserInfo(userID)
	if err != nil {
		utils.ErrorJSON(w, fmt.Errorf("error fetching user data: %v", err), http.StatusInternalServerError)
		return
	}

	// Send user information back to client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(userInfo)
}

// getUserInfo retrieves user information from the user-service
func (h *UserHandler) getUserInfo(userID string) ([]byte, error) {
	// Get user-service URL from configuration
	userServiceURL := h.cfg.Services["user-service"].URL

	// Build the request URL
	requestURL := fmt.Sprintf("%s/user/%s", userServiceURL, userID)

	// Create and execute the request
	client := &http.Client{Timeout: h.cfg.Services["user-service"].Timeout}
	resp, err := client.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user-service: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-200 responses
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user-service returned status: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response from user-service: %w", err)
	}

	return body, nil
}
