package model

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

var db *gorm.DB

func New(dbPool *gorm.DB) Models {
	db = dbPool

	return Models{
		Product:    Product{},
		SKU:        Sku{},
		Variant:    Variant{},
		Review:     Review{},
		Collection: Collection{},
		Store:      Store{},
	}
}

type Models struct {
	Product    Product
	SKU        Sku
	Variant    Variant
	Review     Review
	Collection Collection
	Store      Store
}

// BaseModel to reduce code repetition
type BaseModel struct {
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Store related with products, reviews,collections,categories 'one to many'
type Store struct {
	ID         uint         `json:"id" gorm:"primaryKey;autoIncrement:false"` // Disable auto-increment
	Name       string       `json:"name" gorm:"size:255;not null"`
	Product    []Product    `json:"products" gorm:"foreignKey:StoreID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`    // One-to-many relationship with products
	Review     []Review     `json:"reviews" gorm:"foreignKey:StoreID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`     // One-to-many relationship with reviews
	Collection []Collection `json:"collections" gorm:"foreignKey:StoreID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // One-to-many relationship with collections
	Category   []Category   `json:"categories" gorm:"foreignKey:StoreID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`  // One-to-many relationship with categories
	Slug       string       `json:"slug" gorm:"size:255;not null"`
	BaseModel
}
type Review struct {
	ID             uint   `json:"id" gorm:"primaryKey"`
	ProductID      uint   `json:"product_id" gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	StoreID        uint   `json:"store_id" gorm:"not null;index"`
	UserName       string `json:"user_name" gorm:"size:255;not null"`
	Rating         int    `json:"rating" gorm:"not null;check:rating >= 1 AND rating <= 5"`
	Title          string `json:"title" gorm:"size:255"`
	Description    string `json:"description" gorm:"type:text"`
	Published      bool   `json:"published" gorm:"default:true"`
	Classification bool   `json:"classification"`
	BaseModel
}

type Product struct {
	ID            uint           `json:"_" gorm:"primaryKey"`
	StoreID       uint           `json:"store_id" gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // Foreign key for Store
	CategoryID    *uint          `json:"category_id" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`      // Nullable FK for Category
	Name          string         `json:"name" gorm:"size:255;not null"`
	Description   string         `json:"description" gorm:"type:text"`
	Published     bool           `json:"published" gorm:"default:true"`
	StartPrice    float64        `json:"startPrice" gorm:"not null"`
	SKUs          []Sku          `json:"skus" gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // One-to-many relationship with SKU
	Slug          string         `json:"slug" gorm:"size:255;not null"`
	MainImageURL  string         `json:"main_image_url" gorm:"size:255;not null"`
	ImagesURL     pq.StringArray `json:"images_url" gorm:"type:text[];not null;default:'{}'"`
	Category      Category       `json:"category"`
	CollectionIDs []uint         `json:"-" gorm:"-"`
	BaseModel
}

type Sku struct {
	ID             uint         `json:"_" gorm:"primaryKey"`
	Name           string       `json:"name" gorm:"size:255;not null"`
	ProductID      uint         `json:"product_id" gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // Foreign key for Product
	Stock          int          `json:"stock" gorm:"not null"`
	Price          float64      `json:"price" gorm:"not null"`
	CompareAtPrice float64      `json:"compare_at_price" gorm:"not null"`
	CostPerItem    float64      `json:"cost_per_item" gorm:"not null"`
	Profit         float64      `json:"profit" gorm:"not null"`
	Margin         float64      `json:"margin" gorm:"not null"`
	ImageURL       string       `json:"image_url" gorm:"size:255"`
	Variants       []Variant    `json:"variants" gorm:"many2many:sku_variants;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // Many-to-many with Variants
	SKUVariants    []SKUVariant `json:"sku_variants" gorm:"foreignKey:SkuID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`   // One-to-many relationship with SKUVariant
	BaseModel
}

type Variant struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"size:255;not null;unique"`
	BaseModel
}

type Collection struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	StoreID     uint   `json:"store_id" gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // Foreign key for Store
	Name        string `json:"name" gorm:"size:255;not null"`
	Slug        string `json:"slug" gorm:"size:255;not null"`
	ImageURL    string `json:"image_url" gorm:"size:255"`
	Description string `json:"description" gorm:"type:text"`
	BaseModel
	Products []Product `json:"products" gorm:"many2many:collection_products;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
type Category struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	StoreID     uint   `json:"store_id" gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // Foreign key for Store
	Name        string `json:"name" gorm:"size:255;not null"`
	Slug        string `json:"slug" gorm:"size:255;not null"`
	Description string `json:"description" gorm:"type:text"`
	BaseModel
	Products []Product `json:"products" gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // One-to-many relationship with product
}

type SKUVariant struct {
	SkuID     uint   `gorm:"primaryKey;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // Foreign key for Sku
	VariantID uint   `gorm:"primaryKey;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // Foreign key for Variant
	Value     string `json:"value" gorm:"size:255;not null"`                          // Added correct `size` syntax
	BaseModel
}

func UpdateInventory(db *gorm.DB, skus []Sku) error {
	// Update inventory for each SKU
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if err := recover(); err != nil {
			tx.Rollback()
		}
	}()
	for _, sku := range skus {
		err := tx.Model(&Sku{}).
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
