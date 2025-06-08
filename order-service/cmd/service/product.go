package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"order-service/cmd/model"
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

type SKUsRequest struct {
	IDs []uint `json:"sku-ids" binding:"required"`
}
type SKUsResponse struct {
	SKUs []SKUProductResponse `json:"skus"`
}
type SKUProductResponse struct {
	ID          uint   `json:"sku_id"`
	Name        string `json:"sku_name"`
	ProductID   uint   `json:"product_id"`
	ProductName string `json:"product_name"`
	ImgURL      string `json:"image_url"`
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
func (s *ProductService) GetSkuDetails(storeID uint, skuIDs []uint) (*SKUsResponse, error) {
	skusRequest := &SKUsRequest{
		IDs: skuIDs,
	}

	jsonData, err := json.Marshal(skusRequest)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(s.ProductServiceURL+fmt.Sprintf("/stores/%d/skus/info", storeID), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to verify order items")
	}

	var skusResponse SKUsResponse

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&skusResponse)
	if err != nil {
		return nil, err
	}

	return &skusResponse, nil

}
func (s *ProductService) GetOrderItemDetails(storeID uint, orderItems []model.OrderItemResponse) error {
	// Check if orderItems is empty
	if len(orderItems) == 0 {
		return nil
	}
	// intialize list of sku ids from order items and initialize map that contains id as key and index in the list as value
	skuIDs := make([]uint, 0, len(orderItems))
	skuIndexMap := make(map[uint]int, len(orderItems))
	for i, item := range orderItems {
		skuIDs = append(skuIDs, item.SkuID)
		skuIndexMap[item.SkuID] = i
	}
	skusResponse, err := s.GetSkuDetails(storeID, skuIDs)
	if err != nil {
		return fmt.Errorf("failed to get SKU details: %w", err)
	}
	for _, sku := range skusResponse.SKUs {
		if index, exists := skuIndexMap[sku.ID]; exists {
			orderItems[index].SkuName = sku.Name
			orderItems[index].ProductID = sku.ProductID
			orderItems[index].ProductName = sku.ProductName
			orderItems[index].ImageURL = sku.ImgURL
		}
	}

	return nil
}
