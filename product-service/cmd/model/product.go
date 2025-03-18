package model

type ProductRequest struct {
	Name         string   `json:"name" binding:"required"`
	Description  string   `json:"description" binding:"required"`
	Published    bool     `json:"published" binding:"required"`
	StartPrice   float64  `json:"startPrice" binding:"required"`
	Slug         string   `json:"slug"`
	MainImageURL string   `json:"main_image_url" binding:"required"`
	ImagesURL    []string `json:"images_url"`

	SKUs     []SKURequest  `json:"skus" binding:"required"`
	Category *CategoryInfo `json:"category,omitempty" `
}
type ProductResponse struct {
	ID           uint          `json:"id"`
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Slug         string        `json:"slug"`
	Published    bool          `json:"published"`
	StartPrice   float64       `json:"startPrice"`
	MainImageURL string        `json:"main_image_url"`
	ImagesURL    []string      `json:"images_url,omitempty"`
	Category     *CategoryInfo `json:"category,omitempty"`
}
type ProductDetailsResponse struct {
	ProductResponse

	SKUs             []SKUResponse             `json:"skus"`
	ReviewStatistics *ProductReviewsStatistics `json:"review_statistics,omitempty"`
}

func (p *ProductRequest) CreateProduct(storeID uint) *Product {
	return &Product{
		Name:         p.Name,
		Description:  p.Description,
		StoreID:      storeID,
		Published:    p.Published,
		StartPrice:   p.StartPrice,
		Slug:         p.Slug,
		MainImageURL: p.MainImageURL,
		CategoryID:   &p.Category.ID,

		//ImagesURL:    p.ImagesURL,
	}
}

func (p *Product) ToProductResponse() *ProductResponse {
	return &ProductResponse{
		ID:           p.ID,
		Name:         p.Name,
		Slug:         p.Slug,
		Description:  p.Description,
		Published:    p.Published,
		StartPrice:   p.StartPrice,
		MainImageURL: p.MainImageURL,
		Category:     p.Category.ToCategoryInfo(),
		//ImagesURL:    p.ImagesURL,
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
