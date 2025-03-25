package handlers

import (
	"net/http"

	apperrors "github.com/robaa12/product-service/cmd/errors"
	"github.com/robaa12/product-service/cmd/model"
	"github.com/robaa12/product-service/cmd/service"
	"github.com/robaa12/product-service/cmd/utils"
)

type CollectionHandler struct {
	service *service.CollectionService
}

func NewCollectionHandler(service *service.CollectionService) *CollectionHandler {
	return &CollectionHandler{service: service}
}
func (h *CollectionHandler) CreateCollection(w http.ResponseWriter, r *http.Request) {
	storeID, err := utils.GetID(r, "store_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid store ID"))
		return
	}
	var collectionRequest model.CollectionRequest
	if err := utils.ReadJSON(w, r, &collectionRequest); err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid request payload"))
		return
	}

	collectionResponse, err := h.service.CreateCollection(storeID, &collectionRequest)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}
	_ = utils.WriteJSON(w, http.StatusCreated, collectionResponse)
}

// GetCollections - GET /stores/{store_id}/collections/
func (h *CollectionHandler) GetCollections(w http.ResponseWriter, r *http.Request) {
	storeID, err := utils.GetID(r, "store_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid store ID"))
		return
	}
	collections, err := h.service.GetCollections(storeID)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	_ = utils.WriteJSON(w, http.StatusOK, collections)
}

// GetCollection - GET /stores/{store_id}/collections/{collection_id}
func (h *CollectionHandler) GetCollection(w http.ResponseWriter, r *http.Request) {
	storeID, err := utils.GetID(r, "store_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid store ID"))
		return
	}
	collectionID, err := utils.GetID(r, "collection_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid collection_id"))
		return
	}
	collection, err := h.service.GetCollection(storeID, collectionID)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	_ = utils.WriteJSON(w, http.StatusOK, collection)
}

// AddProductToCollection Add product to collection - POST /stores/{store_id}/collections/{collection_id}
func (h *CollectionHandler) AddProductToCollection(w http.ResponseWriter, r *http.Request) {
	storeID, err := utils.GetID(r, "store_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid store ID"))
		return
	}
	collectionID, err := utils.GetID(r, "collection_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid collection_id"))
		return
	}

	var collectionProductsRequest model.CollectionProductsRequest
	// Read product IDs from request body
	if err := utils.ReadJSON(w, r, &collectionProductsRequest); err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid request payload"))
		return
	}
	err = h.service.AddProductToCollection(storeID, collectionID, &collectionProductsRequest)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	_ = utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Products added to collection successfully",
	})
}

// RemoveProductFromCollection Remove product from collection - DELETE /stores/{store_id}/collections/{collection_id}/products/{product_id}
func (h *CollectionHandler) RemoveProductFromCollection(w http.ResponseWriter, r *http.Request) {
	storeID, err := utils.GetID(r, "store_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid store ID"))
		return
	}
	collectionID, err := utils.GetID(r, "collection_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid collection_id"))
		return
	}
	productID, err := utils.GetID(r, "product_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid product_id"))
		return
	}
	err = h.service.RemoveProductFromCollection(storeID, collectionID, productID)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}
	_ = utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Product removed from collection successfully",
	})
}
