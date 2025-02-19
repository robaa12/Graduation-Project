package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"order-service/cmd/model"
)

type ProductService struct {
	ProductServiceURL string
}

func NewProductService(url string) *ProductService {
	return &ProductService{
		ProductServiceURL: url,
	}
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
	updateRequest := struct {
		Items []model.OrderItemRequest `json:"items"`
	}{
		Items: items,
	}

	jsonData, err := json.Marshal(updateRequest)
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
