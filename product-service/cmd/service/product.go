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

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repository: repo}
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

func (ps *ProductService) GetProduct(id uint) (*model.ProductResponse, error) {
	product, err := ps.repository.GetProduct(id)
	if err != nil {
		return nil, err
	}
	productResponse := product.ToProductResponse()
	return productResponse, nil
}

func (ps *ProductService) UpdateProduct(id uint, productResponse model.ProductResponse) error {
	// Check if the product exists
	product, err := ps.repository.GetProduct(id)
	if err != nil {
		return err
	}

	if product.Name != productResponse.Name {
		slug, err := ps.repository.GenerateProductSlug(productResponse.Name, product.StoreID)
		if err != nil {
			log.Println("Error Genrating Product's Slug")
			return err
		}
		productResponse.Slug = slug
	}
	err = ps.repository.UpdateProduct(productResponse, id)
	if err != nil {
		return err
	}
	return nil
}

func (ps *ProductService) DeleteProduct(productID uint, storeID uint) error {
	// Call the repository to delete the product
	return ps.repository.DeleteProduct(productID, storeID)
}

func (ps *ProductService) GetStoreProducts(storeID uint) ([]model.ProductResponse, error) {
	// Call the repository to get the products
	products, err := ps.repository.GetStoreProducts(storeID)
	if err != nil {
		return nil, err
	}

	// Convert products to response objects
	var productsResponse []model.ProductResponse
	for _, product := range products {
		productResponse := product.ToProductResponse()
		productsResponse = append(productsResponse, *productResponse)
	}

	return productsResponse, nil
}

func (ps *ProductService) GetProductDetails(productID uint) (*model.ProductDetailsResponse, error) {
	// Call the repository to get the product with details
	product, err := ps.repository.GetProductDetails(productID)
	if err != nil {
		return nil, err
	}

	// Convert the product to a detailed response
	productDetailsResponse := product.ToProductDetailsResponse()
	return productDetailsResponse, nil
}

func (ps *ProductService) GetProductBySlug(slug string, storeID uint) (*model.ProductDetailsResponse, error) {
	if slug == "" || storeID == 0 {
		return nil, errors.New("both slug and store_id are required")
	}

	// Call repository to get the product
	product, err := ps.repository.GetProductBySlug(slug, storeID)
	if err != nil {
		return nil, err
	}

	// Convert to detailed response
	productDetailsResponse := product.ToProductDetailsResponse()
	return productDetailsResponse, nil
}
