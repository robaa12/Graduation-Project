package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/robaa12/product-service/cmd/data"
	"github.com/robaa12/product-service/cmd/utils"
	"gorm.io/gorm"
)

type ProductHandler struct {
	DB *gorm.DB
}

// NewProduct creates a new product , skus and variants in the database
func (h *ProductHandler) NewProduct(w http.ResponseWriter, r *http.Request) {
	// Read the JSON request
	var productRequest data.ProductRequest
	err := utils.ReadJSON(w, r, &productRequest)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	// Start Database Transaction
	err = h.DB.Transaction(func(tx *gorm.DB) error {
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
			var sku data.Sku
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
		utils.ErrorJSON(w, err)
		return
	}
	// Return the product
	utils.WriteJSON(w, 201, productRequest)
}

// GetProduct returns a product from the database
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	var product data.Product

	// Get the product ID from the URL
	id, err := utils.GetID(r, "id")
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	// Get the product from the database
	err = product.GetProduct(strconv.Itoa(int(id)))

	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}
	// TO DO	: MAPPING FOR REFACTORING

	productResponse := product.ToProductResponse()
	// Return the product
	utils.WriteJSON(w, 200, productResponse)
}

// UpdateProduct updates a product in the database
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	// Read the JSON request
	var product data.ProductResponse
	err := utils.ReadJSON(w, r, &product)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}
	// Get the product ID from the URL
	id, err := utils.GetID(r, "id")
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}
	// Update the product in the database
	err = h.DB.Model(&data.Product{}).Where("id = ?", id).Updates(&product).Error
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}
	product.ID = id
	// Return the Updated product
	utils.WriteJSON(w, 200, product)
}

// DeleteProduct deletes a product from the database
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	// Get the product ID from the URL
	id, err := utils.GetID(r, "id")
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}
	var product data.Product
	result := h.DB.Where("id=?", id).Preload("SKUs.SKUVariants").Preload("SKUs.Variants").Find(&product)
	if result.Error != nil {
		utils.ErrorJSON(w, result.Error)
		return
	}
	if result.RowsAffected == 0 {
		utils.ErrorJSON(w, errors.New("Product Not Found"), 404)
		return
	}
	// Delete the product from the database
	err = h.DB.Unscoped().Delete(&product).Error
	if err != nil {
		utils.ErrorJSON(w, errors.New("Couldn't delete Product"))
		return
	}
	// Return the product
	utils.WriteJSON(w, 200, "Product deleted successfully")
}

func (h *ProductHandler) GetStoreProducts(w http.ResponseWriter, r *http.Request) {
	// Fetch Store ID Param From URL
	storeID, err := utils.GetID(r, "store_id")
	if err != nil {
		utils.ErrorJSON(w, err, 400)
		return
	}
	// create slice of Products
	var products []data.Product

	// Find All Products With Store ID
	result := h.DB.Where("store_id = ?", storeID).Find(&products)
	if result.RowsAffected == 0 {
		utils.ErrorJSON(w, errors.New("Products Not Found"), 404)
		return
	}
	if result.Error != nil {
		utils.ErrorJSON(w, result.Error)
		return
	}

	// Create Slice of ProductResponse
	var productsResponse []data.ProductResponse
	for _, product := range products {
		productResponse := product.ToProductResponse()
		productsResponse = append(productsResponse, productResponse)
	}
	// Return store's Products
	utils.WriteJSON(w, 200, productsResponse)
}

// GetProductDetails returns a product details from the database
func (h *ProductHandler) GetProductDetails(w http.ResponseWriter, r *http.Request) {
	// Get the product ID from the URL
	productId, err := utils.GetID(r, "id")
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	var product data.Product

	result := h.DB.Where("id=?", productId).Preload("SKUs.SKUVariants").Preload("SKUs.Variants").Find(&product)
	if result.Error != nil {
		utils.ErrorJSON(w, result.Error)
		return
	}
	if result.RowsAffected == 0 {
		utils.ErrorJSON(w, errors.New("Product Not Found"), 404)
		return
	}

	productDetailsResponse := product.ToProductDetailsResponse()
	utils.WriteJSON(w, 200, productDetailsResponse)
}
