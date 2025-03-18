package service

import (
	apperrors "github.com/robaa12/product-service/cmd/errors"
	"github.com/robaa12/product-service/cmd/model"
	"github.com/robaa12/product-service/cmd/repository"
)

type ReviewService struct {
	reviewRepo *repository.ReviewRepository
}

func NewReviewService(reviewRepo *repository.ReviewRepository) *ReviewService {
	return &ReviewService{reviewRepo: reviewRepo}
}

func (rs *ReviewService) CreateReview(productID, storeID uint, reviewRequest *model.ReviewRequest) (*model.ReviewResponse, error) {
	// Check if product exists
	exists, err := rs.reviewRepo.ProductExists(productID, storeID)
	if err != nil {
		return nil, apperrors.ErrCheck(err)
	}
	if !exists {
		return nil, apperrors.NewNotFoundError("product not found")
	}

	// Create new review
	review := reviewRequest.ToReview(productID, storeID)

	err = rs.reviewRepo.CreateReview(review)
	if err != nil {
		return nil, apperrors.ErrCheck(err)
	}
	return review.ToReviewResponse(), nil
}

func (rs *ReviewService) GetProductReviews(productID, storeID uint) ([]model.ReviewResponse, error) {
	exists, err := rs.reviewRepo.ProductExists(productID, storeID)
	if err != nil {
		return nil, apperrors.ErrCheck(err)
	}
	if !exists {
		return nil, apperrors.NewNotFoundError("product not found")
	}

	reviews, err := rs.reviewRepo.GetProductReviews(productID, storeID)
	if err != nil {
		return nil, apperrors.ErrCheck(err)
	}
	var reviewResponses []model.ReviewResponse
	for _, review := range reviews {
		reviewResponses = append(reviewResponses, *review.ToReviewResponse())
	}
	return reviewResponses, nil
}

func (rs *ReviewService) GetReview(reviewID, productID, storeID uint) (*model.ReviewResponse, error) {
	review, err := rs.reviewRepo.GetReview(reviewID, productID, storeID)
	if err != nil {
		return nil, apperrors.ErrCheck(err)
	}
	return review.ToReviewResponse(), nil
}

func (rs *ReviewService) GetReviewStatistics(productID, storeID uint) (*model.ProductReviewsStatistics, error) {
	exists, err := rs.reviewRepo.ProductExists(productID, storeID)
	if err != nil {
		return nil, apperrors.ErrCheck(err)
	}
	if !exists {
		return nil, apperrors.NewNotFoundError("product not found")
	}

	stats, err := rs.reviewRepo.GetReviewStatistics(productID, storeID)
	if err != nil {
		return nil, apperrors.ErrCheck(err)
	}
	return stats, nil
}
