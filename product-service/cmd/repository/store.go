package repository

import (
	"github.com/robaa12/product-service/cmd/database"
	"github.com/robaa12/product-service/cmd/model"
	"gorm.io/gorm"
)

type StoreRepository struct {
	db database.Database
}

func NewStoreRepository(db database.Database) *StoreRepository {
	return &StoreRepository{db: db}
}

// create store to the database
func (sr *StoreRepository) CreateStore(store *model.Store) error {
	result := sr.db.DB.Create(store)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// get store by slug from the database
func (sr *StoreRepository) GetStoreBySlug(slug string) (*model.Store, error) {
	result := &model.Store{}
	err := sr.db.DB.Where("slug = ?", slug).First(result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

// delete store from the database
func (sr *StoreRepository) DeleteStore(storeID uint) error {
	result := sr.db.DB.Delete(&model.Store{}, storeID)
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
	err := sr.db.DB.First(result, storeID).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
