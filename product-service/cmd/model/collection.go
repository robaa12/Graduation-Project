package model

type CollectionRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}
type CollectionResponse struct {
	ID          uint   `json:"id"`
	StoreID     uint   `json:"store_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
type CollectionDetailsResponse struct {
	CollectionResponse
	Products []ProductResponse `json:"products"`
}

func (c *Collection) ToCollectionResponse() *CollectionResponse {
	return &CollectionResponse{
		ID:          c.ID,
		StoreID:     c.StoreID,
		Name:        c.Name,
		Description: c.Description,
	}
}

func (c *Collection) ToCollectionDetailsResponse() *CollectionDetailsResponse {
	var products []ProductResponse
	for _, product := range c.Products {
		products = append(products, *product.ToProductResponse())
	}
	return &CollectionDetailsResponse{
		CollectionResponse: *c.ToCollectionResponse(),
		Products:           products,
	}
}
