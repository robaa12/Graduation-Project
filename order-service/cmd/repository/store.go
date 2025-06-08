package repository

import (
	"order-service/cmd/model"

	"gorm.io/gorm"
)

type StoreRepository struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) *StoreRepository {
	return &StoreRepository{db: db}
}

// create store to the database
func (sr *StoreRepository) CreateStore(store *model.Store) error {
	result := sr.db.Create(store)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// delete store from the database
func (sr *StoreRepository) DeleteStore(storeID uint) error {
	result := sr.db.Delete(&model.Store{}, storeID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// get store by id from the database
func (sr *StoreRepository) GetStoreByID(storeID uint) (*model.Store, error) {
	result := &model.Store{}
	err := sr.db.First(result, storeID).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
func GetStoreByID(storeID uint, tx *gorm.DB) (*model.Store, error) {
	result := &model.Store{}
	err := tx.First(result, storeID).Error

	if err != nil {
		return nil, err
	}
	return result, nil
}
