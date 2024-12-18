package data

type ProductRequest struct {
	StoreID     uint         `json:"store_id" binding:"required"`
	Name        string       `json:"name" binding:"required"`
	Description string       `json:"description" binding:"required"`
	Published   bool         `json:"published" binding:"required"`
	StartPrice  float64      `json:"startprice" binding:"required"`
	Category    string       `json:"category" binding:"required"`
	SKUs        []SKURequest `json:"skus" binding:"required"`
}

type SKURequest struct {
	Stock          int              `json:"stock" binding:"required"`
	Price          float64          `json:"price" binding:"required"`
	CompareAtPrice float64          `json:"compare_at_price" binding:"required"`
	CostPerItem    float64          `json:"cost_per_item" binding:"required"`
	Profit         float64          `json:"profit" binding:"required"`
	Margin         float64          `json:"margin" binding:"required"`
	Variants       []VariantRequest `json:"variants" binding:"required"`
}

type VariantRequest struct {
	Name  string `json:"name" binding:"required"`
	Value string `json:"value" binding:"required"`
}

type ProductResponse struct {
	ID          uint    `json:"id"`
	StoreID     uint    `json:"store_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Published   bool    `json:"published"`
	StartPrice  float64 `json:"start_price"`
	Category    string  `json:"category"`
}
type ProductDetailsResponse struct {
	ID          uint          `json:"id"`
	StoreID     uint          `json:"store_id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Published   bool          `json:"published"`
	Category    string        `json:"category"`
	StartPrice  float64       `json:"start_price"`
	SKUs        []SKUResponse `json:"skus"`
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

type VariantResponse struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// map product object ro product response object
func (p *Product) ToProductDetailsResponse() ProductDetailsResponse {
	SKUs := []SKUResponse{}
	for _, sku := range p.SKUs {
		SKUs = append(SKUs, sku.ToSKUResponse())
	}
	return ProductDetailsResponse{
		ID:          p.ID,
		StoreID:     p.StoreID,
		Name:        p.Name,
		Description: p.Description,
		Published:   p.Published,
		StartPrice:  p.StartPrice,
		Category:    p.Category,
		SKUs:        SKUs,
	}
}

func (p *Product) ToProductResponse() ProductResponse {
	return ProductResponse{
		ID:          p.ID,
		StoreID:     p.StoreID,
		Name:        p.Name,
		Description: p.Description,
		Published:   p.Published,
		StartPrice:  p.StartPrice,
		Category:    p.Category,
	}
}

// map sku object to sku response object

func (s *Sku) ToSKUResponse() SKUResponse {
	var variants []VariantResponse
	for i, skuVariant := range s.SKUVariants {
		variants = append(variants, VariantResponse{
			Name:  s.Variants[i].Name,
			Value: skuVariant.Value,
		})
	}
	return SKUResponse{
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
