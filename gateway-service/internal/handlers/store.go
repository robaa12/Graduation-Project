package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/robaa12/gatway-service/internal/config"
	"github.com/robaa12/gatway-service/internal/middleware/auth"
	"github.com/robaa12/gatway-service/utils"
)

// StoreHandler handles store-related operations
type StoreHandler struct {
	userServiceURL    string
	productServiceURL string
	orderServiceURL   string
	client            *http.Client
	jwtService        *auth.JWTService
}

// NewStoreHandler creates a new store handler
func NewStoreHandler(config config.Config) *StoreHandler {
	return &StoreHandler{
		userServiceURL:    config.Services["user-service"].URL,
		productServiceURL: config.Services["product-service"].URL,
		orderServiceURL:   config.Services["order-service"].URL,
		jwtService:        auth.NewJWTService(config.Auth.JWTSecret, config.Auth.AccessTokenExp, config.Auth.RefreshTokenExp),
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type ServicesStoreRequest struct {
	ID   int    `json:"id"`
	Name string `json:"name" validate:"required"`
}

type StoreData struct {
	ID uint `json:"id"`
	StoreCreateRequest
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// CreateStoreResponse is the response structure for store creation
type StoreCreateResponse struct {
	Store       StoreData `json:"store"`
	Success     bool      `json:"success"`
	AccessToken string    `json:"access_token,omitempty"`
}

// ServiceResult tracks the result of operations on individual services
type ServiceResult struct {
	Success bool        `json:"success"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type StoreCreateRequest struct {
	UserID        int     `json:"user_id,omitempty"`
	StoreName     string  `json:"store_name" validate:"required"`
	Description   string  `json:"description" validate:"required"`
	BusinessPhone string  `json:"business_phone" validate:"required"`
	CategoryID    int     `json:"category_id" validate:"required"`
	PlanID        int     `json:"plan_id" validate:"required"`
	StoreCurrency string  `json:"store_currency" validate:"required"`
	Href          *string `json:"href,omitempty"`
	Slug          *string `json:"slug,omitempty"`
}

// CreateStore handles the distributed transaction for creating a store across all services
func (h *StoreHandler) CreateStore(w http.ResponseWriter, r *http.Request) {

	// Extract user claims from context (set by auth middleware)
	claims, ok := r.Context().Value("user").(*auth.Claims)
	if !ok {
		utils.ErrorJSON(w, fmt.Errorf("unauthorized"), http.StatusUnauthorized)
		return
	}

	// Read the original request body
	var storeRequest StoreCreateRequest
	err := utils.ReadJSON(w, r, &storeRequest)
	if err != nil {
		utils.ErrorJSON(w, fmt.Errorf("invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Add user ID from token to request
	storeRequest.UserID = claims.UserID

	// Step 1: Create store in user service
	userServiceResp, storeData, err := h.createStoreInUserService(&storeRequest)
	if err != nil {
		utils.ErrorJSON(w, fmt.Errorf("failed to create store in user service: %v", err), http.StatusBadGateway)
		return
	}

	if userServiceResp.StatusCode != http.StatusCreated {
		// Forward the error response from the user service
		utils.ErrorJSON(w, fmt.Errorf("user service returned status code %d", userServiceResp.StatusCode), userServiceResp.StatusCode)
		return
	}

	// Parse the user service response
	var storeResponse StoreCreateResponse
	storeResponse.Store = storeData

	// Extract the store ID from the response
	StoreRequest := ServicesStoreRequest{
		UserID:  claims.UserID,
		StoreID: storeResponse.Store.ID,
	}

	// Extract store name for other services
	if name, ok := storeResponse["store_name"].(string); ok {
		storeName = name
	} else if data, ok := storeResponse["data"].(map[string]interface{}); ok {
		if name, ok := data["store_name"].(string); ok {
			storeName = name
		}
	}

	if storeID == 0 {
		utils.ErrorJSON(w, fmt.Errorf("invalid store ID from user service"), http.StatusInternalServerError)
		return
	}

	// Step 2: Concurrently create store in other services
	var wg sync.WaitGroup
	var mu sync.Mutex
	serviceResults := make(map[string]ServiceResult)
	successfulServices := make([]string, 0)

	// Create store in product service
	wg.Add(1)
	go func() {
		defer wg.Done()
		productStoreReq := map[string]interface{}{
			"id":   storeID,
			"name": storeName,
		}
		productReqBody, _ := json.Marshal(productStoreReq)
		resp, respBody, err := h.sendRequest(http.MethodPost, h.productServiceURL+"/api/stores", productReqBody)

		mu.Lock()
		defer mu.Unlock()

		result := ServiceResult{Success: err == nil && resp != nil && resp.StatusCode == http.StatusCreated}
		if err != nil {
			result.Error = err.Error()
		} else if resp != nil && resp.StatusCode != http.StatusCreated {
			result.Error = fmt.Sprintf("Product service returned status code %d", resp.StatusCode)
		} else {
			// Store was created successfully
			successfulServices = append(successfulServices, "product_service")
			if respBody != nil {
				var data map[string]interface{}
				if json.Unmarshal(respBody, &data) == nil {
					result.Data = data
				}
			}
		}
		serviceResults["product_service"] = result
	}()

	// Create store in order service
	wg.Add(1)
	go func() {
		defer wg.Done()
		orderStoreReq := map[string]interface{}{
			"id":   storeID,
			"name": storeName,
		}
		orderReqBody, _ := json.Marshal(orderStoreReq)
		resp, respBody, err := h.sendRequest(http.MethodPost, h.orderServiceURL+"/api/stores", orderReqBody)

		mu.Lock()
		defer mu.Unlock()

		result := ServiceResult{Success: err == nil && resp != nil && resp.StatusCode == http.StatusCreated}
		if err != nil {
			result.Error = err.Error()
		} else if resp != nil && resp.StatusCode != http.StatusCreated {
			result.Error = fmt.Sprintf("Order service returned status code %d", resp.StatusCode)
		} else {
			// Store was created successfully
			successfulServices = append(successfulServices, "order_service")
			if respBody != nil {
				var data map[string]interface{}
				if json.Unmarshal(respBody, &data) == nil {
					result.Data = data
				}
			}
		}
		serviceResults["order_service"] = result
	}()

	// Wait for all goroutines to complete
	wg.Wait()

	// Step 3: Check if all services succeeded
	allSucceeded := true
	for _, result := range serviceResults {
		if !result.Success {
			allSucceeded = false
			break
		}
	}

	// Step 4: Handle compensating transactions if necessary
	if !allSucceeded {
		log.Printf("Store creation failed in some services. Initiating compensating transactions.")

		// Perform compensating transactions for successful services
		var compensationWg sync.WaitGroup

		for _, serviceName := range successfulServices {
			compensationWg.Add(1)

			go func(service string) {
				defer compensationWg.Done()

				var err error
				switch service {
				case "product_service":
					err = h.deleteStoreFromProductService(storeID)
				case "order_service":
					err = h.deleteStoreFromOrderService(storeID)
				}

				if err != nil {
					log.Printf("Compensation transaction failed for %s: %v", service, err)
				} else {
					log.Printf("Successfully performed compensating transaction for %s", service)
				}
			}(serviceName)
		}

		// Also delete from user service
		go func() {
			if err := h.deleteStoreFromUserService(storeID, r.Header); err != nil {
				log.Printf("Failed to delete store from user service: %v", err)
			} else {
				log.Printf("Successfully deleted store from user service")
			}
		}()

		// Wait for compensating transactions to complete
		compensationWg.Wait()

		// Respond with error
		utils.ErrorJSON(w, fmt.Errorf("failed to create store in all services"), http.StatusInternalServerError)
		return
	}

	// Step 5: Generate new tokens with updated store IDs
	// Add the new store ID to the claims if it doesn't exist
	storeExists := false
	for _, sid := range claims.StoresID {
		if sid == storeID {
			storeExists = true
			break
		}
	}

	var tokenResponse *auth.TokenResponse
	if !storeExists {
		updatedStoreIDs := append(claims.StoresID, storeID)

		// Generate new tokens with updated store IDs
		accessToken, refreshToken, err := h.jwtService.GenerateTokenPair(claims.UserID, updatedStoreIDs)
		if err != nil {
			log.Printf("Error generating new tokens: %v", err)
			// Continue anyway to return the store data
		} else {
			tokenResponse = &auth.TokenResponse{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
				ExpiresIn:    int64(h.jwtService.GetAccessTokenExpiry().Seconds()),
			}
		}
	}

	// Step 6: Create enhanced response with both store data and new tokens
	enhancedResponse := map[string]interface{}{
		"store_data":   storeResponse,
		"service_info": serviceResults,
		"success":      true,
	}

	if tokenResponse != nil {
		enhancedResponse["tokens"] = tokenResponse
	}

	// Return the enhanced response
	utils.WriteJSON(w, http.StatusCreated, enhancedResponse)
}

// Helper method to create store in user service
func (h *StoreHandler) createStoreInUserService(storeRequest *StoreCreateRequest) (*http.Response, StoreData, error) {
	requestBody, err := json.Marshal(storeRequest)
	if err != nil {
		return nil, StoreData{}, fmt.Errorf("marshaling request body failed: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, h.userServiceURL+"/api/stores", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, StoreData{}, fmt.Errorf("creating request failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, StoreData{}, fmt.Errorf("request failed: %w", err)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		resp.Body.Close()
		return nil, StoreData{}, fmt.Errorf("reading response body failed: %w", err)
	}
	resp.Body.Close()

	var storeData StoreData
	err = json.Unmarshal(respBody, &storeData)
	if err != nil {
		return nil, StoreData{}, fmt.Errorf("unmarshaling response body failed: %w", err)
	}

	return resp, storeData, nil
}

// Helper method to send HTTP requests
func (h *StoreHandler) sendRequest(method, url string, body []byte) (*http.Response, []byte, error) {
	var req *http.Request
	var err error

	if body != nil {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(body))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return nil, nil, fmt.Errorf("creating request failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("request failed: %w", err)
	}

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		resp.Body.Close()
		return resp, nil, fmt.Errorf("reading response body failed: %w", err)
	}
	resp.Body.Close()

	return resp, respBody, nil
}

// Compensating transaction: Delete store from user service
func (h *StoreHandler) deleteStoreFromUserService(storeID uint, authToken string) error {
	url := fmt.Sprintf("%s/api/stores/%d", h.userServiceURL, storeID)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("creating request failed: %w", err)
	}

	req.Header.Set("Authorization", authToken)

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("delete store from user service failed with status code %d", resp.StatusCode)
	}

	return nil
}

// Compensating transaction: Delete store from product service
func (h *StoreHandler) deleteStoreFromProductService(storeID uint) error {
	url := fmt.Sprintf("%s/api/stores/%d", h.productServiceURL, storeID)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("creating request failed: %w", err)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("delete store from product service failed with status code %d", resp.StatusCode)
	}

	return nil
}

// Compensating transaction: Delete store from order service
func (h *StoreHandler) deleteStoreFromOrderService(storeID uint) error {
	url := fmt.Sprintf("%s/api/stores/%d", h.orderServiceURL, storeID)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("creating request failed: %w", err)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("delete store from order service failed with status code %d", resp.StatusCode)
	}

	return nil
}
