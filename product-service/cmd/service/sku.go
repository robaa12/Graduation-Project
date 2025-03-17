package service

import (
	apperrors "github.com/robaa12/product-service/cmd/errors"
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

func (s *SKUService) UpdateSKU(skuID, productID, storeID uint, skuRequest *model.SKURequest) error {
	sku := skuRequest.CreateSKU(productID)
	sku.ID = skuID
	// Find SKU by ID
	_, err := s.repository.FindSku(skuID, sku.ProductID, storeID)
	err = apperrors.ErrCheck(err)
	if err != nil {
		return err
	}
	err = s.repository.UpdateSku(sku, storeID)
	err = apperrors.ErrCheck(err)
	if err != nil {
		return err
	}

	return nil
}

func (s *SKUService) GetSKU(skuID, productID, storeID uint) (*model.SKUResponse, error) {

	// Find SKU by ID
	sku, err := s.repository.GetSku(skuID, productID, storeID)
	err = apperrors.ErrCheck(err)
	if err != nil {
		return nil, err
	}

	skuResponse := sku.ToSKUResponse()
	return skuResponse, nil
}

func (s *SKUService) DeleteSKU(skuID, productID, storeID uint) error {

	// Find SKU by ID
	sku, err := s.repository.FindSku(skuID, productID, storeID)
	err = apperrors.ErrCheck(err)
	if err != nil {
		return err
	}
	err = s.repository.DeleteSKU(sku, storeID)
	err = apperrors.ErrCheck(err)
	if err != nil {
		return err
	}
	return nil
}

func (s *SKUService) NewSKU(storeID, productID uint, skuRequest *model.SKURequest) (*model.SKUResponse, error) {
	// check if the product exists in the store

	sku := skuRequest.CreateSKU(productID)

	// Start Database Transaction
	// Create a new SKU
	sku, err := s.repository.CreateSku(storeID, sku, skuRequest.Variants)
	err = apperrors.ErrCheck(err)
	if err != nil {
		return nil, err
	}
	skuResponse := sku.ToSKUResponse()
	// Return the SKU
	return skuResponse, nil
}
