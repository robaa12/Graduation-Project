package service

import (
	"errors"
	"log"
	"time"

	apperrors "github.com/robaa12/product-service/cmd/errors"
	"github.com/robaa12/product-service/cmd/model"
	"github.com/robaa12/product-service/cmd/repository"
	"gorm.io/gorm"
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

	if err := validateProductImages(productRequest); err != nil {
		return nil, err
	}

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

func (ps *ProductService) UpdateProduct(id, storeID uint, productResponse model.ProductResponse) (*model.ProductResponse, error) {
	// Check if the product exists
	product, err := ps.repository.GetProduct(id, storeID)
	if err != nil {
		return nil, err
	}

	if product.Name != productResponse.Name {
		slug, err := ps.repository.GenerateProductSlug(productResponse.Name, product.StoreID)
		if err != nil {
			log.Println("Error Generating Product's Slug")
			return nil, err
		}
		productResponse.Slug = slug
	}

	// Update the product
	updatedProduct, err := ps.repository.UpdateProduct(productResponse, id, storeID)
	if err != nil {
		return nil, err
	}

	// Convert to response and return
	return updatedProduct.ToProductResponse(), nil
}

func (ps *ProductService) DeleteProduct(productID uint, storeID uint) error {
	// Call the repository to delete the product
	err := ps.repository.DeleteProduct(productID, storeID)
	return apperrors.ErrCheck(err)
}

func (ps *ProductService) GetStoreProducts(storeID uint, limit, offset int) (*model.PaginatedProductsResponse, error) {
	// Check if we're fetching all products or using pagination
	isPaginated := limit > 0

	if isPaginated {
		log.Printf("GetStoreProducts: Paginated request - storeID=%d, limit=%d, offset=%d", storeID, limit, offset)
	} else {
		log.Printf("GetStoreProducts: Fetching all products for storeID=%d", storeID)
	}

	// Call the repository to get the products
	products, total, err := ps.repository.GetStoreProducts(storeID, limit, offset)
	err = apperrors.ErrCheck(err)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("Error getting products: %v", err)
		return nil, err
	}

	log.Printf("Retrieved %d products out of total %d", len(products), total)

	// Create response with pagination info
	paginatedResponse := model.GetPaginatedProductsResponse(products, total, limit, offset, isPaginated)
	return paginatedResponse, nil
}

func (ps *ProductService) GetStoreProductsDashboard(storeID uint, startDate, endDate time.Time) (*model.ProductsDashboardResponse, error) {
	if startDate.IsZero() || endDate.IsZero() {
		startDate = time.Now().AddDate(0, -30, 0) // Default to 30 days ago
		if endDate.IsZero() {
			endDate = time.Now()
		}
	}

	// Call the repository to get the products
	productsDashboardResponse, err := ps.repository.GetStoreProductsDashboard(storeID, startDate, endDate)
	err = apperrors.ErrCheck(err)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("Error getting products: %v", err)
		return nil, err
	}

	log.Printf("Retrieved %v products dashboard", productsDashboardResponse)

	return productsDashboardResponse, nil
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

	if product.CategoryID != nil {
		relatedProducts, err := ps.repository.GetRelatedProducts(product.ID, *product.CategoryID, storeID, 4)
		if err == nil {
			relatedProductResponses := make([]model.ProductResponse, 0, len(relatedProducts))
			for _, rp := range relatedProducts {
				relatedProductResponses = append(relatedProductResponses, *rp.ToProductResponse())
			}
			productDetailsResponse = productDetailsResponse.WithRelatedProducts(relatedProductResponses)
		}
	}

	if ps.reviewService != nil {
		reviewsStats, err := ps.reviewService.GetReviewStatistics(product.ID, storeID)
		if err != nil {
			return nil, err
		}
		productDetailsResponse.ReviewStatistics = reviewsStats
	}
	return productDetailsResponse, nil
}

func validateProductImages(product model.ProductRequest) error {
	if product.MainImageURL == "" {
		return apperrors.NewBadRequestError("main image URL is required")
	}

	const maxAdditionalImages = 10
	if len(product.ImagesURL) > maxAdditionalImages {
		return apperrors.NewBadRequestError("maximum of 10 additional images are allowed")
	}
	seenURLs := make(map[string]bool)
	seenURLs[product.MainImageURL] = true
	for _, url := range product.ImagesURL {
		if url == "" {
			continue
		}
		if seenURLs[url] {
			continue
		}
		seenURLs[url] = true
	}
	return nil
}

// GetProductsByStoreSlug retrieves all products for a store identified by its slug
func (ps *ProductService) GetProductsByStoreSlug(storeSlug string, limit, offset int) (*model.PaginatedProductsResponse, error) {
	// Check if we're fetching all products or using pagination
	isPaginated := limit > 0

	if isPaginated {
		log.Printf("GetProductsByStoreSlug: Paginated request - storeSlug=%s, limit=%d, offset=%d", storeSlug, limit, offset)
	} else {
		log.Printf("GetProductsByStoreSlug: Fetching all products for storeSlug=%s", storeSlug)
	}

	// Call the repository to get the products
	products, storeID, total, err := ps.repository.GetProductsByStoreSlug(storeSlug, limit, offset)
	if err != nil {
		return nil, err
	}

	log.Printf("Retrieved %d products out of total %d for store %d (slug: %s)", len(products), total, storeID, storeSlug)

	// Create response with pagination info
	paginatedResponse := model.GetPaginatedProductsResponse(products, total, limit, offset, isPaginated)
	return paginatedResponse, nil
}
