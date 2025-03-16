package service

import (
	"github.com/robaa12/product-service/cmd/model"
	"github.com/robaa12/product-service/cmd/repository"
)

type SKUService struct {
	repository *repository.SkuRepository
}

func NewSKUService(r *repository.SkuRepository) *SKUService {
	return &SKUService{repository: r}
}

// GetStoreProducts returns all products of a store

func (s *SKUService) UpdateSKU(skuID int, skuRequest *model.SKURequest) error {
	sku := skuRequest.ToSKU()
	sku.ID = uint(skuID)
	// Find SKU by ID
	_, err := s.repository.FindSku(skuID)
	if err != nil {
		return err
	}
	err = s.repository.UpdateSku(sku)
	if err != nil {
		return err
	}

	return nil
}

func (s *SKUService) GetSKU(skuID int) (*model.SKUResponse, error) {

	// Find SKU by ID
	sku, err := s.repository.GetSku(skuID)
	if err != nil {
		return nil, err
	}

	skuResponse := sku.ToSKUResponse()
	return skuResponse, nil
}

func (s *SKUService) DeleteSKU(skuID int) error {

	// Find SKU by ID
	sku, err := s.repository.FindSku(skuID)
	if err != nil {
		return err
	}
	err = s.repository.DeleteSKU(sku)
	if err != nil {
		return err
	}
	return nil
}

func (s *SKUService) NewSKU(storeID, productID uint, skuRequest *model.SKURequest) (*model.SKUResponse, error) {
	sku := skuRequest.CreateSKU(productID, storeID)

	// Start Database Transaction
	// Create a new SKU
	sku, err := s.repository.CreateSku(sku, skuRequest.Variants)
	if err != nil {
		return nil, err
	}
	skuResponse := sku.ToSKUResponse()
	// Return the SKU
	return skuResponse, nil
}
