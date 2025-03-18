package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	apperrors "github.com/robaa12/product-service/cmd/errors"
	"github.com/robaa12/product-service/cmd/model"
	"github.com/robaa12/product-service/cmd/service"
	"github.com/robaa12/product-service/cmd/utils"
)

type ProductHandler struct {
	ProductService service.ProductService
}

// NewProduct creates a new product , skus and variants in the database
func (h *ProductHandler) NewProduct(w http.ResponseWriter, r *http.Request) {
	// Get the storeID  from the URL
	storeID, err := utils.GetID(r, "store_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid store ID"))
		return
	}
	// Read the JSON request
	var productRequest model.ProductRequest
	err = utils.ReadJSON(w, r, &productRequest)
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid request payload"))
		return
	}
	productResponse, err := h.ProductService.NewProduct(storeID, productRequest)

	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}
	// Return the product
	_ = utils.WriteJSON(w, http.StatusCreated, productResponse)
}

// GetProduct returns a product from the database
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	// Get the product ID from the URL
	storeID, err := utils.GetID(r, "store_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid store ID"))
		return
	}
	id, err := utils.GetID(r, "product_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid product ID"))
		return
	}
	// Get the product from the database
	productResponse, err := h.ProductService.GetProduct(id, storeID)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}
	// Return the product
	_ = utils.WriteJSON(w, http.StatusOK, productResponse)
}

// UpdateProduct updates a product in the database
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	// Read the JSON request
	var product model.ProductResponse
	err := utils.ReadJSON(w, r, &product)
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid request payload"))
		return
	}

	// Get the product ID from the URL
	id, err := utils.GetID(r, "product_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid product ID"))
		return
	}
	storeID, err := utils.GetID(r, "store_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid store ID"))
		return
	}

	err = h.ProductService.UpdateProduct(id, storeID, product)
	product.ID = id
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	// Return the Updated product
	_ = utils.WriteJSON(w, http.StatusOK, product)
}

// DeleteProduct deletes a product from the database
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	// Get the product ID from the URL
	id, err := utils.GetID(r, "product_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid product ID"))
		return
	}

	// Get the store ID from the URL
	storeID, err := utils.GetID(r, "store_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid store ID"))
		return
	}

	// Call the service to delete the product
	err = h.ProductService.DeleteProduct(id, storeID)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	// Return success message
	_ = utils.WriteJSON(w, http.StatusOK, "Product deleted successfully")
}

func (h *ProductHandler) GetStoreProducts(w http.ResponseWriter, r *http.Request) {
	// Fetch Store ID Param From URL
	storeID, err := utils.GetID(r, "store_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid store ID"))
		return
	}

	productsResponse, err := h.ProductService.GetStoreProducts(storeID)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	// Return store's Products
	_ = utils.WriteJSON(w, http.StatusOK, productsResponse)
}

// GetProductDetails returns a product details from the database
func (h *ProductHandler) GetProductDetails(w http.ResponseWriter, r *http.Request) {
	// Get the product ID from the URL
	productId, err := utils.GetID(r, "product_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid product ID"))
		return
	}
	storeId, err := utils.GetID(r, "store_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid store ID"))
		return
	}

	productDetailsResponse, err := h.ProductService.GetProductDetails(productId, storeId)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}
	_ = utils.WriteJSON(w, http.StatusOK, productDetailsResponse)
}

func (h *ProductHandler) GetProductBySlug(w http.ResponseWriter, r *http.Request) {
	storeIDStr := chi.URLParam(r, "store_id")
	slug := chi.URLParam(r, "slug")

	if slug == "" || storeIDStr == "" {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("Store ID and slug is required"), http.StatusBadRequest)
		return
	}
	StoreID, err := strconv.ParseUint(storeIDStr, 10, 0)
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("Invalid store ID"), http.StatusBadRequest)
		return
	}
	productDetailsResponse, err := h.ProductService.GetProductBySlug(slug, uint(StoreID))
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewNotFoundError("Product not found"), http.StatusNotFound)
		return
	}
	_ = utils.WriteJSON(w, http.StatusOK, productDetailsResponse)
}
