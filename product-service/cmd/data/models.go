package data

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
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type Sku struct {
	ID             uint           `json:"_" gorm:"primaryKey"`
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

type SKUVariant struct {
	SkuID     uint           `gorm:"primaryKey;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // Foreign key for Sku
	VariantID uint           `gorm:"primaryKey;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // Foreign key for Variant
	Value     string         `json:"value" gorm:"size:255;not null"`                          // Added correct `size` syntax
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (p *Product) CreateProduct(productR ProductRequest) {
	p.Name = productR.Name
	p.Description = productR.Description
	p.StoreID = productR.StoreID
	p.Published = productR.Published
	p.StartPrice = productR.StartPrice
	p.Category = productR.Category
}

func (s *Sku) CreateSKU(skuR SKURequest, productID uint) {
	s.ProductID = productID
	s.Stock = skuR.Stock
	s.Price = skuR.Price
	s.CompareAtPrice = skuR.CompareAtPrice
	s.CostPerItem = skuR.CostPerItem
	s.Profit = skuR.Profit
	s.Margin = skuR.Margin
}

func (v *Variant) CreateVariant(variantR VariantRequest) {
	v.Name = variantR.Name
}

func (sv *SKUVariant) CreateSkuVariant(skuID uint, variantID uint, value string) {
	sv.SkuID = skuID
	sv.VariantID = variantID
	sv.Value = value
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
