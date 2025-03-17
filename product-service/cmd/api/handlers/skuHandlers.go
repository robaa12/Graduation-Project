package handlers

import (
	"net/http"

	apperrors "github.com/robaa12/product-service/cmd/errors"
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
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid skuID"))
		return
	}
	storeID, err := utils.GetID(r, "store_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid storeID"))
		return
	}
	productID, err := utils.GetID(r, "product_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid productID"))
		return
	}

	var skuRequest model.SKURequest
	err = utils.ReadJSON(w, r, &skuRequest)
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid request payload"))
		return
	}
	err = h.service.UpdateSKU(skuID, productID, storeID, &skuRequest)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}
	_ = utils.WriteJSON(w, http.StatusOK, "SKU updated successfully.")
}

func (h *SKUHandler) GetSKU(w http.ResponseWriter, r *http.Request) {
	// Get the SKU ID from the URL
	skuID, err := utils.GetID(r, "sku_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid skuID"))
		return
	}
	storeID, err := utils.GetID(r, "store_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid storeID"))
		return
	}
	productID, err := utils.GetID(r, "product_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid productID"))
		return
	}

	skuResponse, err := h.service.GetSKU(skuID, productID, storeID)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}
	_ = utils.WriteJSON(w, http.StatusOK, skuResponse)
}

func (h *SKUHandler) DeleteSKU(w http.ResponseWriter, r *http.Request) {
	// Get the SKU ID from the URL
	skuID, err := utils.GetID(r, "sku_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid skuID"))
		return
	}
	storeID, err := utils.GetID(r, "store_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid storeID"))
		return
	}
	productID, err := utils.GetID(r, "product_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid productID"))
		return
	}

	// Find SKU by ID
	err = h.service.DeleteSKU(skuID, productID, storeID)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}
	_ = utils.WriteJSON(w, http.StatusOK, "SKU deleted successfully.")
}

func (h *SKUHandler) NewSKU(w http.ResponseWriter, r *http.Request) {
	storeID, err := utils.GetID(r, "store_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid storeID"))
		return
	}
	productID, err := utils.GetID(r, "product_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid productID"))
		return
	}
	var skuRequest model.SKURequest
	err = utils.ReadJSON(w, r, &skuRequest)
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid request payload"))
		return
	}

	// Start Database Transaction
	// Create a new SKU
	skuResponse, err := h.service.NewSKU(storeID, productID, &skuRequest)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	// Return the SKU
	_ = utils.WriteJSON(w, http.StatusCreated, skuResponse)
}
