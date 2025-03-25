package repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/robaa12/product-service/cmd/database"
	"github.com/robaa12/product-service/cmd/model"
	"gorm.io/gorm"
)

type CollectionRepository struct {
	db *database.Database
}

func NewCollectionRepository(db *database.Database) *CollectionRepository {
	return &CollectionRepository{db: db}
}

func (cr *CollectionRepository) GetCollectionByID(storeID uint, collectionID uint) (*model.Collection, error) {
	var collection model.Collection
	err := cr.db.DB.Where("store_id = ? AND id = ?", storeID, collectionID).Preload("Products").First(&collection).Error
	return &collection, err
}

func (cr *CollectionRepository) GetStoreCollections(storeID uint) ([]model.Collection, error) {
	var collections []model.Collection
	err := cr.db.DB.Where("store_id = ?", storeID).Find(&collections).Error
	return collections, err
}

func (cr *CollectionRepository) AddProductsToCollection(collection *model.Collection, products []model.Product) error {
	return cr.db.DB.Model(collection).Association("Products").Append(&products)
}

func (cr *CollectionRepository) RemoveProductFromCollection(collectionID uint, productID uint) error {
	return cr.db.DB.Exec("DELETE FROM collection_products WHERE product_id = ? AND collection_id = ?", productID, collectionID).Error
}

func (cr *CollectionRepository) UpdateCollection(collection model.Collection) error {
	return cr.db.DB.Save(collection).Error
}

func (cr *CollectionRepository) CreateCollection(collection *model.Collection) error {

	// Create Collection
	return cr.db.DB.Create(&collection).Error
}

// GenerateCollectionSlug checks if the slug is unique within the store and generates a new one if necessary.
func (cr *CollectionRepository) GenerateCollectionSlug(name string, storeID uint) (string, error) {
	// Generate the base slug from the collection name
	baseSlug := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	slug := baseSlug
	var count int64

	// Loop to find a unique slug
	for i := 1; ; i++ {
		// Check if a collection with the same slug and store_id already exists
		cr.db.DB.Model(&model.Collection{}).
			Where("slug = ? AND store_id = ?", slug, storeID).
			Count(&count)

		// If no duplicate is found, return the unique slug
		if count == 0 {
			return slug, nil
		}

		// If a duplicate exists, append a counter to the slug
		slug = fmt.Sprintf("%s-%d", baseSlug, i)
	}
}

func (cr *CollectionRepository) DeleteCollection(collectionID uint) error {
	if err := cr.db.DB.Exec("DELETE FROM collection_products WHERE collection_id = ?", collectionID).Error; err != nil {
		return err
	}
	// Then delete from the collections table
	return cr.db.DB.Exec("DELETE FROM collections WHERE id = ?", collectionID).Error
}
func (cr *CollectionRepository) FindProducts(storeID uint, collectionProductsRequests *model.CollectionProductsRequest) ([]model.Product, error) {
	var products []model.Product
	if err := cr.db.DB.Where("store_id = ? AND id IN (?)", storeID, collectionProductsRequests.ProductIDs).Find(&products).Error; err != nil {
		return nil, err
	}
	if len(products) != len(collectionProductsRequests.ProductIDs) {
		return nil, errors.New("one or more products not found or don't belong to this store")
	}
	return products, nil
}

func (cr *CollectionRepository) FindCollection(storeID, collectionID uint) (*model.Collection, error) {
	var collection model.Collection
	if err := cr.db.DB.Where("store_id = ? AND id = ?", storeID, collectionID).First(&collection).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &collection, nil
}

func (cr *CollectionRepository) FindCollectionProduct(storeID, collectionID, productID uint) (*model.Collection, error) {
	var collection model.Collection

	// One query: find the collection and preload the specific product
	err := cr.db.DB.
		Model(&model.Collection{}).
		Where("id = ? AND store_id = ?", collectionID, storeID).
		Preload("Products", func(db *gorm.DB) *gorm.DB {
			return db.Where("products.id = ?", productID)
		}).
		First(&collection).Error

	if err != nil {
		return nil, gorm.ErrRecordNotFound
	}

	if len(collection.Products) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &collection, nil
}
