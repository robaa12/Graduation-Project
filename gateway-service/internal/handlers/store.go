// File: Graduation-Project/gateway-service/internal/handlers/store/store_handler.go
package store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/robaa12/gatway-service/internal/config"
	"github.com/robaa12/gatway-service/internal/middleware/auth"
	"github.com/robaa12/gatway-service/utils"
)

type Handler struct {
	userService config.ServiceConfig
	client      *http.Client
	cfg         *config.Config
	jwtService  *auth.JWTService
}

func NewStoreHandler(cfg *config.Config) *Handler {
	return &Handler{
		userService: cfg.Services["user-service"],
		client:      &http.Client{Timeout: cfg.Services["user-service"].Timeout},
		cfg:         cfg,
		jwtService:  auth.NewJWTService(cfg.Auth.JWTSecret, cfg.Auth.AccessTokenExp, cfg.Auth.RefreshTokenExp),
	}
}

// CreateStore handles store creation requests, extracting user_id from JWT token
func (h *Handler) CreateStore(w http.ResponseWriter, r *http.Request) {
	// Extract user claims from context (set by auth middleware)
	claims, ok := r.Context().Value("user").(*auth.Claims)
	if !ok {
		_ = utils.ErrorJSON(w, fmt.Errorf("unauthorized"), http.StatusUnauthorized)
		return
	}

	// Read the original request body
	var storeRequest map[string]interface{}
	err := utils.ReadJSON(w, r, &storeRequest)
	if err != nil {
		_ = utils.ErrorJSON(w, fmt.Errorf("invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Remove user_id if it exists in the request (we'll use the one from token)
	delete(storeRequest, "user_id")

	// Add user ID from token to request
	storeRequest["user_id"] = claims.UserID

	// Convert back to JSON
	modifiedBody, err := json.Marshal(storeRequest)
	if err != nil {
		_ = utils.ErrorJSON(w, fmt.Errorf("error processing request: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Modified store creation request: %s", string(modifiedBody))

	// Create new request to user service
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/store", h.userService.URL),
		bytes.NewBuffer(modifiedBody),
	)
	if err != nil {
		_ = utils.ErrorJSON(w, fmt.Errorf("error creating proxy request: %v", err), http.StatusInternalServerError)
		return
	}

	// Copy headers
	req.Header = r.Header.Clone()
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := h.client.Do(req)
	if err != nil {
		_ = utils.ErrorJSON(w, fmt.Errorf("error forwarding request: %v", err), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		_ = utils.ErrorJSON(w, fmt.Errorf("error reading response: %v", err), http.StatusInternalServerError)
		return
	}

	// If store was created successfully
	if resp.StatusCode == http.StatusCreated {
		// Parse the store response to extract the store ID
		var storeResponse map[string]interface{}
		if err := json.Unmarshal(respBody, &storeResponse); err != nil {
			log.Printf("Error parsing store response: %v", err)
			// Continue anyway to return the original response
		} else {
			// Try to extract the store ID from the response
			// The structure might be different based on your API, adjust this part if needed
			var newStoreID int
			
			// Assuming the response has a structure like: { "id": 123, ... } or { "data": { "id": 123, ... } }
			if id, ok := storeResponse["id"].(float64); ok {
				newStoreID = int(id)
			} else if data, ok := storeResponse["data"].(map[string]interface{}); ok {
				if id, ok := data["id"].(float64); ok {
					newStoreID = int(id)
				}
			}
			
			if newStoreID > 0 {
				// Get the current user claims from context
				claims, ok := r.Context().Value("user").(*auth.Claims)
				if ok {
					// Check if the store ID is already in the claims
					storeExists := false
					for _, sid := range claims.StoresID {
						if sid == newStoreID {
							storeExists = true
							break
						}
					}
					
					// Add the new store ID to the claims if it doesn't exist
					if !storeExists {
						updatedStoreIDs := append(claims.StoresID, newStoreID)
						
						// Generate new tokens with updated store IDs
						accessToken, refreshToken, err := h.jwtService.GenerateTokenPair(claims.UserID, updatedStoreIDs)
						if err == nil {
							// Create an enhanced response with both the store data and new tokens
							tokenResponse := auth.TokenResponse{
								AccessToken:  accessToken,
								RefreshToken: refreshToken,
								ExpiresIn:    int64(h.jwtService.GetAccessTokenExpiry().Seconds()),
							}
							
							// Create enhanced response structure
							enhancedResponse := map[string]interface{}{
								"store_data": storeResponse,
								"tokens":     tokenResponse,
							}
							
							// Replace the response body with the enhanced version
							enhancedRespBody, err := json.Marshal(enhancedResponse)
							if err == nil {
								respBody = enhancedRespBody
							} else {
								log.Printf("Error creating enhanced response: %v", err)
							}
						} else {
							log.Printf("Error generating new tokens: %v", err)
						}
					}
				}
			}
		}
	}

	// Set content type to JSON
	w.Header().Set("Content-Type", "application/json")
	
	// Copy other response headers
	for k, v := range resp.Header {
		if k != "Content-Type" && k != "Content-Length" {
			w.Header()[k] = v
		}
	}

	// Set status code
	w.WriteHeader(resp.StatusCode)

	// Write response
	_, err = w.Write(respBody)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
