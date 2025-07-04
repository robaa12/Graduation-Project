package model

import "sort"

type SKURequest struct {
	Stock          int              `json:"stock" binding:"required"`
	Price          float64          `json:"price" binding:"required"`
	CompareAtPrice float64          `json:"compare_at_price" binding:"required"`
	CostPerItem    float64          `json:"cost_per_item" binding:"required"`
	Profit         float64          `json:"profit" binding:"required"`
	Margin         float64          `json:"margin" binding:"required"`
	ImageURL       string           `json:"image_url,omitempty"`
	Variants       []VariantRequest `json:"variants" binding:"required"`
}
type SKUsRequest struct {
	IDs []uint `json:"sku-ids" binding:"required"`
}
type SKUsResponse struct {
	SKUs []SKUProductResponse `json:"skus"`
}
type SKUProductResponse struct {
	SKUID       uint   `json:"sku_id" gorm:"column:sku_id"`
	SKUName     string `json:"sku_name" gorm:"column:sku_name"`
	ProductID   uint   `json:"product_id" gorm:"column:product_id"`
	ProductName string `json:"product_name" gorm:"column:product_name"`
	ImageURL    string `json:"image_url" gorm:"column:image_url"`
}

type SKUResponse struct {
	ID             uint              `json:"id"`
	Name           string            `json:"name"`
	Price          float64           `json:"price"`
	Stock          int               `json:"stock"`
	CostPerItem    float64           `json:"cost_per_item"`
	Profit         float64           `json:"profit"`
	Margin         float64           `json:"margin"`
	CompareAtPrice float64           `json:"compare_at_price"`
	ImageURL       string            `json:"image_url,omitempty"`
	Variants       []VariantResponse `json:"variants"`
}

func (s *SKURequest) ToSKU() *Sku {
	return &Sku{
		Stock:          s.Stock,
		Price:          s.Price,
		CompareAtPrice: s.CompareAtPrice,
		CostPerItem:    s.CostPerItem,
		Profit:         s.Profit,
		Margin:         s.Margin,
		ImageURL:       s.ImageURL,
	}
}
func (s *SKURequest) CreateSKU(productID uint) *Sku {
	sku := s.ToSKU()
	sku.ProductID = productID
	sku.Name = s.generateSKUName()
	return sku
}

// generateSKUName generates a name for the SKU based on its variants.value only like red,small,128 after sorting the variants based on Name .
func (s *SKURequest) generateSKUName() string {
	//generate map contain key the variant name and value the variant value
	variantMap := make(map[string]string)
	variantNames := []string{}
	for _, variant := range s.Variants {
		variantMap[variant.Name] = variant.Value
		variantNames = append(variantNames, variant.Name)
	}
	// Sort the variant names
	sort.Strings(variantNames)

	// Generate the SKU name by concatenating the sorted variant values
	var skuName string
	for _, name := range variantNames {
		if value, exists := variantMap[name]; exists {
			if skuName != "" {
				skuName += ","
			}
			skuName += value
		}
	}
	return skuName
}

// map sku object to sku response object

func (s *Sku) ToSKUResponse() *SKUResponse {
	variants := []VariantResponse{}
	for i, skuVariant := range s.SKUVariants {
		variants = append(variants, VariantResponse{
			Name:  s.Variants[i].Name,
			Value: skuVariant.Value,
		})
	}
	return &SKUResponse{
		ID:             s.ID,
		Name:           s.Name,
		Price:          s.Price,
		Stock:          s.Stock,
		CostPerItem:    s.CostPerItem,
		Profit:         s.Profit,
		Margin:         s.Margin,
		CompareAtPrice: s.CompareAtPrice,
		Variants:       variants,
		ImageURL:       s.ImageURL,
	}
}
