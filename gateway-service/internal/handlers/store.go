package handlers

import (
	"fmt"
	"log"
	"net/http"

	apperrors "github.com/robaa12/gatway-service/internal/errors"
	"github.com/robaa12/gatway-service/internal/middleware/auth"
	"github.com/robaa12/gatway-service/internal/model"
	"github.com/robaa12/gatway-service/internal/service"
	"github.com/robaa12/gatway-service/utils"
)

// StoreHandler handles store-related operations
type StoreHandler struct {
	storeService *service.StoreService
	jwtService   *auth.JWTService
}

// NewStoreHandler creates a new store handler
func NewStoreHandler(storeService *service.StoreService, jwtService *auth.JWTService) *StoreHandler {
	return &StoreHandler{
		storeService: storeService,
		jwtService:   jwtService,
	}
}

// CreateStore handles the distributed transaction for creating a store across all services
func (h *StoreHandler) CreateStore(w http.ResponseWriter, r *http.Request) {

	// Extract user claims from context (set by auth middleware)
	claims, ok := r.Context().Value("user").(*auth.Claims)
	if !ok {
		utils.ErrorJSON(w, fmt.Errorf("unauthorized"), http.StatusUnauthorized)
		return
	}

	// Read the original request body
	var storeRequest model.StoreRequest
	err := utils.ReadJSON(w, r, &storeRequest)
	if err != nil {
		utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid request payload"))
		return
	}

	// Add user ID from token to request
	storeRequest.UserID = uint(claims.UserID)

	store, err := h.storeService.CreateStore(&storeRequest)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	//  Generate new tokens with updated store IDs
	tokenResponse, err := h.jwtService.GenerateUpdatedTokenResponse(claims.UserID, claims.StoresID, store.ID)
	if err != nil {
		log.Printf("Error generating new tokens: %v", err)
	}
	storeResponse := store.GetStoreResponse(tokenResponse)
	// Return the enhanced response
	utils.WriteJSON(w, http.StatusCreated, storeResponse)
}
