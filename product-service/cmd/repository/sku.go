package repository

import (
	"github.com/robaa12/product-service/cmd/database"
	"github.com/robaa12/product-service/cmd/model"
	"gorm.io/gorm"
)

type SkuRepository struct {
	db database.Database
}

func NewSkuRepository(db database.Database) *SkuRepository {
	return &SkuRepository{db: db}
}

func (sr *SkuRepository) GetSku(sku model.Sku, id int) error {
	return sr.db.DB.First(sku, id).Error
}

func (sr *SkuRepository) UpdateSku(sku model.Sku) error {
	return sr.db.DB.Save(sku).Error
}

func (sr *SkuRepository) CreateSku(sku model.Sku) error {
	return sr.db.DB.Create(sku).Error
}

// UpdateInventory updates the inventory of the SKUs
func (sr *SkuRepository) UpdateInventory(skus []model.Sku) error {
	// Update inventory for each SKU
	tx := sr.db.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()
	for _, sku := range skus {
		err := tx.Model(&model.Sku{}).
			Where("id = ?", sku.ID).
			Update("stock", gorm.Expr("stock - ?", sku.Stock)).
			Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}
