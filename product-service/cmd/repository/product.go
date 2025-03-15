package repository

import (
	"fmt"
	"log"
	"strings"

	"github.com/robaa12/product-service/cmd/database"
	"github.com/robaa12/product-service/cmd/model"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db database.Database
}

func NewProductRepository(db database.Database) *ProductRepository {
	return &ProductRepository{db: db}
}

func (pr *ProductRepository) GetProduct(p model.Product, id int) error {
	return pr.db.DB.First(p, id).Error
}

func (pr *ProductRepository) GetStoreProducts(p *[]model.Product) error {
	return pr.db.DB.Find(p).Error
}

func (pr *ProductRepository) UpdateProduct(p model.Product) error {
	return pr.db.DB.Save(p).Error
}

func (pr *ProductRepository) CreateProduct(productRequest model.ProductRequest) (*model.Product, error) {
	// Generate product slug

	product := productRequest.CreateProduct()

	err := pr.db.DB.Transaction(func(tx *gorm.DB) error {

		// Add Product to the database
		if err := tx.Create(&product).Error; err != nil {
			log.Println("Error creating product in database")
			return err
		}
		for _, skuRequest := range productRequest.SKUs {
			// Create a new SKU
			sku := skuRequest.CreateSKU(product.ID, product.StoreID)

			// Add the SKU to the database
			if err := tx.Create(&sku).Error; err != nil {
				log.Println("Error creating sku in database")
				return err
			}

			for _, variantRequest := range skuRequest.Variants {
				// Create a new variant
				variant := variantRequest.CreateVariant()

				// Check if the variant already exists in the database or not and create it if it doesn't
				if err := tx.FirstOrCreate(&variant, model.Variant{Name: variantRequest.Name}).Error; err != nil {
					log.Println("Error creating variant in database")
					return err
				}

				// Create a new SKU Variant
				skuVariant := model.CreateSkuVariant(sku.ID, variant.ID, variantRequest.Value)

				// Add the SKU Variant to the database
				if err := tx.Create(&skuVariant).Error; err != nil {
					log.Println("Error creating sku variant in database")
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return product, nil
}

// Generate product slug
func (pr *ProductRepository) GenerateProductSlug(name string, storeID uint) (string, error) {
	// Generate the base slug from the product name
	baseSlug := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	slug := baseSlug
	var count int64

	// Loop to find a unique slug
	for i := 1; ; i++ {
		// Check if a product with the same slug and store_id already exists
		pr.db.DB.Model(&model.Product{}).
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
