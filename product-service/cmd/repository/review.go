package repository

import (
	"github.com/robaa12/product-service/cmd/database"
	apperrors "github.com/robaa12/product-service/cmd/errors"
	"github.com/robaa12/product-service/cmd/model"
)

type ReviewRepository struct {
	db database.Database
}

func NewReviewRepository(db database.Database) *ReviewRepository {
	return &ReviewRepository{db: db}
}

func (rr *ReviewRepository) CreateReview(review *model.Review) error {
	return rr.db.DB.Create(&review).Error
}

func (rr *ReviewRepository) GetProductReviews(productID, storeID uint) ([]model.Review, error) {
	var reviews []model.Review
	err := rr.db.DB.Where("product_id = ? AND store_id = ?", productID, storeID).Find(&reviews).Error
	err = apperrors.ErrCheck(err)
	return reviews, err
}

func (rr *ReviewRepository) GetReview(reviewID, productID, storeID uint) (*model.Review, error) {
	var review model.Review
	err := rr.db.DB.Where("id = ? AND product_id = ? AND store_id = ?", reviewID, productID, storeID).Find(review).Error
	err = apperrors.ErrCheck(err)
	return &review, err
}

func (rr *ReviewRepository) UpdateReview(review *model.Review) error {
	// Ensure the review exists
	_, err := rr.GetReview(review.ID, review.ProductID, review.StoreID)
	err = apperrors.ErrCheck(err)
	if err != nil {
		return err
	}
	return rr.db.DB.Save(&review).Error
}

func (rr *ReviewRepository) DeleteReview(reviewID, productID, storeID uint) error {
	// Ensure the review exists
	_, err := rr.GetReview(reviewID, productID, storeID)
	err = apperrors.ErrCheck(err)
	if err != nil {
		return err
	}
	err = rr.db.DB.Where("id = ? AND product_id = ? AND store_id = ?", reviewID, productID, storeID).Delete(&model.Review{}).Error
	err = apperrors.ErrCheck(err)
	return err
}

func (rr *ReviewRepository) GetReviewStatistics(productID, storeID uint) (*model.ProductReviewsStatistics, error) {
	var stats model.ProductReviewsStatistics

	// Count total reviews
	if err := rr.db.DB.Model(&model.Review{}).Where("product_id = ? AND store_id = ?", productID, storeID).Count(&stats.TotalReviews).Error; err != nil {
		err = apperrors.ErrCheck(err)
		return nil, err
	}

	if stats.TotalReviews == 0 {
		return &stats, nil
	}

	// Calculate average rating
	if err := rr.db.DB.Model(&model.Review{}).Where("product_id = ? AND store_id = ?", productID, storeID).Select("COALESCE(AVG(rating), 0) as average_rating").Scan(&stats.AverageRating).Error; err != nil {
		err = apperrors.ErrCheck(err)
		return nil, err
	}

	// Count ratings by value
	if err := rr.db.DB.Model(&model.Review{}).Where("product_id = ? AND store_id = ? AND rating = ?", productID, storeID, 5).Count(&stats.Rating5Count).Error; err != nil {
		err = apperrors.ErrCheck(err)
		return nil, err

	}

	if err := rr.db.DB.Model(&model.Review{}).Where("product_id = ? AND store_id = ? AND rating = ?", productID, storeID, 4).Count(&stats.Rating4Count).Error; err != nil {
		err = apperrors.ErrCheck(err)
		return nil, err
	}

	if err := rr.db.DB.Model(&model.Review{}).Where("product_id = ? AND store_id = ? AND rating = ?", productID, storeID, 3).Count(&stats.Rating3Count).Error; err != nil {
		err = apperrors.ErrCheck(err)
		return nil, err
	}

	if err := rr.db.DB.Model(&model.Review{}).Where("product_id = ? AND store_id = ? AND rating = ?", productID, storeID, 2).Count(&stats.Rating2Count).Error; err != nil {
		err = apperrors.ErrCheck(err)
		return nil, err
	}

	if err := rr.db.DB.Model(&model.Review{}).Where("product_id = ? AND store_id = ? AND rating = ?", productID, storeID, 1).Count(&stats.Rating1Count).Error; err != nil {
		err = apperrors.ErrCheck(err)
		return nil, err
	}

	return &stats, nil
}

func (rr *ReviewRepository) ProductExists(productID, storeID uint) (bool, error) {
	var count int64
	err := rr.db.DB.Model(&model.Product{}).Where("id = ? AND store_id = ?", productID, storeID).Count(&count).Error
	err = apperrors.ErrCheck(err)
	return count > 0, err

}
