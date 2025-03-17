package repository

import (
	"fmt"
	"log"
	"strings"

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
	err := pr.db.DB.Where("id = ? AND store_id = ?", productId, storeId).First(&product).Error
	if err != nil {
		return nil, gorm.ErrRecordNotFound
	}
	return &product, nil
}

func (pr *ProductRepository) UpdateProduct(p model.ProductResponse, id uint, storeId uint) error {
	return pr.db.DB.Model(&model.Product{}).Where("id = ? AND store_id = ?", id, storeId).Updates(p).Error
}

func (pr *ProductRepository) CreateProduct(productRequest model.ProductRequest) (*model.Product, error) {
	// Generate product slug

	product := productRequest.CreateProduct()

	err := pr.db.DB.Transaction(func(tx *gorm.DB) error {

		// Add Product to the database
		if err := tx.Create(&product).Error; err != nil {
			log.Println("Error creating product in database")
			return err
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

	result := pr.db.DB.Where("id = ? AND store_id = ?", productID, storeID).Delete(&product).Preload("SKUs.SKUVariants").
		Preload("SKUs.Variants").
		Find(&product)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return pr.db.DB.Unscoped().Delete(&product).Error
}

func (pr *ProductRepository) GetStoreProducts(storeID uint) ([]model.Product, error) {
	var products []model.Product

	result := pr.db.DB.Model(&model.Product{}).Where("store_id = ?", storeID).Find(&products)

	if result.Error != nil {
		return nil, result.Error
	}

	if len(products) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return products, nil
}

func (pr *ProductRepository) GetProductDetails(productID uint, storeID uint) (*model.Product, error) {
	var product model.Product

	result := pr.db.DB.Where("id=? AND store_id = ?", productID, storeID).Preload("SKUs.SKUVariants").Preload("SKUs.Variants").Find(&product)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &product, nil
}

func (pr *ProductRepository) GetProductBySlug(slug string, storeID uint) (*model.Product, error) {
	var product model.Product

	result := pr.db.DB.Where("slug = ? AND store_id = ?", slug, storeID).
		Preload("SKUs.SKUVariants").
		Preload("SKUs.Variants").
		First(&product)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &product, nil
}
