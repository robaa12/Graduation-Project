package repository

import (
	"order-service/cmd/model"

	"gorm.io/gorm"
)

// StoreRepository handles database operations for stores
type StoreRepository struct {
	db *gorm.DB
}

// NewStoreRepository creates a new instance of StoreRepository
func NewStoreRepository(db *gorm.DB) *StoreRepository {
	return &StoreRepository{db: db}
}

// CreateStore inserts a new store into the database
func CreateStore(store *model.Store, tx *gorm.DB) error {
	return tx.FirstOrCreate(store, model.Store{ID: store.ID}).Error
}

// GetStoreByID retrieves a store by ID with its orders
func (r *StoreRepository) GetStoreWithOrders(store *model.Store) error {
	return r.db.Preload("Orders.Customers").First(&store, store.ID).Error

}

// Find Store By ID retrieves rows Affected
func (r *OrderItemRepository) FindStore(store *model.Store) (int64, error) {

	result := r.db.Find(store, store.ID)
	return result.RowsAffected, result.Error
}

// DeleteStore removes a store by ID
func (r *StoreRepository) DeleteStore(storeID uint) error {
	return r.db.Delete(&model.Store{}, storeID).Error
}
