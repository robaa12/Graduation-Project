package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/robaa12/product-service/cmd/data"
	"gorm.io/gorm"
)

// NewProduct creates a new product , skus and variants in the database
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
			sku.CreateSKU(skuRequest, product.ID)

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
				var skuVariant data.SKUVariant
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

// GetProduct returns a product from the database
func (app *Config) GetProduct(w http.ResponseWriter, r *http.Request) {
	var product data.Product

	// Get the product ID from the URL
	id := chi.URLParam(r, "id")

	// Get the product from the database
	err := product.GetProduct(id)

	if err != nil {
		app.errorJSON(w, err)
		return
	}
	productResponse := data.ProductResponse{
		ID:          product.ID,
		StoreID:     product.StoreID,
		Name:        product.Name,
		Description: product.Description,
	}
	// Return the product
	app.writeJSON(w, 200, productResponse)
}

// UpdateProduct updates a product in the database
func (app *Config) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	// Read the JSON request
	var product data.ProductResponse
	err := app.readJSON(w, r, &product)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	// Get the product ID from the URL
	strID := chi.URLParam(r, "id")
	if strID == "" {
		app.errorJSON(w, errors.New("product ID is required "))
		return
	}
	// Convert the product ID to an unsigned integer
	id, err := strconv.ParseUint(strID, 10, 0)
	if err != nil {
		app.errorJSON(w, errors.New("product ID must be a number"))
		return
	}
	// Update the product in the database
	err = app.db.Model(&data.Product{}).Where("id = ?", id).Updates(&product).Error
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	product.ID = uint(id)
	// Return the Updated product
	app.writeJSON(w, 200, product)
}

// DeleteProduct deletes a product from the database
func (app *Config) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	// Get the product ID from the URL
	strID := chi.URLParam(r, "id")
	if strID == "" {
		app.errorJSON(w, errors.New("product ID is required"))
		return
	}

	id, err := strconv.ParseUint(strID, 10, 0)
	if err != nil {
		app.errorJSON(w, errors.New("product ID must be a number"))
		return
	}
	tx := app.db.Begin()
	if tx.Error != nil {
		app.errorJSON(w, errors.New("Couldn't start Transaction"))
		return
	}
	err = tx.Where("sku_id IN (?)", tx.Model(&data.SKU{}).Where("product_id = ?", id).Select("id")).Delete(&data.SKUVariant{}).Error
	if err != nil {
		tx.Rollback()
		app.errorJSON(w, err)
		return
	}

	err = tx.Where("product_id = ?", id).Delete(&data.SKU{}).Error
	if err != nil {
		tx.Rollback()
		log.Println(err)
		app.errorJSON(w, err)
		return
	}
	// Delete the product from the database
	err = tx.Where("id = ?", id).Delete(&data.Product{}).Error
	if err != nil {
		tx.Rollback()
		app.errorJSON(w, errors.New("Couldn't delete Product"))
		return
	}
	err = tx.Commit().Error
	if err != nil {
		app.errorJSON(w, errors.New("Couldn't commit Transaction"))
		return
	}
	// Return the product
	app.writeJSON(w, 200, "Product and SKUs deleted successfully")
}

// GetStoreProducts returns all products of a store
func (app *Config) GetStoreProducts(w http.ResponseWriter, r *http.Request) {
	// Fetch Store ID Param From URL
	store_id := chi.URLParam(r, "store_id")
	if store_id == "" {
		app.errorJSON(w, errors.New("store_id Not Found"), 400)
		return
	}
	// create slice of Products
	var products []data.Product

	// Find All Products With Store ID
	result := app.db.Where("store_id = ?", store_id).Find(&products)
	if result.RowsAffected == 0 {
		app.errorJSON(w, errors.New("Products Not Found"), 404)
		return
	}
	if result.Error != nil {
		app.errorJSON(w, result.Error)
		return
	}
	if result.RowsAffected == 0 {
		app.errorJSON(w, errors.New("No Products Found"), 404)
		return
	}

	// Create Slice of ProductResponse
	var productsResponse []data.ProductResponse
	for _, product := range products {
		productsResponse = append(productsResponse, data.ProductResponse{
			ID:          product.ID,
			StoreID:     product.StoreID,
			Name:        product.Name,
			Description: product.Description,
		})
	}

	if result.Error != nil {
		app.errorJSON(w, result.Error)
		return
	}
	// Return store's Products
	app.writeJSON(w, 200, productsResponse)
}

