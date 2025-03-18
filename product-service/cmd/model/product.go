package model

type ProductRequest struct {
	StoreID     uint         `json:"store_id" binding:"required"`
	Name        string       `json:"name" binding:"required"`
	Description string       `json:"description" binding:"required"`
	Published   bool         `json:"published" binding:"required"`
	StartPrice  float64      `json:"startPrice" binding:"required"`
	Slug        string       `json:"slug"`
	Category    string       `json:"category" binding:"required"`
	SKUs        []SKURequest `json:"skus" binding:"required"`
}
type ProductResponse struct {
	ID          uint    `json:"id"`
	StoreID     uint    `json:"store_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Slug        string  `json:"slug"`
	Published   bool    `json:"published"`
	StartPrice  float64 `json:"startPrice"`
	Category    string  `json:"category"`
}
type ProductDetailsResponse struct {
	ProductResponse
	SKUs             []SKUResponse             `json:"skus"`
	ReviewStatistics *ProductReviewsStatistics `json:"review_statistics,omitempty"`
}

func (p *ProductRequest) CreateProduct() *Product {
	return &Product{
		Name:        p.Name,
		Description: p.Description,
		StoreID:     p.StoreID,
		Published:   p.Published,
		StartPrice:  p.StartPrice,
		Category:    p.Category,
		Slug:        p.Slug,
	}
}

func (p *Product) ToProductResponse() *ProductResponse {
	return &ProductResponse{
		ID:          p.ID,
		StoreID:     p.StoreID,
		Name:        p.Name,
		Slug:        p.Slug,
		Description: p.Description,
		Published:   p.Published,
		StartPrice:  p.StartPrice,
		Category:    p.Category,
	}
}

// ToProductDetailsResponse map product object to product response object
func (p *Product) ToProductDetailsResponse() *ProductDetailsResponse {
	var SKUs []SKUResponse
	for _, sku := range p.SKUs {
		SKUs = append(SKUs, *sku.ToSKUResponse())
	}

	return &ProductDetailsResponse{
		ProductResponse: *p.ToProductResponse(),
		SKUs:            SKUs,
	}
}
func (p *ProductDetailsResponse) WithReviewStatistics(stats *ProductReviewsStatistics) *ProductDetailsResponse {
	p.ReviewStatistics = stats
	return p
}
