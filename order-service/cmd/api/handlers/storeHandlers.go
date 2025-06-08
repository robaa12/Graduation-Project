package handlers

import (
	"log"
	"net/http"
	apperrors "order-service/cmd/errors"
	"order-service/cmd/model"
	"order-service/cmd/service"
	"order-service/cmd/utils"
)

type StoreHandler struct {
	service *service.StoreService
}

func NewStoreHandler(s *service.StoreService) *StoreHandler {
	return &StoreHandler{service: s}
}

// CreateStore creates a new store
func (h *StoreHandler) CreateStore(w http.ResponseWriter, r *http.Request) {
	// Parse the request body into a StoreRequest object
	var storeRequest model.StoreRequest
	err := utils.ReadJSON(w, r, &storeRequest)
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid request payload"))
		return
	}

	// Call the CreateStore method from the service layer
	storeResponse, err := h.service.CreateStore(&storeRequest)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}
	log.Printf("Store request created successfully: %+v", storeRequest)
	log.Printf("Store created successfully: %+v", storeResponse)
	// Return the created store response
	_ = utils.WriteJSON(w, http.StatusCreated, storeResponse)

}

// DeleteStore deletes a store by ID
func (h *StoreHandler) DeleteStore(w http.ResponseWriter, r *http.Request) {
	// Get the store ID from the URL parameters
	storeID, err := utils.GetID(r, "store_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid store id"))
		return
	}
	// Call the DeleteStore method from the service layer
	err = h.service.DeleteStore(storeID)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}
	log.Printf("Store deleted successfully")
	// Return a success response
	_ = utils.WriteJSON(w, http.StatusNoContent, nil)
}
