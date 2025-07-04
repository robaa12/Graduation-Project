package repository

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/robaa12/product-service/cmd/database"
	"github.com/robaa12/product-service/cmd/model"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db database.Database
}

func NewProductRepository(db database.Database) *ProductRepository {
	return &ProductRepository{db: db}
}

func (pr *ProductRepository) GetProduct(productId uint, storeId uint) (*model.Product, error) {
	var product model.Product
	// Find the product with the given id and store_id
	err := pr.db.DB.Preload("Category").Preload("SKUs").Where("id = ? AND store_id = ?", productId, storeId).First(&product).Error
	if err != nil {
		return nil, gorm.ErrRecordNotFound
	}

	// Fetch collection IDs
	product.CollectionIDs = pr.fetchCollectionIDs(product.ID)

	return &product, nil
}

func (pr *ProductRepository) UpdateProduct(p model.ProductResponse, id uint, storeId uint) (*model.Product, error) {
	// Find the product with the given id and store_id
	var product model.Product
	err := pr.db.DB.Where("id = ? AND store_id = ?", id, storeId).First(&product).Error
	if err != nil {
		return nil, err
	}

	// Create a map for updates with the correct field types
	updates := map[string]interface{}{}

	// Only add fields that are non-empty
	if p.Name != "" {
		updates["name"] = p.Name
	}

	if p.Description != "" {
		updates["description"] = p.Description
	}

	// Include boolean fields
	updates["published"] = p.Published

	// Only update price if it's set
	if p.StartPrice > 0 {
		updates["start_price"] = p.StartPrice
	}

	if p.Slug != "" {
		updates["slug"] = p.Slug
	}

	if p.MainImageURL != "" {
		updates["main_image_url"] = p.MainImageURL
	}

	// Convert []string to pq.StringArray for ImagesURL
	if p.ImagesURL != nil {
		updates["images_url"] = pq.StringArray(p.ImagesURL)
	}

	// Update CategoryID if Category is provided
	if p.Category != nil {
		updates["category_id"] = p.Category.ID
	}

	// Apply updates to the product
	if err := pr.db.DB.Model(&product).Updates(updates).Error; err != nil {
		return nil, err
	}

	// Fetch the updated product with related data
	var updatedProduct model.Product
	if err := pr.db.DB.Preload("Category").Where("id = ?", id).First(&updatedProduct).Error; err != nil {
		return nil, err
	}

	updatedProduct.CollectionIDs = pr.fetchCollectionIDs(updatedProduct.ID)

	return &updatedProduct, nil
}

