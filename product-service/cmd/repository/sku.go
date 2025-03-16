package repository

import (
	"errors"

	"github.com/robaa12/product-service/cmd/database"
	"github.com/robaa12/product-service/cmd/model"
	"gorm.io/gorm"
)

type SkuRepository struct {
	db *database.Database
}

func NewSkuRepository(db *database.Database) *SkuRepository {
	return &SkuRepository{db: db}
}

func (sr *SkuRepository) GetSku(skuID int) (*model.Sku, error) {
	var sku model.Sku
	result := sr.db.DB.Model(&model.Sku{}).Where("id = ?", skuID).Preload("Variants").Preload("SKUVariants").Find(&sku)
	if result.RowsAffected == 0 {
		return nil, errors.New("SKU Not Found")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &sku, nil

}

func (sr *SkuRepository) UpdateSku(sku *model.Sku) error {
	return sr.db.DB.Where("id= ?", sku.ID).Updates(&sku).Error

}
func (sr *SkuRepository) FindSku(skuID int) (*model.Sku, error) {
	var sku model.Sku
	result := sr.db.DB.Where("id = ?", skuID).Find(&sku)
	if result.RowsAffected == 0 {
		return nil, errors.New("Sku is not found.")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &sku, nil

}
func (sr *SkuRepository) DeleteSKU(sku *model.Sku) error {
	return sr.db.DB.Unscoped().Delete(&sku).Error
}
func (sr *SkuRepository) CreateSku(sku *model.Sku, variantsRequest []model.VariantRequest) (*model.Sku, error) {
	// make transaction
	tx := sr.db.DB.Begin()
	// Add the SKU to the database
	if err := tx.Create(&sku).Error; err != nil {
		return nil, errors.New("Error creating sku in database")
	}
	// Create SKU Variant
	for _, variant := range variantsRequest {

		variantData := variant.CreateVariant()
		if err := AddVariant(variantData, tx); err != nil {
			return nil, err
		}
		// Create a new SKU Variant
		skuVariant := model.CreateSkuVariant(sku.ID, variantData.ID, variant.Value)
		if err := AddSKUVariant(skuVariant, tx); err != nil {
			return nil, err
		}

	}
	// commit transaction
	err := tx.Commit().Error
	if err != nil {
		return nil, err
	}
	return sku, nil
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
