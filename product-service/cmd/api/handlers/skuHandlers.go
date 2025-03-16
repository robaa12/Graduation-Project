package handlers

import (
	"errors"
	"net/http"

	"github.com/robaa12/product-service/cmd/model"
	"github.com/robaa12/product-service/cmd/service"
	"github.com/robaa12/product-service/cmd/utils"
)

type SKUHandler struct {
	service *service.SKUService
}

func NewSKUHandler(s *service.SKUService) *SKUHandler {
	return &SKUHandler{service: s}
}

// GetStoreProducts returns all products of a store

func (h *SKUHandler) UpdateSKU(w http.ResponseWriter, r *http.Request) {
	// Get the SKU ID from the URL
	skuID, err := utils.GetID(r, "sku_id")
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	var skuRequest model.SKURequest
	err = utils.ReadJSON(w, r, &skuRequest)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}
	err = h.service.UpdateSKU(int(skuID), &skuRequest)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
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
	skuResponse, err := h.service.GetSKU(int(skuID))
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
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
	err = h.service.DeleteSKU(int(skuID))
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	utils.WriteJSON(w, 200, "SKU deleted successfully.")
}

func (h *SKUHandler) NewSKU(w http.ResponseWriter, r *http.Request) {
	// Get store id from URI
	storeID, err := utils.GetID(r, "store_id")
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}
	// Read the JSON request
	var skuRequest model.SKURequest
	err = utils.ReadJSON(w, r, &skuRequest)
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
	skuResonse, err := h.service.NewSKU(storeID, productID, &skuRequest)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Return the SKU
	utils.WriteJSON(w, 201, skuResonse)
}
