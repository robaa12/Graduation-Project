package model

type SKURequest struct {
	Stock          int              `json:"stock" binding:"required"`
	Price          float64          `json:"price" binding:"required"`
	CompareAtPrice float64          `json:"compare_at_price" binding:"required"`
	CostPerItem    float64          `json:"cost_per_item" binding:"required"`
	Profit         float64          `json:"profit" binding:"required"`
	Margin         float64          `json:"margin" binding:"required"`
	Variants       []VariantRequest `json:"variants" binding:"required"`
}

type SKUResponse struct {
	ID             uint              `json:"id"`
	Price          float64           `json:"price"`
	Stock          int               `json:"stock"`
	CostPerItem    float64           `json:"cost_per_item"`
	Profit         float64           `json:"profit"`
	Margin         float64           `json:"margin"`
	CompareAtPrice float64           `json:"compare_at_price"`
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
	}
}
func (s *SKURequest) CreateSKU(productID uint, storeID uint) *Sku {
	sku := s.ToSKU()
	sku.ProductID = productID
	sku.StoreID = storeID
	return sku
}

// map sku object to sku response object

func (s *Sku) ToSKUResponse() *SKUResponse {
	var variants []VariantResponse
	for i, skuVariant := range s.SKUVariants {
		variants = append(variants, VariantResponse{
			Name:  s.Variants[i].Name,
			Value: skuVariant.Value,
		})
	}
	return &SKUResponse{
		ID:             s.ID,
		Price:          s.Price,
		Stock:          s.Stock,
		CostPerItem:    s.CostPerItem,
		Profit:         s.Profit,
		Margin:         s.Margin,
		CompareAtPrice: s.CompareAtPrice,
		Variants:       variants,
	}
}
