package repository

import (
	"github.com/robaa12/product-service/cmd/database"
	"github.com/robaa12/product-service/cmd/model"
)

type VariantRepository struct {
	db database.Database
}

func NewVariantRepository(db database.Database) *VariantRepository {
	return &VariantRepository{db: db}
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
