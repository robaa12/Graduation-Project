package model

import (
	"time"

	"gorm.io/gorm"
)

const dbTimeout = time.Second * 3

var db *gorm.DB

func New(dbPool *gorm.DB) Models {
	db = dbPool

	return Models{
		Product: Product{},
		SKU:     Sku{},
		Variant: Variant{},
	}
}

type Models struct {
	Product Product
	SKU     Sku
	Variant Variant
}

type Product struct {
	ID          uint           `json:"_" gorm:"primaryKey"`
	StoreID     uint           `json:"store_id" gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // Foreign key for Store
	Name        string         `json:"name" gorm:"size:255;not null"`
	Description string         `json:"description" gorm:"type:text"`
	Published   bool           `json:"published" gorm:"default:true"`
	StartPrice  float64        `json:"startprice" gorm:"not null"`
	Category    string         `json:"category" gorm:"size:255;not null"`
	SKUs        []Sku          `json:"skus" gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // One-to-many relationship with SKU
	Slug        string         `json:"slug" gorm:"size:255;not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type Sku struct {
	ID             uint           `json:"_" gorm:"primaryKey"`
	StoreID        uint           `json:"store_id" gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`   // Add store_id
	ProductID      uint           `json:"product_id" gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // Foreign key for Product
	Stock          int            `json:"stock" gorm:"not null"`
	Price          float64        `json:"price" gorm:"not null"`
	CompareAtPrice float64        `json:"compare_at_price" gorm:"not null"`
	CostPerItem    float64        `json:"cost_per_item" gorm:"not null"`
	Profit         float64        `json:"profit" gorm:"not null"`
	Margin         float64        `json:"margin" gorm:"not null"`
	Variants       []Variant      `json:"variants" gorm:"many2many:sku_variants;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // Many-to-many with Variants
	SKUVariants    []SKUVariant   `json:"sku_variants" gorm:"foreignKey:SkuID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`   // One-to-many relationship with SKUVariant
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

type Variant struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"size:255;not null;unique"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Collection struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	StoreID     uint           `json:"store_id" gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // Foreign key for Store
	Name        string         `json:"name" gorm:"size:255;not null"`
	Slug        string         `json:"slug" gorm:"size:255;not null"`
	ImageURL    string         `json:"image_url" gorm:"size:255"`
	Description string         `json:"description" gorm:"type:text"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Products    []Product      `json:"products" gorm:"many2many:collection_products;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type SKUVariant struct {
	SkuID     uint           `gorm:"primaryKey;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // Foreign key for Sku
	VariantID uint           `gorm:"primaryKey;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // Foreign key for Variant
	Value     string         `json:"value" gorm:"size:255;not null"`                          // Added correct `size` syntax
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (p *Product) GetProduct(id string) error {
	return db.First(p, id).Error
}

func (p *Product) UpdateProduct(id string) error {
	return db.Save(p).Error
}

func (s *Sku) GetSKU(id string) error {
	return db.First(s, id).Error
}

func (Product) TableName() string {
	return "products"
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	return tx.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_products_store_id_slug ON products (store_id, slug)").Error
}

func GetCollectionByID(db *gorm.DB, storeID uint, collectionID uint) (Collection, error) {
	var collection Collection
	err := db.Where("store_id = ? AND id = ?", storeID, collectionID).First(&collection).Error
	return collection, err
}

func GetAllCollections(db *gorm.DB) ([]Collection, error) {
	var collections []Collection
	err := db.Find(&collections).Error
	return collections, err
}

func AddProductToCollection(db *gorm.DB) ([]Collection, error) {
	var collections []Collection
	err := db.Find(&collections).Error
	return collections, err
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
