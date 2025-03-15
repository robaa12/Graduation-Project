package service

import (
	"errors"
	"log"

	"github.com/robaa12/product-service/cmd/model"
	"github.com/robaa12/product-service/cmd/repository"
)

type ProductService struct {
	repository *repository.ProductRepository
}

// NewProduct creates a new product , skus and variants in the database
func (ps *ProductService) NewProduct(productRequest model.ProductRequest) (*model.ProductResponse, error) {
	// Prouct Request Validation
	if len(productRequest.Name) > 255 {
		return nil, errors.New("product name cannot exceed 255 characters")
	}

	if len(productRequest.Description) > 1000 {
		return nil, errors.New("product description cannot exceed 1000 characters")
	}

	// Price validation
	for _, sku := range productRequest.SKUs {
		if sku.CompareAtPrice > 0 && sku.Price > sku.CompareAtPrice {
			return nil, errors.New("regular price cannot be greater than compare-at price")
		}

		if sku.CostPerItem > sku.Price {
			return nil, errors.New("cost per item cannot be greater than selling price")
		}
	}
	// TO DO : VALIDATION
	slug, err := ps.repository.GenerateProductSlug(productRequest.Name, productRequest.StoreID)
	if err != nil {
		log.Println("Error Genrating Product's Slug")
		return nil, err
	}
	productRequest.Slug = slug
	product, err := ps.repository.CreateProduct(productRequest)
	if err != nil {
		return nil, err
	}
	productResponse := product.ToProductResponse()
	return productResponse, nil
}
