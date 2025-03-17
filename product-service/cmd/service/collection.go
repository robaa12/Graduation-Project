package service

import (
	"errors"
	"log"

	apperrors "github.com/robaa12/product-service/cmd/errors"
	"github.com/robaa12/product-service/cmd/model"
	"github.com/robaa12/product-service/cmd/repository"
)

type CollectionService struct {
	repository *repository.CollectionRepository
}

func NewCollectionService(repository *repository.CollectionRepository) *CollectionService {
	return &CollectionService{repository: repository}
}

func (cs *CollectionService) CreateCollection(storeID uint, collectionRequest *model.CollectionRequest) (*model.CollectionResponse, error) {
	if collectionRequest.Name == "" {
		return nil, errors.New("collection name is required")
	}
	if collectionRequest.Description == "" {
		return nil, errors.New("collection description is required")
	}
	collection := collectionRequest.ToCollection(storeID)
	slug, err := cs.repository.GenerateCollectionSlug(collection.Name, collection.StoreID)
	if err != nil {
		log.Println("Error Generating Collection's Slug")
		return nil, err
	}
	collection.Slug = slug
	err = cs.repository.CreateCollection(collection)
	err = apperrors.ErrCheck(err)
	if err != nil {
		return nil, err
	}

	collectionResponse := collection.ToCollectionResponse()
	return collectionResponse, nil
}

// GetCollections - GET /stores/{store_id}/collections/
func (cs *CollectionService) GetCollections(storeID uint) ([]model.CollectionResponse, error) {

	// Get collections from the database by store ID
	collections, err := cs.repository.GetStoreCollections(storeID)
	err = apperrors.ErrCheck(err)
	if err != nil {
		return nil, err
	}
	// Convert to response
	var collectionsResponse []model.CollectionResponse
	for _, collection := range collections {
		collectionsResponse = append(collectionsResponse, *collection.ToCollectionResponse())
	}
	return collectionsResponse, nil
}

// GetCollection - GET /stores/{store_id}/collections/{collection_id}
func (cs *CollectionService) GetCollection(storeID, collectionID uint) (*model.CollectionDetailsResponse, error) {

	// Get collection from the database by store ID and collection ID
	collection, err := cs.repository.GetCollectionByID(storeID, collectionID)
	err = apperrors.ErrCheck(err)
	if err != nil {
		return nil, err
	}
	collectionResponse := collection.ToCollectionDetailsResponse()
	return collectionResponse, nil
}

// AddProductToCollection Add product to collection - POST /stores/{store_id}/collections/{collection_id}
func (cs *CollectionService) AddProductToCollection(storeID, collectionID uint, collectionProductsRequest *model.CollectionProductsRequest) error {

	// Validate collection
	collection, err := cs.repository.FindCollection(storeID, collectionID)
	err = apperrors.ErrCheck(err)
	if err != nil {
		return err
	}

	// Validate products
	products, err := cs.repository.FindProducts(storeID, collectionProductsRequest)
	err = apperrors.ErrCheck(err)
	if err != nil {
		return err
	}
	// Add products to collection
	err = cs.repository.AddProductsToCollection(collection, products)
	err = apperrors.ErrCheck(err)
	return err
}

// RemoveProductFromCollection Remove product from collection - DELETE /stores/{store_id}/collections/{collection_id}/products/{product_id}
func (cs *CollectionService) RemoveProductFromCollection(storeID, collectionID, productID uint) error {
	// Validate collection Product
	collection, err := cs.repository.FindCollectionProduct(storeID, collectionID, productID)
	err = apperrors.ErrCheck(err)
	if err != nil {
		return err
	}

	// Remove product from collection
	err = cs.repository.RemoveProductFromCollection(collection.ID, collection.Products[0].ID)
	err = apperrors.ErrCheck(err)
	return err
}