func (pr *ProductRepository) CreateProduct(storeID uint, productRequest model.ProductRequest) (*model.Product, error) {
	// Generate product slug

	product := productRequest.CreateProduct(storeID)

	err := pr.db.DB.Transaction(func(tx *gorm.DB) error {
		// TODO: remove when adding Distributed Transaction
		//add new store if not exist in database using firstorcreate
		store := model.Store{ID: storeID}
		if err := tx.FirstOrCreate(&store, store).Error; err != nil {
			log.Println("Error creating store in database")
			return err
		}

		// Check if the category already exists in the database or not
		if product.CategoryID != nil {
			var category model.Category
			result := tx.Where("store_id = ? AND id = ?", storeID, product.CategoryID).First(&category)

			if err := result.Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return fmt.Errorf("category with id %d not found", *product.CategoryID)
				}
				return result.Error
			}
			product.Category = category
		}

		// Add Product to the database
		if err := tx.Create(&product).Error; err != nil {
			log.Println("Error creating product in database")
			return err
		}
		var verifyProduct model.Product
		if err := tx.First(&verifyProduct, product.ID).Error; err != nil {
			return err
		}

		if len(verifyProduct.ImagesURL) != len(product.ImagesURL) {
			return fmt.Errorf("failed to store all images")
		}

		for _, skuRequest := range productRequest.SKUs {
			// Create a new SKU
			sku := skuRequest.CreateSKU(product.ID)

			// Add the SKU to the database
			if err := tx.Create(&sku).Error; err != nil {
				log.Println("Error creating sku in database")
				return err
			}

			for _, variantRequest := range skuRequest.Variants {
				// Create a new variant
				variant := variantRequest.CreateVariant()

				// Check if the variant already exists in the database or not and create it if it doesn't
				if err := tx.FirstOrCreate(&variant, model.Variant{Name: variantRequest.Name}).Error; err != nil {
					log.Println("Error creating variant in database")
					return err
				}

				// Create a new SKU Variant
				skuVariant := model.CreateSkuVariant(sku.ID, variant.ID, variantRequest.Value)

				// Add the SKU Variant to the database
				if err := tx.Create(&skuVariant).Error; err != nil {
					log.Println("Error creating sku variant in database")
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (pr *ProductRepository) GenerateProductSlug(name string, storeID uint) (string, error) {
	// Generate the base slug from the product name
	baseSlug := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	slug := baseSlug
	var count int64

	// Loop to find a unique slug
	for i := 1; ; i++ {
		// Check if a product with the same slug and store_id already exists
		pr.db.DB.Model(&model.Product{}).
			Where("slug = ? AND store_id = ?", slug, storeID).
			Count(&count)

		// If no duplicate is found, return the unique slug
		if count == 0 {
			return slug, nil
		}

		// If a duplicate exists, append a counter to the slug
		slug = fmt.Sprintf("%s-%d", baseSlug, i)
	}
}

func (pr *ProductRepository) DeleteProduct(productID uint, storeID uint) error {
	var product model.Product

	// First check if the product exists
	if err := pr.db.DB.Where("id = ? AND store_id = ?", productID, storeID).First(&product).Error; err != nil {
		return err // This will return gorm.ErrRecordNotFound if the product doesn't exist
	}

	// If you need to load related data for cascade deletion, do it here
	// But don't combine it with the delete operation
	if err := pr.db.DB.Preload("SKUs.SKUVariants").Preload("SKUs.Variants").
		Where("id = ?", productID).First(&product).Error; err != nil {
		return err
	}

	// Perform the actual delete (choose one approach):
	// Option 1: Soft delete (keeps record but marks as deleted)
	// return pr.db.DB.Delete(&product).Error

	// Option 2: Hard delete (completely removes the record)
	return pr.db.DB.Unscoped().Delete(&product).Error
}

func (pr *ProductRepository) GetStoreProducts(storeID uint, limit, offset int) ([]model.Product, int64, error) {
	products := []model.Product{}
	var total int64

	// Count total products for pagination info
	if err := pr.db.DB.Model(&model.Product{}).Where("store_id = ?", storeID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Build the base query
	query := pr.db.DB.Model(&model.Product{}).
		Preload("Category").
		Preload("SKUs").
		Where("store_id = ?", storeID).
		Order("id ASC")

	// Apply limit only if pagination is requested
	if limit > 0 {
		query = query.Limit(limit)
	}

	// Apply offset only if pagination is requested
	if limit > 0 || offset > 0 {
		query = query.Offset(offset)
	}

	result := query.Find(&products)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	// For each product, fetch the collection IDs
	for i := range products {
		products[i].CollectionIDs = pr.fetchCollectionIDs(products[i].ID)
	}

	if len(products) == 0 && offset == 0 {
		return nil, 0, gorm.ErrRecordNotFound
	}

	return products, total, nil
}

func (pr *ProductRepository) GetProductDetails(productID uint, storeID uint) (*model.Product, error) {
	var product model.Product

	result := pr.db.DB.Where("id=? AND store_id = ?", productID, storeID).
		Preload("SKUs.SKUVariants").
		Preload("SKUs.Variants").
		Preload("Category").
		Find(&product)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	// Fetch collection IDs
	product.CollectionIDs = pr.fetchCollectionIDs(product.ID)

	return &product, nil
}

func (pr *ProductRepository) GetProductBySlug(slug string, storeID uint) (*model.Product, error) {
	var product model.Product

	result := pr.db.DB.Where("slug = ? AND store_id = ?", slug, storeID).
		Preload("SKUs.SKUVariants").
		Preload("SKUs.Variants").
		Preload("Category").
		First(&product)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	// Fetch collection IDs
	product.CollectionIDs = pr.fetchCollectionIDs(product.ID)

	return &product, nil
}

func (pr *ProductRepository) GetRelatedProducts(productID uint, categoryID uint, storeID uint, limit int) ([]model.Product, error) {
	var products []model.Product

	result := pr.db.DB.Where("id != ? AND category_id = ? AND store_id = ? AND published = ?",
		productID, categoryID, storeID, true).
		Order("RANDOM()").
		Limit(limit).
		Preload("Category").
		Preload("SKUs"). // Add this to preload SKUs
		Find(&products)

	// Fetch collection IDs for each product
	for i := range products {
		products[i].CollectionIDs = pr.fetchCollectionIDs(products[i].ID)
	}

	return products, result.Error
}

// Helper method to fetch collection IDs for a product
func (pr *ProductRepository) fetchCollectionIDs(productID uint) []uint {
	var collectionIDs []uint
	if err := pr.db.DB.Table("collection_products").
		Select("collection_id").
		Where("product_id = ?", productID).
		Pluck("collection_id", &collectionIDs).Error; err != nil {
		// Just log the error and return empty array
		log.Printf("Error fetching collection IDs for product %d: %v", productID, err)
		return []uint{}
	}
	return collectionIDs
}

// GetProductsByStoreSlug retrieves all products for a store identified by its slug
func (pr *ProductRepository) GetProductsByStoreSlug(storeSlug string, limit, offset int) ([]model.Product, uint, int64, error) {
	products := []model.Product{}
	var total int64
	var storeID uint

	// First, find the store by slug
	var store model.Store
	if err := pr.db.DB.Where("slug = ?", storeSlug).First(&store).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, 0, errors.New("store not found")
		}
		return nil, 0, 0, err
	}

	storeID = store.ID

	// Count total products for pagination info
	if err := pr.db.DB.Model(&model.Product{}).Where("store_id = ?", storeID).Count(&total).Error; err != nil {
		return nil, storeID, 0, err
	}

	// Build the base query
	query := pr.db.DB.Model(&model.Product{}).
		Preload("Category").
		Preload("SKUs").
		Where("store_id = ?", storeID).
		Order("id ASC")

	// Apply limit only if pagination is requested
	if limit > 0 {
		query = query.Limit(limit)
	}

	// Apply offset only if pagination is requested
	if limit > 0 || offset > 0 {
		query = query.Offset(offset)
	}

	result := query.Find(&products)

	if result.Error != nil {
		return nil, storeID, 0, result.Error
	}

	// For each product, fetch the collection IDs
	for i := range products {
		products[i].CollectionIDs = pr.fetchCollectionIDs(products[i].ID)
	}

	if len(products) == 0 && offset == 0 {
		return nil, storeID, total, nil // Return empty products with total count
	}

	return products, storeID, total, nil
}
func (pr *ProductRepository) GetStoreProductsDashboard(storeID uint, startDate, endDate time.Time) (*model.ProductsDashboardResponse, error) {
	var totalProducts int64
	var productsChange float64

	// Count total products in the current period
	if err := pr.db.DB.Model(&model.Product{}).
		Where("store_id = ? AND created_at BETWEEN ? AND ?", storeID, startDate, endDate).
		Count(&totalProducts).Error; err != nil {
		return nil, err
	}

	// Calculate previous period
	periodDuration := endDate.Sub(startDate)
	prevPeriodEnd := startDate
	prevPeriodStart := startDate.Add(-periodDuration)

	// Count products in the previous period
	var previousCount int64
	if err := pr.db.DB.Model(&model.Product{}).
		Where("store_id = ? AND created_at BETWEEN ? AND ?", storeID, prevPeriodStart, prevPeriodEnd).
		Count(&previousCount).Error; err != nil {
		return nil, err
	}

	// Calculate percentage change
	if previousCount > 0 {
		productsChange = float64(totalProducts-previousCount) / float64(previousCount) * 100.0
	} else {
		productsChange = float64(totalProducts) * 100.0 // If no previous products, consider it a full increase
	}

	return &model.ProductsDashboardResponse{
		TotalProducts:  totalProducts,
		ProductsChange: productsChange,
	}, nil
}
