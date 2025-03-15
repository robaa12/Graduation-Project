package handlers

import (
	"errors"
	"net/http"

	"github.com/robaa12/product-service/cmd/model"
	"github.com/robaa12/product-service/cmd/utils"
	"github.com/robaa12/product-service/cmd/validation"
	"gorm.io/gorm"
)

type CollectionHandler struct {
	DB        *gorm.DB
	Validator *validation.CollectionValidator
}

func (h *CollectionHandler) CreateCollection(w http.ResponseWriter, r *http.Request) {
	storeID, err := utils.GetID(r, "store_id")
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	var collectionRequest model.CollectionRequest
	if err := utils.ReadJSON(w, r, &collectionRequest); err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Validate request
	if err := h.Validator.ValidateCollectionRequest(&collectionRequest); err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	collection := model.Collection{
		StoreID:     storeID,
		Name:        collectionRequest.Name,
		Description: collectionRequest.Description,
	}

	if err := h.DB.Create(&collection).Error; err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, collection.ToCollectionResponse())
}

// GetCollections - GET /stores/{store_id}/collections/
func (h *CollectionHandler) GetCollections(w http.ResponseWriter, r *http.Request) {
	StoreID, err := utils.GetID(r, "store_id")
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	var collections []model.Collection
	// Get collections from the database by store ID
	if err := h.DB.Where("store_id = ?", StoreID).Find(&collections).Error; err != nil {
		utils.ErrorJSON(w, err)
		return
	}
	// Convert to response
	var output []model.CollectionResponse
	for _, collection := range collections {
		output = append(output, *collection.ToCollectionResponse())
	}
	utils.WriteJSON(w, 200, output)
}

// GetCollection - GET /stores/{store_id}/collections/{collection_id}
func (h *CollectionHandler) GetCollection(w http.ResponseWriter, r *http.Request) {
	StoreID, err := utils.GetID(r, "store_id")
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}
	CollectionID, err := utils.GetID(r, "collection_id")
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	// Get collection from the database by store ID and collection ID
	collection, err := h.Validator.GetCollectionWithProducts(CollectionID, StoreID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ErrorJSON(w, errors.New("collection not found"), http.StatusNotFound)
		} else {
			utils.ErrorJSON(w, err)
		}
		return
	}
	utils.WriteJSON(w, 200, collection.ToCollectionDetailsResponse())
}

// Add product to collection - POST /stores/{store_id}/collections/{collection_id}
func (h *CollectionHandler) AddProductToCollection(w http.ResponseWriter, r *http.Request) {
	storeID, err := utils.GetID(r, "store_id")
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	collectionID, err := utils.GetID(r, "collection_id")
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	// Read product IDs from request body
	var request struct {
		ProductIDs []uint `json:"product_ids"`
	}
	if err := utils.ReadJSON(w, r, &request); err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Validate collection
	collection, err := h.Validator.CollectionExists(collectionID, storeID)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusNotFound)
		return
	}

	// Validate products
	products, err := h.Validator.ValidateProductsExist(storeID, request.ProductIDs)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Add products to collection
	if err := h.DB.Model(collection).Association("Products").Append(&products); err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Products added to collection successfully",
	})
}

// Remove product from collection - DELETE /stores/{store_id}/collections/{collection_id}/products/{product_id}
func (h *CollectionHandler) RemoveProductFromCollection(w http.ResponseWriter, r *http.Request) {
	storeID, err := utils.GetID(r, "store_id")
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	collectionID, err := utils.GetID(r, "collection_id")
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	productID, err := utils.GetID(r, "product_id")
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	// Validate collection
	collection, err := h.Validator.CollectionExists(collectionID, storeID)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusNotFound)
		return
	}

	// Validate product
	product, err := h.Validator.ValidateProductExists(productID, storeID)
	if err != nil {
		utils.ErrorJSON(w, err, http.StatusNotFound)
		return
	}

	// Remove product from collection
	if err := h.DB.Model(collection).Association("Products").Delete(product); err != nil {
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Product removed from collection successfully",
	})
}
