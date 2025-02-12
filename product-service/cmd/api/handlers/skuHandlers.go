package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/robaa12/product-service/cmd/data"
	"github.com/robaa12/product-service/cmd/utils"
	"gorm.io/gorm"
)

type SKUHandler struct {
	DB *gorm.DB
}

// GetStoreProducts returns all products of a store

func (h *SKUHandler) UpdateSKU(w http.ResponseWriter, r *http.Request) {
	// Get the SKU ID from the URL
	skuID, err := utils.GetID(r, "id")
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	var skuRequest data.SKURequest
	err = utils.ReadJSON(w, r, &skuRequest)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	result := h.DB.Where("id= ?", skuID).Updates(&data.Sku{
		Stock:          skuRequest.Stock,
		Price:          skuRequest.Price,
		CompareAtPrice: skuRequest.CompareAtPrice,
		CostPerItem:    skuRequest.CostPerItem,
		Profit:         skuRequest.Profit,
		Margin:         skuRequest.Margin,
	})
	if result.Error != nil {
		utils.ErrorJSON(w, result.Error)
		return
	}

	utils.WriteJSON(w, 200, "SKU updated successfully.")
}

func (h *SKUHandler) GetSKU(w http.ResponseWriter, r *http.Request) {
	// Get the SKU ID from the URL
	skuID, err := utils.GetID(r, "sku_id")
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	// Find SKU by ID
	var sku data.Sku
	result := h.DB.Model(&data.Sku{}).Where("id = ?", skuID).Preload("Variants").Preload("SKUVariants").Find(&sku)
	if result.RowsAffected == 0 {
		utils.ErrorJSON(w, errors.New("SKU Not Found"), 404)
		return
	}
	if result.Error != nil {
		utils.ErrorJSON(w, result.Error)
		return
	}
	skuResponse := sku.ToSKUResponse()
	utils.WriteJSON(w, 200, skuResponse)
}

func (h *SKUHandler) DeleteSKU(w http.ResponseWriter, r *http.Request) {
	// Get the SKU ID from the URL
	skuID, err := utils.GetID(r, "sku_id")
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	// Find SKU by ID
	var sku data.Sku
	result := h.DB.Where("id = ?", skuID).Find(&sku)
	if result.RowsAffected == 0 {
		utils.ErrorJSON(w, errors.New("Sku is not found."), 404)
		return
	}
	if result.Error != nil {
		utils.ErrorJSON(w, result.Error)
		return
	}

	err = h.DB.Unscoped().Delete(&sku).Error
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}
	utils.WriteJSON(w, 200, "SKU deleted successfully.")
}

func (h *SKUHandler) NewSKU(w http.ResponseWriter, r *http.Request) {
	// Read the JSON request
	var skuRequest data.SKURequest
	err := utils.ReadJSON(w, r, &skuRequest)
	if err != nil {
		utils.ErrorJSON(w, errors.New("Enter valid SKU data"))
		return
	}
	// Get the Product ID from the URL
	productID, err := utils.GetID(r, "product_id")
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	// Start Database Transaction
	// Create a new SKU
	var sku data.Sku
	sku.CreateSKU(skuRequest, productID)

	// make transaction
	tx := h.DB.Begin()
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
		utils.ErrorJSON(w, err)
		return
	}
	// Return the SKU
	utils.WriteJSON(w, 201, "SKU created successfully.")
}

func (h *SKUHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	var products []data.Product
	err := h.DB.Preload("SKUs.Variants").Preload("SKUs.SKUVariants").Find(&products).Error
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	utils.WriteJSON(w, 200, products)
}
