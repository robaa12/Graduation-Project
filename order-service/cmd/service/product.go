package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"order-service/cmd/model"
	"time"
)

type ProductService struct {
	ProductServiceURL string
	client            *http.Client
}
type VerificationResponse struct {
	Valid   bool           `json:"valid"`
	Message string         `json:"messages"`
	Items   []VerifiedItem `json:"items"`
}

type VerifiedItem struct {
	SkuID   uint    `json:"sku_id"`
	Valid   bool    `json:"valid"`
	InStock bool    `json:"in_stock"`
	Price   float64 `json:"actual_price"`
}

type SkuDetail struct {
	SkuID       uint   `json:"sku_id"`
	SkuName     string `json:"sku_name"`
	ProductID   uint   `json:"product_id"`
	ProductName string `json:"product_name"`
	ImageURL    string `json:"image_url"`
}

type SkuDetailsResponse struct {
	Status  bool       `json:"status"`
	Message string     `json:"message"`
	Data    []SkuDetail `json:"data"`
}

func (s *ProductService) VerifyOrderItems(storeID uint, items []model.OrderItemRequest) error {
	verificationRequest := struct {
		StoreID uint                     `json:"store_id"`
		Items   []model.OrderItemRequest `json:"items"`
	}{
		StoreID: storeID,
		Items:   items,
	}

	jsonData, err := json.Marshal(verificationRequest)
	if err != nil {
		return err
	}

	resp, err := http.Post(s.ProductServiceURL+"/verify-order", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to verify order items")
	}

	var VerificationResponse VerificationResponse

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&VerificationResponse)
	if err != nil {
		return err
	}
	if !VerificationResponse.Valid {
		return errors.New(VerificationResponse.Message)

	}

	return nil
}

func (s *ProductService) UpdateInventory(items []model.OrderItemRequest) error {
	jsonData, err := json.Marshal(items)
	if err != nil {
		return err
	}

	resp, err := http.Post(s.ProductServiceURL+"/update-inventory", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to update inventory")
	}

	return nil
}

// GetSkuDetails fetches detailed information about SKUs from product service
func (s *ProductService) GetSkuDetails(storeID uint, skuIDs []uint) (map[uint]SkuDetail, error) {
	// Create a map to store the sku details by sku ID for easy lookup
	skuDetailsMap := make(map[uint]SkuDetail)
	
	// Prepare request body
	requestBody := struct {
		StoreID uint  `json:"store_id"`
		SkuIDs  []uint `json:"sku_ids"`
	}{
		StoreID: storeID,
		SkuIDs:  skuIDs,
	}

	// Convert request to JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal sku details request: %w", err)
	}

	// Initialize HTTP client if not already done
	if s.client == nil {
		s.client = &http.Client{Timeout: 10 * time.Second}
	}

	// Make request to product service
	resp, err := s.client.Post(
		fmt.Sprintf("%s/skus/details", s.ProductServiceURL),
		"application/json", 
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to product service: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("product service returned non-OK status: %d", resp.StatusCode)
	}

	// Parse the response
	var response SkuDetailsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to parse sku details response: %w", err)
	}

	// Check if response was successful
	if !response.Status {
		return nil, fmt.Errorf("product service error: %s", response.Message)
	}

	// Map the data for fast lookup by SKU ID
	for _, skuDetail := range response.Data {
		skuDetailsMap[skuDetail.SkuID] = skuDetail
	}

	return skuDetailsMap, nil
}
