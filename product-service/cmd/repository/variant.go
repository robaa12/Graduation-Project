package repository

import (
	"log"

	"github.com/robaa12/product-service/cmd/database"
	"github.com/robaa12/product-service/cmd/model"
	"gorm.io/gorm"
)

type VariantRepository struct {
	db database.Database
}

func NewVariantRepository(db database.Database) *VariantRepository {
	return &VariantRepository{db: db}
}
func AddVariant(v *model.Variant, tx *gorm.DB) error {
	if err := tx.FirstOrCreate(&v, model.Variant{Name: v.Name}).Error; err != nil {
		log.Println("Error creating variant in database")
		return err
	}
	return nil
}
func AddSKUVariant(sv *model.SKUVariant, tx *gorm.DB) error {
	if err := tx.Create(&sv).Error; err != nil {
		log.Println("Error creating sku variant in database")
		return err
	}
	return nil
}

func (vr *VariantRepository) GetVariant(v model.Variant, id int) error {
	return vr.db.DB.First(v, id).Error
}

func (vr *VariantRepository) UpdateVariant(v model.Variant) error {
	return vr.db.DB.Save(v).Error
}

func (vr *VariantRepository) CreateVariant(v model.Variant) error {
	return vr.db.DB.Create(v).Error
}

func (vr *VariantRepository) DeleteVariant(v model.Variant) error {
	return vr.db.DB.Delete(v).Error
}
