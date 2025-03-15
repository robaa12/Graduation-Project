package handlers

import (
	"errors"
	"fmt"
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
	var product model.Product

	// Get the product ID from the URL
	id, err := utils.GetID(r, "product_id")
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

	// Check if the product exists
	var dbProduct model.Product
	err = h.DB.Where("id = ?", id).First(&dbProduct).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ErrorJSON(w, errors.New("product not found"), http.StatusNotFound)
		} else {
			utils.ErrorJSON(w, err, http.StatusInternalServerError)
		}
	}

	// Compare product names
	if dbProduct.Name != product.Name {
		slug, err := utils.ValidateAndGenerateSlug(h.DB, product.Name, dbProduct.StoreID)
		if err != nil {
			utils.ErrorJSON(w, err, http.StatusInternalServerError)
			return
		}
		fmt.Println("slug: ", slug)
		product.Slug = slug
	}

	// Update the product in the database
	err = h.DB.Model(&model.Product{}).Where("id = ?", id).Updates(&product).Error
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}
	product.ID = id
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
	var product model.Product
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

	// Create base Query
	query := h.DB.Model(&model.Product{}).Where("store_id", storeID)
	// Execute Query
	var products []model.Product
	if err := query.Find(&products).Error; err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	// If no products found
	if len(products) == 0 {
		utils.ErrorJSON(w, errors.New("no products found"), http.StatusNotFound)
		return
	}

	// Create Slice of ProductResponse
	var productsResponse []model.ProductResponse
	for _, product := range products {
		productResponse := product.ToProductResponse()
		productsResponse = append(productsResponse, *productResponse)
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

	var product model.Product

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

func (h *ProductHandler) GetProductBySlug(w http.ResponseWriter, r *http.Request) {
	store_id := chi.URLParam(r, "store_id")
	slug := chi.URLParam(r, "slug")

	if slug == "" || store_id == "" {
		utils.ErrorJSON(w, errors.New("both slug and store_id is required"), http.StatusBadRequest)
		return
	}
	var product model.Product
	result := h.DB.Where("slug = ? AND store_id = ?", slug, store_id).Preload("SKUs.SKUVariants").Preload("SKUs.Variants").First(&product)
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
