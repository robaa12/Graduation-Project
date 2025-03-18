package model

import (
	"time"
)

type ReviewRequest struct {
	UserName    string `json:"user_name" binding:"required"`
	Rating      int    `json:"rating" binding:"required,min=1,max=5"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type ReviewResponse struct {
	ID             uint      `json:"id"`
	ProductID      uint      `json:"product_id"`
	StoreID        uint      `json:"store_id"`
	UserName       string    `json:"user_name"`
	Rating         int       `json:"rating"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Classification bool      `json:"classification"`
	CreatedAt      time.Time `json:"created_at"`
}

type ProductReviewsStatistics struct {
	TotalReviews  int64   `json:"total_reviews"`
	AverageRating float64 `json:"average_rating"`
	Rating5Count  int64   `json:"rating_5_count"`
	Rating4Count  int64   `json:"rating_4_count"`
	Rating3Count  int64   `json:"rating_3_count"`
	Rating2Count  int64   `json:"rating_2_count"`
	Rating1Count  int64   `json:"rating_1_count"`
}

func (r *Review) ToReviewResponse() *ReviewResponse {
	return &ReviewResponse{
		ID:          r.ID,
		ProductID:   r.ProductID,
		StoreID:     r.StoreID,
		UserName:    r.UserName,
		Rating:      r.Rating,
		Title:       r.Title,
		Description: r.Description,
		CreatedAt:   r.CreatedAt,
	}
}

func (rr *ReviewRequest) ToReview(productID, storeID uint) *Review {
	return &Review{
		ProductID:   productID,
		StoreID:     storeID,
		UserName:    rr.UserName,
		Rating:      rr.Rating,
		Title:       rr.Title,
		Description: rr.Description,
	}
}
