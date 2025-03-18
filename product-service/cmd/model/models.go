package model

import (
	"time"

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
	}
}

type Models struct {
	Product    Product
	SKU        Sku
	Variant    Variant
	Review     Review
	Collection Collection
}

// BaseModel to reduce code repetition
type BaseModel struct {
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
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
	ID           uint    `json:"_" gorm:"primaryKey"`
	StoreID      uint    `json:"store_id" gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // Foreign key for Store
	CategoryID   *uint   `json:"category_id" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`      // Nullable FK for Category
	Name         string  `json:"name" gorm:"size:255;not null"`
	Description  string  `json:"description" gorm:"type:text"`
	Published    bool    `json:"published" gorm:"default:true"`
	StartPrice   float64 `json:"startPrice" gorm:"not null"`
	SKUs         []Sku   `json:"skus" gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // One-to-many relationship with SKU
	Slug         string  `json:"slug" gorm:"size:255;not null"`
	MainImageURL string  `json:"main_image_url" gorm:"size:255;not null"`
	//ImagesURL    []string `json:"images_url" gorm:"size:255"`
	Category Category `json:"category"`
	BaseModel
}

type Sku struct {
	ID             uint         `json:"_" gorm:"primaryKey"`
	ProductID      uint         `json:"product_id" gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // Foreign key for Product
	Stock          int          `json:"stock" gorm:"not null"`
	Price          float64      `json:"price" gorm:"not null"`
	CompareAtPrice float64      `json:"compare_at_price" gorm:"not null"`
	CostPerItem    float64      `json:"cost_per_item" gorm:"not null"`
	Profit         float64      `json:"profit" gorm:"not null"`
	Margin         float64      `json:"margin" gorm:"not null"`
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
	Slug        string `json:"slug" gorm:"size:255;not null;uniqueIndex"`
	ImageURL    string `json:"image_url" gorm:"size:255"`
	Description string `json:"description" gorm:"type:text"`
	BaseModel
	Products []Product `json:"products" gorm:"many2many:collection_products;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
type Category struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	StoreID     uint   `json:"store_id" gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // Foreign key for Store
	Name        string `json:"name" gorm:"size:255;not null"`
	Slug        string `json:"slug" gorm:"size:255;not null;uniqueIndex"`
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

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	return tx.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_products_store_id_slug ON products (store_id, slug)").Error
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
