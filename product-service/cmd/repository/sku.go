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

func (sr *SkuRepository) GetSku(skuID, productID, storeID uint) (*model.Sku, error) {
	var sku model.Sku
	// Join Product and SKU tables and find the SKU with the given id, product_id and store_id in the database and preload the variants
	result := sr.db.DB.Model(&model.Sku{}).
		Joins("JOIN products ON skus.product_id = products.id").
		Where("skus.id = ? AND skus.product_id = ? AND products.store_id = ?", skuID, productID, storeID).
		Preload("Variants").Preload("SKUVariants").First(&sku)
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &sku, nil
}
func (sr *SkuRepository) GetSkus(storeID uint, skuIDs []uint) (*[]model.SKUProductResponse, error) {

	var skusResponse []model.SKUProductResponse
	result := sr.db.DB.Model(&model.Sku{}).
		Select("skus.id as sku_id, skus.name as sku_name, products.id as product_id, products.name as product_name, skus.image_url as image_url").
		Joins("JOIN products ON skus.product_id = products.id").
		Where("products.store_id = ? AND skus.id IN ?", storeID, skuIDs).
		Scan(&skusResponse)
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return &skusResponse, nil

}

func (sr *SkuRepository) UpdateSku(sku *model.Sku, storeID uint) error {
	result := sr.db.DB.Model(&model.Sku{}).
		Joins("JOIN products ON skus.product_id = products.id").
		Where("skus.id = ? AND skus.product_id = ? AND products.store_id = ?", sku.ID, sku.ProductID, storeID).
		Updates(&sku)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}
func (sr *SkuRepository) FindSku(skuID, productID, storeID uint) (*model.Sku, error) {
	var sku model.Sku
	result := sr.db.DB.Model(&model.Sku{}).
		Joins("JOIN products ON skus.product_id = products.id").
		Where("skus.id = ? AND skus.product_id = ? AND products.store_id = ?", skuID, productID, storeID).
		First(&sku)
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &sku, nil
}
func (sr *SkuRepository) DeleteSKU(sku *model.Sku, storeID uint) error {
	result := sr.db.DB.Model(&model.Sku{}).
		Joins("JOIN products ON skus.product_id = products.id").
		Where("skus.id = ? AND skus.product_id = ? AND products.store_id = ?", sku.ID, sku.ProductID, storeID).Unscoped().Delete(&sku)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}
func (sr *SkuRepository) CreateSku(storeID uint, sku *model.Sku, variantsRequest []model.VariantRequest) (*model.Sku, error) {
	// Check if the product exists in the store
	var product model.Product
	result := sr.db.DB.Where("id = ? AND store_id = ?", sku.ProductID, storeID).First(&product)
	if result.RowsAffected == 0 {
		return nil, errors.New("product not found or doesn't belong to the store")
	}
	if result.Error != nil {
		return nil, result.Error
	}

	// make transaction
	tx := sr.db.DB.Begin()
	// Add the SKU to the database
	if err := tx.Create(&sku).Error; err != nil {
		return nil, errors.New("error creating sku in database")
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
