package service

import (
	"errors"
	"log"

	apperrors "github.com/robaa12/product-service/cmd/errors"
	"github.com/robaa12/product-service/cmd/model"
	"github.com/robaa12/product-service/cmd/repository"
)

type ProductService struct {
	repository    *repository.ProductRepository
	reviewService *ReviewService
}

func NewProductService(repo *repository.ProductRepository, reviewSvc *ReviewService) *ProductService {
	return &ProductService{
		repository:    repo,
		reviewService: reviewSvc,
	}
}

// NewProduct creates a new product , skus and variants in the database
func (ps *ProductService) NewProduct(storeID uint, productRequest model.ProductRequest) (*model.ProductResponse, error) {
	// Product Request Validation
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
	slug, err := ps.repository.GenerateProductSlug(productRequest.Name, storeID)
	if err != nil {
		log.Println("Error Generating Product's Slug")
		return nil, err
	}
	productRequest.Slug = slug
	product, err := ps.repository.CreateProduct(storeID, productRequest)
	if err != nil {
		return nil, err
	}
	productResponse := product.ToProductResponse()
	return productResponse, nil
}

func (ps *ProductService) GetProduct(id uint, storeID uint) (*model.ProductResponse, error) {
	product, err := ps.repository.GetProduct(id, storeID)
	err = apperrors.ErrCheck(err)
	if err != nil {
		return nil, err
	}
	productResponse := product.ToProductResponse()
	return productResponse, nil
}

func (ps *ProductService) UpdateProduct(id, storeID uint, productResponse model.ProductResponse) error {
	// Check if the product exists
	product, err := ps.repository.GetProduct(id, storeID)
	if err != nil {
		return err
	}

	if product.Name != productResponse.Name {
		slug, err := ps.repository.GenerateProductSlug(productResponse.Name, product.StoreID)
		if err != nil {
			log.Println("Error Generating Product's Slug")
			return err
		}
		productResponse.Slug = slug
	}
	err = ps.repository.UpdateProduct(productResponse, id, storeID)
	if err != nil {
		return err
	}
	return nil
}

func (ps *ProductService) DeleteProduct(productID uint, storeID uint) error {
	// Call the repository to delete the product
	err := ps.repository.DeleteProduct(productID, storeID)
	return apperrors.ErrCheck(err)
}

func (ps *ProductService) GetStoreProducts(storeID uint) ([]model.ProductResponse, error) {
	// Call the repository to get the products
	products, err := ps.repository.GetStoreProducts(storeID)
	err = apperrors.ErrCheck(err)
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

func (ps *ProductService) GetProductDetails(productID, storeID uint) (*model.ProductDetailsResponse, error) {
	// Call the repository to get the product with details
	product, err := ps.repository.GetProductDetails(productID, storeID)
	err = apperrors.ErrCheck(err)
	if err != nil {
		return nil, err
	}

	// Convert the product to a detailed response
	productDetailsResponse := product.ToProductDetailsResponse()

	// Get the review statistics
	if ps.reviewService != nil {
		reviewsStats, err := ps.reviewService.GetReviewStatistics(productID, storeID)
		if err != nil {
			return nil, err
		}
		productDetailsResponse.ReviewStatistics = reviewsStats
	}

	return productDetailsResponse, nil
}

func (ps *ProductService) GetProductBySlug(slug string, storeID uint) (*model.ProductDetailsResponse, error) {
	if slug == "" || storeID == 0 {
		return nil, errors.New("both slug and store_id are required")
	}

	// Call repository to get the product
	product, err := ps.repository.GetProductBySlug(slug, storeID)
	err = apperrors.ErrCheck(err)
	if err != nil {
		return nil, err
	}

	// Convert to detailed response
	productDetailsResponse := product.ToProductDetailsResponse()

	if ps.reviewService != nil {
		reviewsStats, err := ps.reviewService.GetReviewStatistics(product.ID, storeID)
		if err != nil {
			return nil, err
		}
		productDetailsResponse.ReviewStatistics = reviewsStats
	}
	return productDetailsResponse, nil
}
