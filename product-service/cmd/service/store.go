package service

import (
	"github.com/robaa12/product-service/cmd/database"
	apperrors "github.com/robaa12/product-service/cmd/errors"
	"github.com/robaa12/product-service/cmd/model"
	"github.com/robaa12/product-service/cmd/repository"
)

type StoreService struct {
	db   *database.Database
	repo *repository.StoreRepository
}

func NewStoreService(db *database.Database, repo *repository.StoreRepository) *StoreService {
	return &StoreService{
		db:   db,
		repo: repo,
	}
}

// / CreateStore creates a new store in the database
func (s *StoreService) CreateStore(storeRequest *model.StoreRequest) (*model.StoreResponse, error) {
	// Create a new store in the database
	store := storeRequest.ToStore()
	// Check if the store already exists
	existingStore, err := s.repo.GetStoreByID(store.ID)
	if err != nil {
		return nil, apperrors.NewInternalServerError("Failed to check if store exists")
	}
	if existingStore != nil {
		return nil, apperrors.NewBadRequestError("Store already exists")
	}
	// call CreateStore method from repository
	err = s.repo.CreateStore(store)
	if err != nil {
		return nil, apperrors.NewInternalServerError("Failed to create store")
	}

	storeResponse := store.ToStoreResponse()
	return storeResponse, nil
}

// deleteStore deletes a store from the database using the store ID
func (s *StoreService) DeleteStore(storeID uint) error {
	// Check if the store exists
	store, err := s.repo.GetStoreByID(storeID)
	if err != nil {
		return apperrors.NewInternalServerError("Failed to check if store exists")
	}
	if store == nil {
		return apperrors.NewNotFoundError("Store not found")
	}
	// Call the DeleteStore method from the repository
	err = s.repo.DeleteStore(storeID)
	if err != nil {
		return apperrors.NewInternalServerError("Failed to delete store")
	}
	return nil
}
