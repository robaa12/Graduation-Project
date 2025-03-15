package validation

import (
	"errors"

	"github.com/robaa12/product-service/cmd/model"
	"gorm.io/gorm"
)

type CollectionValidator struct {
	DB *gorm.DB
}

func NewCollectionValidator(db *gorm.DB) *CollectionValidator {
	return &CollectionValidator{DB: db}
}

func (v *CollectionValidator) CollectionExists(storeID, collectionID uint) (*model.Collection, error) {
	var collection model.Collection
	if err := v.DB.Where("store_id = ? AND id = ?", storeID, collectionID).First(&collection).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &collection, nil
}

func (v *CollectionValidator) ValidateProductsExist(storeID uint, productIDs []uint) ([]model.Product, error) {
	var products []model.Product
	if err := v.DB.Where("store_id = ? AND id IN (?)", storeID, productIDs).Find(&products).Error; err != nil {
		return nil, err
	}
	if len(products) != len(productIDs) {
		return nil, errors.New("one or more products not found or don't belong to this store")
	}
	return products, nil
}

func (v *CollectionValidator) ValidateProductExists(storeID, productID uint) (*model.Product, error) {
	var product model.Product
	if err := v.DB.Where("store_id = ? AND id = ?", storeID, productID).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return &product, nil
}

func (v *CollectionValidator) ValidateCollectionRequest(req *model.CollectionRequest) error {
	if req.Name == "" {
		return errors.New("collection name is required")
	}
	if req.Description == "" {
		return errors.New("collection description is required")
	}
	return nil
}

func (v *CollectionValidator) GetCollectionWithProducts(collectionID, storeID uint) (*model.Collection, error) {
	var collection model.Collection
	if err := v.DB.Preload("Products").Where("store_id = ? AND id = ?", storeID, collectionID).First(&collection).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("collection not found")
		}
	}
	return &collection, nil
}
