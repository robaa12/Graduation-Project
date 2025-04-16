package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/robaa12/product-service/cmd/model"
	"github.com/robaa12/product-service/cmd/utils"
	"gorm.io/gorm"
)

type OrderHandler struct {
	DB *gorm.DB
}

type VerificationRequest struct {
	StoreID uint               `json:"store_id"`
	Items   []VerificationItem `json:"items"`
}

type VerificationItem struct {
	SkuID    uint    `json:"sku_id"`
	Quantity uint    `json:"quantity"`
	Price    float64 `json:"price"`
}

type VerificationResponse struct {
	Valid   bool           `json:"valid"`
	Message string         `json:"messages"`
	Items   []VerifiedItem `json:"items"`
}

type VerifiedItem struct {
	SkuID   uint     `json:"sku_id"`
	Valid   bool     `json:"valid"`
	InStock bool     `json:"in_stock"`
	Price   float64  `json:"actual_price"`
	Message []string `json:"message,omitempty"`
}

func (h *OrderHandler) VerifyOrderItems(w http.ResponseWriter, r *http.Request) {
	var req VerificationRequest
	if err := utils.ReadJSON(w, r, &req); err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	response := VerificationResponse{
		Valid: true,
		Items: []VerifiedItem{},
	}

	// Extract all SKU IDs
	skuIDs := make([]uint, len(req.Items))
	skuQuantityMap := make(map[uint]uint)
	skuPriceMap := make(map[uint]float64)

	for i, item := range req.Items {
		skuIDs[i] = item.SkuID
		skuQuantityMap[item.SkuID] = item.Quantity
		skuPriceMap[item.SkuID] = item.Price
	}

	// Fetch all SKUs in a single query
	var skus []model.Sku
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		// Fix: Properly specify the join condition and table references
		if err := tx.Joins("JOIN products ON skus.product_id = products.id").
			Where("skus.id IN ? AND products.store_id = ?", skuIDs, req.StoreID).
			Find(&skus).Error; err != nil {
			return err
		}

		// Create a map for quick SKU lookup
		skuMap := make(map[uint]model.Sku)
		for _, sku := range skus {
			skuMap[sku.ID] = sku
		}

		// Verify each requested SKU
		for _, requestedSkuID := range skuIDs {
			sku, exists := skuMap[requestedSkuID]
			verifiedItem := VerifiedItem{
				SkuID:   requestedSkuID,
				Valid:   true,
				InStock: true,
			}

			if !exists {
				response.Valid = false
				verifiedItem.Valid = false
				verifiedItem.Message = append(verifiedItem.Message, "SKU not found or does not belong to store")
			} else {
				requestedQty := skuQuantityMap[requestedSkuID]
				requestedPrice := skuPriceMap[requestedSkuID]

				// Verify stock
				if sku.Stock < int(requestedQty) {
					response.Valid = false
					verifiedItem.Valid = false
					verifiedItem.InStock = false
					verifiedItem.Message = append(verifiedItem.Message, fmt.Sprintf("Insufficient stock (available: %d)", sku.Stock))
				}

				// Verify price
				if requestedPrice != sku.Price {
					response.Valid = false
					verifiedItem.Valid = false
					verifiedItem.Message = append(verifiedItem.Message, fmt.Sprintf("Price mismatch (actual: %.2f)", sku.Price))
				}

				verifiedItem.Price = sku.Price
			}

			response.Items = append(response.Items, verifiedItem)
		}

		return nil
	})

	if err != nil {
		_ = utils.ErrorJSON(w, err)
		log.Println(err)
		return
	}
	if !response.Valid {
		response.Message = "Unverified Order Items"
	} else {
		response.Message = "Verified Order Items"
	}

	_ = utils.WriteJSON(w, http.StatusOK, response)
}

func (h *OrderHandler) UpdateInventory(w http.ResponseWriter, r *http.Request) {

	var orderItems []VerificationItem

	if err := utils.ReadJSON(w, r, &orderItems); err != nil {
		log.Println(err)
		_ = utils.ErrorJSON(w, err, http.StatusBadRequest)
	}

	var skus []model.Sku
	// Populate skus
	for _, item := range orderItems {
		sku := model.Sku{
			ID:    item.SkuID,
			Stock: int(item.Quantity),
		}
		skus = append(skus, sku)
	}

	if err := model.UpdateInventory(h.DB, skus); err != nil {
		_ = utils.ErrorJSON(w, errors.New("error Updating Inventory"), http.StatusNotFound)
	}

	_ = utils.WriteJSON(w, 200, "Inventory Updated Successfully.")
}