// GetProduct Returns a product details from the databse
func (app *Config) GetProductDetails(w http.ResponseWriter, r *http.Request) {

	strID := chi.URLParam(r, "id")

	if strID == "" {
		app.errorJSON(w, errors.New("Product Id Not Found "))
		return
	}
	// Convert the product ID to an unsigned integer
	productId, err := strconv.ParseUint(strID, 10, 0)
	if err != nil {
		app.errorJSON(w, errors.New("product ID must be a number"))
		return
	}
	var product data.Product

	result := app.db.Where("id=?", productId).Find(&product)
	if result.Error != nil {
		app.errorJSON(w, result.Error)
		return
	}
	if result.RowsAffected == 0 {
		app.errorJSON(w, errors.New("Product Not Found"), 404)
		return
	}

	productResponse := data.ProductDetailsResponse{ID: product.ID, Name: product.Name, StoreID: product.StoreID, Description: product.Description}

	var skus []data.SKU
	result = app.db.Where("product_id = ?", productId).Find(&skus)
	if result.RowsAffected == 0 {
		app.errorJSON(w, errors.New("SKUs Not Found"), 404)
		return
	}
	if result.Error != nil {
		app.errorJSON(w, result.Error)
		return
	}
	// Iterate through each SKU to build its response structure
	for _, sku := range skus {
		// create SKU Response That Use
		var skuResponse data.SKUResponse
		skuResponse.Price = sku.Price
		skuResponse.Stock = sku.Stock

		// Retrieve SKU variants associated with the current SKU
		var sku_Varients []data.SKUVariant
		result = app.db.Where("sku_id = ?", sku.ID).Find(&sku_Varients)
		if result.RowsAffected == 0 {
			app.errorJSON(w, errors.New("SKU Variants Not Found"), 404)
			return
		}
		if result.Error != nil {
			app.errorJSON(w, result.Error)
			return
		}
		// Iterate through each SKU variant to fetch variant details
		for _, sku_varient := range sku_Varients {
			var varient data.Variant
			result = app.db.Where("id = ?", sku_varient.VariantID).Find(&varient)

			if result.Error != nil {
				app.errorJSON(w, result.Error)
				return
			}
			skuResponse.Variants = append(skuResponse.Variants, data.VariantResponse{Name: varient.Name, Value: sku_varient.Value})
		}
		productResponse.SKUs = append(productResponse.SKUs, skuResponse)
	}

	app.writeJSON(w, 200, productResponse)
}
func (app *Config) UpdateSKU(w http.ResponseWriter, r *http.Request) {
	strID := chi.URLParam(r, "id")
	if strID == "" {
		app.errorJSON(w, errors.New("SKU_ID Not Found."))
		return
	}

	// Convert the product ID to an unsigned integer
	sku_id, err := strconv.ParseUint(strID, 10, 0)
	if err != nil {
		app.errorJSON(w, errors.New("SKU ID must be a number"))
		return
	}
	var skuRequest data.SKURequest
	err = app.readJSON(w, r, &skuRequest)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	result := app.db.Where("id= ?", sku_id).Updates(&data.SKU{Stock: skuRequest.Stock, Price: skuRequest.Price})
	if result.Error != nil {
		app.errorJSON(w, result.Error)
		return
	}

	app.writeJSON(w, 200, "SKU updated successfully.")

}

