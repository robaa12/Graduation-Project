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

func AddStore(store *model.Store, tx *gorm.DB) error {
	result := tx.Create(store)
	if result.Error != nil {
		return result.Error
	}
	return nil
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
