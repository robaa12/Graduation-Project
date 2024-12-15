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
		SKU:     SKU{},
		Variant: Variant{},
	}
}

type Models struct {
	Product Product
	SKU     SKU
	Variant Variant
}

type Product struct {
	gorm.Model
	StoreID         uint    `json:"store_id" gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Name            string  `json:"name" gorm:"size:255;not null"`
	Description     string  `json:"description" gorm:"type:text"`
	Published       bool    `json:"published" gorm:"default:true"`
	StartingAtPrice float64 `json:"starting_at_price" gorm:"not null"`
	CategoryID      uint    `json:"category_id" gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	SKUs            []SKU   `json:"skus" gorm:"foreignKey:ProductID"` // One-to-many relationship with SKU
}

type SKU struct {
	gorm.Model
	ProductID      uint      `json:"product_id" gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // Foreign key for Product
	Stock          int       `json:"stock" gorm:"not null"`
	Price          float64   `json:"price" gorm:"not null"`
	CompareAtPrice float64   `json:"compare_at_price" gorm:"not null"`
	CostPerItem    float64   `json:"cost_per_item" gorm:"not null"`
	Profit         float64   `json:"profit" gorm:"not null"`
	Margin         float64   `json:"margin" gorm:"not null"`
	Variants       []Variant `json:"variants" gorm:"many2many:sku_variants;joinForeignKey:sku_id;joinReferences:variant_id"` // Many-to-many with Variants
}

type Variant struct {
	gorm.Model
	Name string `json:"name" gorm:"size:255;not null;unique"`
	SKUs []SKU  `json:"skus" gorm:"many2many:sku_variants;joinForeignKey:variant_id;joinReferences:sku_id"` // Many-to-many with SKUs
}

type SKUVariant struct {
	SKUID     uint   `gorm:"column:sku_id;primaryKey"`
	VariantID uint   `gorm:"column:variant_id;primaryKey"`
	Value     string `json:"value" gorm:"type:text"` // Add the value column here
}

type Category struct {
	gorm.Model
	Name string `json:"name" gorm:"size:255;not null"`
}

func (p *Product) CreateProduct(productR ProductRequest) {
	p.Name = productR.Name
	p.Description = productR.Description
	p.StoreID = productR.StoreID
}

func (s *SKU) CreateSKU(skuR SKURequest, productID uint) {
	s.ProductID = productID
	s.Stock = skuR.Stock
	s.Price = skuR.Price
}

func (v *Variant) CreateVariant(variantR VariantRequest) {
	v.Name = variantR.Name
}

func (sv *SKUVariant) CreateSkuVariant(skuID uint, variantID uint, value string) {
	sv.SKUID = skuID
	sv.VariantID = variantID
	sv.Value = value
}

func (p *Product) GetProduct(id string) error {
	return db.First(p, id).Error
}

func (p *Product) UpdateProduct(id string) error {
	return db.Save(p).Error
}
