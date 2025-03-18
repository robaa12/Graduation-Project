package validation

import (
	"errors"

	"github.com/robaa12/product-service/cmd/model"
)

func ValidateNewProduct(storeID uint, product model.ProductRequest) error {
	if product.Name == "" {
		return errors.New("product name is required ++ ")
	}
	if product.Description == "" {
		return errors.New("product description is required")
	}
	if storeID == 0 {
		return errors.New("store ID is required")
	}

	if len(product.SKUs) == 0 {
		return errors.New("at least one SKU is required")
	}

	// Validate each SKU
	for _, sku := range product.SKUs {
		if err := ValidateSKU(sku); err != nil {
			return err
		}
	}
	return nil
}

func ValidateSKU(sku model.SKURequest) error {
	if sku.Price < 0 {
		return errors.New("price cannot be negative")
	}
	if sku.Stock < 0 {
		return errors.New("stock cannot be negative")
	}
	if sku.CostPerItem < 0 {
		return errors.New("cost per item cannot be negative")
	}
	if len(sku.Variants) == 0 {
		return errors.New("at least one variant is required for SKU")
	}
	for _, variant := range sku.Variants {
		if err := ValidateVariant(variant); err != nil {
			return err
		}
	}
	return nil
}

func ValidateVariant(variant model.VariantRequest) error {
	if variant.Name == "" {
		return errors.New("vairant name is required")
	}
	if variant.Value == "" {
		return errors.New("variant value is required")
	}
	return nil
}
