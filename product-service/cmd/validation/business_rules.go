package validation

import (
	"errors"

	"github.com/robaa12/product-service/cmd/model"
)

func ValidateBusinessRules(product model.ProductRequest) error {
	if len(product.Name) > 255 {
		return errors.New("product name cannot exceed 255 characters")
	}

	if len(product.Description) > 1000 {
		return errors.New("product description cannot exceed 1000 characters")
	}

	// Price validation
	for _, sku := range product.SKUs {
		if sku.CompareAtPrice > 0 && sku.Price > sku.CompareAtPrice {
			return errors.New("regular price cannot be greater than compare-at price")
		}

		if sku.CostPerItem > sku.Price {
			return errors.New("cost per item cannot be greater than selling price")
		}
	}

	return nil
}
