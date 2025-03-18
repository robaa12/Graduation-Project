package handlers

import (
	"net/http"

	apperrors "github.com/robaa12/product-service/cmd/errors"
	"github.com/robaa12/product-service/cmd/model"
	"github.com/robaa12/product-service/cmd/service"
	"github.com/robaa12/product-service/cmd/utils"
)

type ReviewHandler struct {
	service *service.ReviewService
}

func NewReviewHandler(service *service.ReviewService) *ReviewHandler {
	return &ReviewHandler{service: service}
}

// CreateReview - POST /products/{productID}/stores/{storeID}/reviews
func (h *ReviewHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	storeID, err := utils.GetID(r, "store_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid store id"))
		return
	}

	productID, err := utils.GetID(r, "product_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid product id"))
		return
	}

	var reviewRequest model.ReviewRequest
	if err := utils.ReadJSON(w, r, &reviewRequest); err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid request payload"))
		return
	}

	reviewResponse, err := h.service.CreateReview(productID, storeID, &reviewRequest)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}
	_ = utils.WriteJSON(w, http.StatusCreated, reviewResponse)
}

// GetProductReviews - GET /products/{productID}/stores/{storeID}/reviews
func (h *ReviewHandler) GetProductReviews(w http.ResponseWriter, r *http.Request) {
	storeID, err := utils.GetID(r, "store_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid store id"))
		return
	}

	productID, err := utils.GetID(r, "product_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid product id"))
		return
	}

	reviewResponses, err := h.service.GetProductReviews(productID, storeID)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}
	_ = utils.WriteJSON(w, http.StatusOK, reviewResponses)
}

// GetReview - GET /products/{productID}/stores/{storeID}/reviews/{reviewID}
func (h *ReviewHandler) GetReview(w http.ResponseWriter, r *http.Request) {
	storeID, err := utils.GetID(r, "store_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid store id"))
		return
	}
	productID, err := utils.GetID(r, "product_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid product id"))
		return
	}
	reviewID, err := utils.GetID(r, "review_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid review id"))
		return
	}

	review, err := h.service.GetReview(reviewID, productID, storeID)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}
	_ = utils.WriteJSON(w, http.StatusOK, review)
}

// GetReviewStatistics - GET /products/{productID}/stores/{storeID}/reviews/statistics
func (h *ReviewHandler) GetReviewStatistics(w http.ResponseWriter, r *http.Request) {
	storeID, err := utils.GetID(r, "store_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid store ID"))
		return
	}

	productID, err := utils.GetID(r, "product_id")
	if err != nil {
		_ = utils.ErrorJSON(w, apperrors.NewBadRequestError("invalid product ID"))
		return
	}

	statistics, err := h.service.GetReviewStatistics(productID, storeID)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	_ = utils.WriteJSON(w, http.StatusOK, statistics)
}
