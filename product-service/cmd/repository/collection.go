package repository

import (
	"errors"

	"github.com/robaa12/product-service/cmd/database"
	"github.com/robaa12/product-service/cmd/model"
	"github.com/robaa12/product-service/cmd/utils"
)

type CollectionRepository struct {
	db database.Database
}

func NewCollectionRepository(db database.Database) *CollectionRepository {
	return &CollectionRepository{db: db}
}

func (cr *CollectionRepository) GetCollectionByID(storeID uint, collectionID uint) (model.Collection, error) {
	var collection model.Collection
	err := cr.db.DB.Where("store_id = ? AND id = ?", storeID, collectionID).First(&collection).Error
	return collection, err
}

func (cr *CollectionRepository) GetStoreCollections(storeID uint) ([]model.Collection, error) {
	var collections []model.Collection
	err := cr.db.DB.Where("store_id = ?", storeID).Find(&collections).Error
	return collections, err
}

func (cr *CollectionRepository) AddProductToCollection(productID uint, collectionID uint) error {
	return cr.db.DB.Exec("INSERT INTO collection_products (product_id, collection_id) VALUES (?, ?)", productID, collectionID).Error
}

func (cr *CollectionRepository) RemoveProductFromCollection(productID uint, collectionID uint) error {
	return cr.db.DB.Exec("DELETE FROM collection_products WHERE product_id = ? AND collection_id = ?", productID, collectionID).Error
}

func (cr *CollectionRepository) UpdateCollection(collection model.Collection) error {
	return cr.db.DB.Save(collection).Error
}

func (cr *CollectionRepository) CreateCollection(collection model.Collection) error {
	// Generate Slug
	baseSlug, err := utils.ValidateAndGenerateCollectionSlug(cr.db.DB, collection.Name, collection.StoreID)
	if err != nil {
		return errors.New("Error Generating Collection Slug")
	}
	collection.Slug = baseSlug

	// Create Collection
	return cr.db.DB.Create(collection).Error
}

func (cr *CollectionRepository) DeleteCollection(collectionID uint) error {
	if err := cr.db.DB.Exec("DELETE FROM collection_products WHERE collection_id = ?", collectionID).Error; err != nil {
		return err
	}
	// Then delete from the collections table
	return cr.db.DB.Exec("DELETE FROM collections WHERE id = ?", collectionID).Error
}
