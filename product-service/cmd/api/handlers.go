package main

import (
	"log"
	"net/http"

	"github.com/robaa12/product-service/cmd/data"
	"gorm.io/gorm"
)

func (app *Config) NewProduct(w http.ResponseWriter, r *http.Request) {
	// Read the JSON request
	var productRequest data.ProductRequest
	err := app.readJSON(w, r, &productRequest)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Start Database Transaction
	err = app.db.Transaction(func(tx *gorm.DB) error {
		// Create a new Product
		var product data.Product
		product.CreateProduct(productRequest)

		// Add Product to the database
		if err := tx.Create(&product).Error; err != nil {
			log.Println("Error creating product in database")
			return err
		}

		for _, skuRequest := range productRequest.SKUs {
			// Create a new SKU
			var sku data.SKU
			sku.CreateSKU(skuRequest)

			// Add the SKU to the database
			if err := tx.Create(&sku).Error; err != nil {
				log.Println("Error creating sku in database")
				return err
			}

			for _, variantRequest := range skuRequest.Variants {
				// Create a new variant
				var variant data.Variant
				variant.CreateVariant(variantRequest)

				// Check if the variant already exists in the database or not and create it if it doesn't
				if err := tx.FirstOrCreate(&variant, data.Variant{Name: variantRequest.Name}).Error; err != nil {
					log.Println("Error creating variant in database")
					return err
				}

				// Create a new SKU Variant
				var skuVariant data.SkuVariant
				skuVariant.CreateSkuVariant(sku.ID, variant.ID, variantRequest.Value)

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
		app.errorJSON(w, err)
		return
	}
	// Return the product
	app.writeJSON(w, 201, productRequest)
}
