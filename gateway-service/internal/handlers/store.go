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
}

func NewStoreHandler(cfg *config.Config) *Handler {
	return &Handler{
		userService: cfg.Services["user-service"],
		client:      &http.Client{Timeout: cfg.Services["user-service"].Timeout},
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

	// Copy response headers
	for k, v := range resp.Header {
		w.Header()[k] = v
	}

	// Set status code
	w.WriteHeader(resp.StatusCode)

	// Write response
	_, err = w.Write(respBody)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
