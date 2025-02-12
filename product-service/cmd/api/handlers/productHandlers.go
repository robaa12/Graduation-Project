package handlers

import (
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/robaa12/product-service/cmd/data"
	"github.com/robaa12/product-service/cmd/utils"
	"github.com/robaa12/product-service/cmd/validation"
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

	if err := validation.ValidateNewProduct(productRequest); err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	if err := validation.ValidateBusinessRules(productRequest); err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	slug, err := utils.ValidateAndGenerateSlug(h.DB, productRequest.Name, productRequest.StoreID)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}
	productRequest.Slug = slug
	// Create a new Product
	var product data.Product
	product.CreateProduct(productRequest)

	// Start Database Transaction
	err = h.DB.Transaction(func(tx *gorm.DB) error {

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
	productResponse := product.ToProductResponse()
	// Return the product
	utils.WriteJSON(w, 201, productResponse)
}

// GetProduct returns a product from the database
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	var product data.Product

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
	var product data.ProductResponse
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
	var dbProduct data.Product
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
	err = h.DB.Model(&data.Product{}).Where("id = ?", id).Updates(&product).Error
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

	// parse pagination parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	// Vaildate sorting field
	validSortFields := map[string]bool{
		"created_at": true,
		"name":       true,
		"price":      true,
	}
	if sort != "" && !validSortFields[sort] {
		utils.ErrorJSON(w, errors.New("invalid sort field"), http.StatusBadRequest)
		return
	}

	// Create pagination query
	pagination := utils.NewPaginationQuery(page, pageSize, sort, order)

	// Get total count
	var total int64
	if err := h.DB.Model(&data.Product{}).Where("store_id = ?", storeID).Count(&total).Error; err != nil {
		utils.ErrorJSON(w, err)
	}

	// Calculate offset
	offset := (pagination.Page - 1) * pagination.PageSize

	// Create base Query
	query := h.DB.Model(&data.Product{}).Where("store_id", storeID)

	// Add sorting
	if pagination.Sort != "" {
		query = query.Order(fmt.Sprintf("%s %s", pagination.Sort, pagination.Order))
	}

	// Execute paginated query
	var products []data.Product
	if err := query.Offset(offset).Limit(pagination.PageSize).Find(&products).Error; err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	// If no products found
	if len(products) == 0 {
		utils.ErrorJSON(w, errors.New("no products found"), http.StatusNotFound)
		return
	}

	// Create Slice of ProductResponse
	var productsResponse []data.ProductResponse
	for _, product := range products {
		productResponse := product.ToProductResponse()
		productsResponse = append(productsResponse, productResponse)
	}
	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(pagination.PageSize)))

	// Create paginated response
	response := utils.PaginatedResponse{
		Data:       productsResponse,
		Total:      total,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: totalPages,
	}
	// Return store's Products
	utils.WriteJSON(w, 200, response)
}

// GetProductDetails returns a product details from the database
func (h *ProductHandler) GetProductDetails(w http.ResponseWriter, r *http.Request) {
	// Get the product ID from the URL
	productId, err := utils.GetID(r, "product_id")
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

func (h *ProductHandler) GetProductBySlug(w http.ResponseWriter, r *http.Request) {
	store_id := chi.URLParam(r, "store_id")
	slug := chi.URLParam(r, "slug")

	if slug == "" || store_id == "" {
		utils.ErrorJSON(w, errors.New("both slug and store_id is required"), http.StatusBadRequest)
		return
	}
	var product data.Product
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
