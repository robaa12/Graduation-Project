package service

import (
	"log"
	"strconv"

	"github.com/robaa12/product-service/cmd/model"

	"gorm.io/gorm"
)

type ProductHandler struct {
	DB *gorm.DB
}

// NewProduct creates a new product , skus and variants in the database
func (h *ProductHandler) NewProduct(productRequest model.ProductRequest) (*model.ProductResponse, error) {

	// Create a new Product
	product := productRequest.CreateProduct()

	// Start Database Transaction
	err := h.DB.Transaction(func(tx *gorm.DB) error {

		// Add Product to the database
		if err := tx.Create(&product).Error; err != nil {
			log.Println("Error creating product in database")
			return err
		}

		for _, skuRequest := range productRequest.SKUs {
			// Create a new SKU
			sku := skuRequest.CreateSKU(product.ID, product.StoreID)

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
		log.Println(err)
		return nil, err
	}
	productResponse := product.ToProductResponse()
	return productResponse, nil
}

// GetProduct returns a product from the database
func (h *ProductHandler) GetProduct(productID uint) (*model.ProductResponse, error) {
	var product model.Product
	// Get the product from the database
	err := product.GetProduct(strconv.Itoa(int(productID)))

	if err != nil {
		return nil, err
	}
	// TO DO	: MAPPING FOR REFACTORING

	productResponse := product.ToProductResponse()
	// Return the product
	return productResponse, nil
}
