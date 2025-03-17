package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/robaa12/product-service/cmd/model"
	"github.com/robaa12/product-service/cmd/service"
	"github.com/robaa12/product-service/cmd/utils"
	"gorm.io/gorm"
)

type ProductHandler struct {
	ProductService service.ProductService
	DB             *gorm.DB
}

// NewProduct creates a new product , skus and variants in the database
func (h *ProductHandler) NewProduct(w http.ResponseWriter, r *http.Request) {
	// Read the JSON request
	var productRequest model.ProductRequest
	err := utils.ReadJSON(w, r, &productRequest)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}
	productResponse, err := h.ProductService.NewProduct(productRequest)

	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}
	// Return the product
	utils.WriteJSON(w, 201, productResponse)
}

// GetProduct returns a product from the database
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	// Get the product ID from the URL
	id, err := utils.GetID(r, "product_id")
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}
	// Get the product from the database
	productResponse, err := h.ProductService.GetProduct(id)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}
	// Return the product
	utils.WriteJSON(w, 200, productResponse)
}

// UpdateProduct updates a product in the database
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	// Read the JSON request
	var product model.ProductResponse
	err := utils.ReadJSON(w, r, &product)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	// Get the product ID from the URL
	id, err := utils.GetID(r, "product_id")
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	err = h.ProductService.UpdateProduct(id, product)
	product.ID = id
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	// Return the Updated product
	utils.WriteJSON(w, http.StatusOK, product)
}

// DeleteProduct deletes a product from the database
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	// Get the product ID from the URL
	id, err := utils.GetID(r, "product_id")
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	// Get the store ID from the URL
	storeID, err := utils.GetID(r, "store_id")
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	// Call the service to delete the product
	err = h.ProductService.DeleteProduct(id, storeID)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusNotFound)
		return
	}

	// Return success message
	utils.WriteJSON(w, http.StatusOK, "Product deleted successfully")
}

func (h *ProductHandler) GetStoreProducts(w http.ResponseWriter, r *http.Request) {
	// Fetch Store ID Param From URL
	storeID, err := utils.GetID(r, "store_id")
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	productsResponse, err := h.ProductService.GetStoreProducts(storeID)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusNotFound)
	}

	// Return store's Products
	utils.WriteJSON(w, 200, productsResponse)
}

// GetProductDetails returns a product details from the database
func (h *ProductHandler) GetProductDetails(w http.ResponseWriter, r *http.Request) {
	// Get the product ID from the URL
	productId, err := utils.GetID(r, "product_id")
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	productDetailsResponse, err := h.ProductService.GetProductDetails(productId)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusNotFound)
		return
	}
	utils.WriteJSON(w, http.StatusOK, productDetailsResponse)
}

func (h *ProductHandler) GetProductBySlug(w http.ResponseWriter, r *http.Request) {
	storeIDStr := chi.URLParam(r, "store_id")
	slug := chi.URLParam(r, "slug")

	if slug == "" || storeIDStr == "" {
		utils.ErrorJSON(w, errors.New("both slug and store_id is required"), http.StatusBadRequest)
		return
	}
	StoreID, err := strconv.ParseUint(storeIDStr, 10, 0)
	if err != nil {
		utils.ErrorJSON(w, errors.New("invalid store ID format"), http.StatusBadRequest)
		return
	}
	productDetailsResponse, err := h.ProductService.GetProductBySlug(slug, uint(StoreID))
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusNotFound)
		return
	}
	utils.WriteJSON(w, http.StatusOK, productDetailsResponse)
}
