package model

type CategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}
type CategoryInfo struct {
	ID   uint   `json:"id" binding:"required"`
	Slug string `json:"slug" binding:"required"`
}
type CategoryResponse struct {
	StoreID uint `json:"store_id"`
	CategoryInfo
	Name        string `json:"name"`
	Description string `json:"description"`
}
type CategoriesResponse struct {
	CategoriesResponse []CategoryResponse `json:"categories"`
}
type CategoryDetailsResponse struct {
	CategoryResponse
	Products []ProductResponse `json:"products"`
}

func (c *CategoryRequest) ToCategory(storeID uint) *Category {
	return &Category{
		StoreID:     storeID,
		Name:        c.Name,
		Description: c.Description,
	}
}

func (c *Category) ToCategoryInfo() *CategoryInfo {
	return &CategoryInfo{
		ID:   c.ID,
		Slug: c.Slug,
	}
}

func (c *Category) ToCategoryResponse() *CategoryResponse {
	return &CategoryResponse{
		CategoryInfo: *c.ToCategoryInfo(),
		StoreID:      c.StoreID,
		Name:         c.Name,
		Description:  c.Description,
	}
}

func (c *Category) ToCategoryDetailsResponse() *CategoryDetailsResponse {
	products := []ProductResponse{}
	for _, product := range c.Products {
		products = append(products, *product.ToProductResponse())
	}
	return &CategoryDetailsResponse{
		CategoryResponse: *c.ToCategoryResponse(),
		Products:         products,
	}
}
func GetCategoriesResponse(categories []Category) *CategoriesResponse {
	categoriesResponse := []CategoryResponse{}
	for _, category := range categories {
		categoriesResponse = append(categoriesResponse, *category.ToCategoryResponse())
	}
	return &CategoriesResponse{
		CategoriesResponse: categoriesResponse,
	}
}