func (app *Config) GetSKU(w http.ResponseWriter, r *http.Request) {
	// Fetch SKU ID Param From URL
	strID := chi.URLParam(r, "id")
	if strID == "" {
		app.errorJSON(w, errors.New("SKU_ID Not Found."))
		return
	}
	// Convert ID to uint
	skuID, err := strconv.ParseUint(strID, 10, 0)
	if err != nil {
		app.errorJSON(w, errors.New("SKU ID must be a number"))
		return
	}
	// Find SKU by ID
	var sku data.SKU
	result := app.db.Model(&data.SKU{}).Where("id = ?", skuID).Find(&sku)
	if result.RowsAffected == 0 {
		app.errorJSON(w, errors.New("SKU Not Found"), 404)
		return
	}
	// Find SKU Variants
	var SKUVariants []data.SKUVariant
	result = app.db.Where("sku_id = ?", skuID).Find(&SKUVariants)
	if result.Error != nil {
		app.errorJSON(w, result.Error)
		return
	}
	var Variants []data.VariantResponse
	for _, skuVariant := range SKUVariants {
		var variant data.Variant
		result = app.db.Where("id = ?", skuVariant.VariantID).Find(&variant)
		if result.Error != nil {
			app.errorJSON(w, result.Error)
			return
		}
		Variants = append(Variants, data.VariantResponse{Name: variant.Name, Value: skuVariant.Value})
	}
	skuResponse := data.SKUResponse{Stock: sku.Stock, Price: sku.Price, Variants: Variants}
	app.writeJSON(w, 200, skuResponse)
}

func (app *Config) DeleteSKU(w http.ResponseWriter, r *http.Request) {
	// Fetch SKU ID Param From URL
	strID := chi.URLParam(r, "id")
	if strID == "" {
		app.errorJSON(w, errors.New("SKU_ID Not Found."))
		return
	}
	// Convert ID to uint
	skuID, err := strconv.ParseUint(strID, 10, 0)
	if err != nil {
		app.errorJSON(w, errors.New("SKU ID must be a number"))
		return
	}
	// Find SKU by ID
	var sku data.SKU
	result := app.db.Where("id = ?", skuID).Find(&sku)
	if result.RowsAffected == 0 {
		app.errorJSON(w, errors.New("Sku is not found."), 404)
		return
	}
	if result.Error != nil {
		app.errorJSON(w, result.Error)
		return
	}
	// begin transaction
	tx := app.db.Begin()
	if tx.Error != nil {
		app.errorJSON(w, errors.New("Couldn't start Transaction"))
		return
	}
	// Delete SKU Variants
	result = tx.Where("sku_id = ?", skuID).Delete(&data.SKUVariant{})
	if result.Error != nil {
		tx.Rollback()
		app.errorJSON(w, result.Error)
		return
	}
	// Delete SKU
	result = tx.Where("id = ?", skuID).Delete(&data.SKU{})
	if result.Error != nil {
		tx.Rollback()
		app.errorJSON(w, result.Error)
		return
	}
	// commmit Transaction
	result = tx.Commit()
	if result.Error != nil {
		app.errorJSON(w, result.Error)
		return
	}
	app.writeJSON(w, 200, "SKU deleted successfully.")
}

func (app *Config) NewSKU(w http.ResponseWriter, r *http.Request) {
	// Read the JSON request
	var skuRequest data.SKURequest
	err := app.readJSON(w, r, &skuRequest)
	if err != nil {
		app.errorJSON(w, errors.New("Enter valid SKU data"))
		return
	}
	// Fetch Product ID Param From URL
	strID := chi.URLParam(r, "id")
	if strID == "" {
		app.errorJSON(w, errors.New("Product ID Not Found."))
		return
	}
	// Convert ID to uint
	ID, err := strconv.ParseUint(strID, 10, 0)
	if err != nil {
		app.errorJSON(w, errors.New("Product ID must be a number"))
		return
	}
	// Start Database Transaction
	// Create a new SKU
	var sku data.SKU
	sku.CreateSKU(skuRequest, uint(ID))

	// make transaction
	tx := app.db.Begin()
	// Add the SKU to the database
	if err := tx.Create(&sku).Error; err != nil {
		log.Println("Error creating sku in database")
		return
	}
	// Create SKU Variant
	for _, variant := range skuRequest.Variants {
		var variantData data.Variant
		variantData.CreateVariant(variant)
		if err := tx.FirstOrCreate(&variantData, data.Variant{Name: variant.Name}).Error; err != nil {
			log.Println("Error creating variant in database")
			return
		}
		// Create a new SKU Variant
		var skuVariant data.SKUVariant
		skuVariant.CreateSkuVariant(sku.ID, variantData.ID, variant.Value)
		if err = tx.Create(&skuVariant).Error; err != nil {
			log.Println("Error creating sku variant in database")
			return
		}
	}
	// commit transaction
	err = tx.Commit().Error
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	// Return the SKU
	app.writeJSON(w, 201, "SKU created successfully.")
}
