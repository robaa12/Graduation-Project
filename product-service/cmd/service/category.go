package service

import (
	"errors"
	"fmt"
	"log"

	apperrors "github.com/robaa12/product-service/cmd/errors"
	"github.com/robaa12/product-service/cmd/model"
	"github.com/robaa12/product-service/cmd/repository"
)

type CategoryService struct {
	repository *repository.CategoryRepository
}

func NewCategoryService(repository *repository.CategoryRepository) *CategoryService {
	return &CategoryService{repository: repository}
}

func (cs *CategoryService) CreateCategory(storeID uint, categoryRequest *model.CategoryRequest) (*model.CategoryResponse, error) {
	if categoryRequest.Name == "" {
		return nil, errors.New("category name is required")
	}
	if categoryRequest.Description == "" {
		return nil, errors.New("category description is required")
	}
	category := categoryRequest.ToCategory(storeID)
	slug, err := cs.repository.GenerateCategorySlug(category.Name, category.StoreID)
	fmt.Println(slug)
	if err != nil {
		log.Println("Error Generating Category's Slug")
		return nil, err
	}
	category.Slug = slug
	err = cs.repository.CreateCategory(category)
	err = apperrors.ErrCheck(err)
	if err != nil {
		return nil, err
	}

	categoryResponse := category.ToCategoryResponse()
	return categoryResponse, nil
}

// GetCategories - GET /stores/{store_id}/categories/
func (cs *CategoryService) GetCategories(storeID uint) ([]model.CategoryResponse, error) {

	// Get categories from the database by store ID
	categories, err := cs.repository.GetStoreCategories(storeID)
	err = apperrors.ErrCheck(err)
	if err != nil {
		return nil, err
	}
	// Convert to response
	var categoriesResponse []model.CategoryResponse
	for _, category := range categories {
		categoriesResponse = append(categoriesResponse, *category.ToCategoryResponse())
	}
	return categoriesResponse, nil
}

// GetCategoryByID - GET /stores/{store_id}/categories/{category_id}
func (cs *CategoryService) GetCategoryByID(storeID, categoryID uint) (*model.CategoryDetailsResponse, error) {

	// Get category from the database by store ID and category ID
	category, err := cs.repository.GetCategoryByID(storeID, categoryID)
	err = apperrors.ErrCheck(err)
	if err != nil {
		return nil, err
	}
	categoryResponse := category.ToCategoryDetailsResponse()
	return categoryResponse, nil
}

// GetCategoryBySlug - GET /stores/{store_id}/categories/{category_slug}
func (cs *CategoryService) GetCategoryBySlug(storeID uint, CategorySlug string) (*model.CategoryDetailsResponse, error) {

	// Get category from the database by store ID and category slug
	category, err := cs.repository.GetCategoryBySlug(storeID, CategorySlug)
	err = apperrors.ErrCheck(err)
	if err != nil {
		return nil, err
	}
	categoryResponse := category.ToCategoryDetailsResponse()
	return categoryResponse, nil
}

// UpdateCategory Update Category - POST /stores/{store_id}/categories/{category_id}
func (cs *CategoryService) UpdateCategory(storeID, categoryID uint) error {
	// Validate category
	category, err := cs.repository.FindCategory(storeID, categoryID)
	err = apperrors.ErrCheck(err)
	if err != nil {
		return err
	}

	// update category
	err = cs.repository.UpdateCategory(storeID, categoryID, category)
	err = apperrors.ErrCheck(err)
	return err
}

// DeleteCategory Delete category  - DELETE /stores/{store_id}/categories/{category_id}
func (cs *CategoryService) DeleteCategory(storeID, categoryID uint) error {
	// Validate category
	category, err := cs.repository.FindCategory(storeID, categoryID)
	err = apperrors.ErrCheck(err)
	if err != nil {
		return err
	}

	// Delete category
	err = cs.repository.DeleteCategory(category)
	err = apperrors.ErrCheck(err)
	return err
}
