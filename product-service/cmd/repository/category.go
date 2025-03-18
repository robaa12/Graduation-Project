package repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/robaa12/product-service/cmd/database"
	"github.com/robaa12/product-service/cmd/model"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *database.Database
}

func NewCategoryRepository(db *database.Database) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (cr *CategoryRepository) GetCategoryByID(storeID uint, categoryID uint) (*model.Category, error) {
	var category model.Category
	err := cr.db.DB.Where("store_id = ? AND id = ?", storeID, categoryID).Preload("Products").First(&category).Error
	return &category, err
}
func (cr *CategoryRepository) GetCategoryBySlug(storeID uint, slug string) (*model.Category, error) {
	var category model.Category
	err := cr.db.DB.Where("store_id = ? AND slug = ?", storeID, slug).Preload("Products").First(&category).Error
	return &category, err
}

func (cr *CategoryRepository) GetStoreCategories(storeID uint) ([]model.Category, error) {
	var category []model.Category
	err := cr.db.DB.Where("store_id = ?", storeID).Find(&category).Error
	return category, err
}

func (cr *CategoryRepository) UpdateCategory(storeID, categoryID uint, category *model.Category) error {
	return cr.db.DB.Model(&model.Product{}).Where("id = ? AND store_id = ?", categoryID, storeID).Updates(category).Error
}

func (cr *CategoryRepository) CreateCategory(category *model.Category) error {

	// Create category
	return cr.db.DB.Create(&category).Error
}

// GenerateCategorySlug checks if the slug is unique within the store and generates a new one if necessary.
func (cr *CategoryRepository) GenerateCategorySlug(name string, storeID uint) (string, error) {
	// Generate the base slug from the Category name
	baseSlug := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	slug := baseSlug
	var count int64

	// Loop to find a unique slug
	for i := 1; ; i++ {
		// Check if a Category with the same slug and store_id already exists
		cr.db.DB.Model(&model.Category{}).
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

func (cr *CategoryRepository) FindCategory(storeID, categoryID uint) (*model.Category, error) {
	var category model.Category
	result := cr.db.DB.Where("store_id = ? AND id = ?", storeID, categoryID).Preload("Products").First(&category)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &category, nil
}
func (cr *CategoryRepository) DeleteCategory(category *model.Category) error {
	return cr.db.DB.Unscoped().Delete(category).Error
}
