package model

type CollectionRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	ImageURL    string `json:"image_url"`
}

type CollectionProductsRequest struct {
	ProductIDs []uint `json:"product_ids"`
}

type CollectionResponse struct {
	ID          uint   `json:"id"`
	StoreID     uint   `json:"store_id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
}
type CollectionDetailsResponse struct {
	CollectionResponse
	Products []ProductResponse `json:"products"`
}
type CollectionsResponse struct {
	Collections []CollectionResponse `json:"collections"`
}

func (c *CollectionRequest) ToCollection(storeID uint) *Collection {
	return &Collection{
		StoreID:     storeID,
		Name:        c.Name,
		Description: c.Description,
		ImageURL:    c.ImageURL,
	}
}

func (c *Collection) ToCollectionResponse() *CollectionResponse {
	return &CollectionResponse{
		ID:          c.ID,
		StoreID:     c.StoreID,
		Name:        c.Name,
		Description: c.Description,
		Slug:        c.Slug,
		ImageURL:    c.ImageURL,
	}
}

func (c *Collection) ToCollectionDetailsResponse() *CollectionDetailsResponse {
	products := []ProductResponse{}
	for _, product := range c.Products {
		products = append(products, *product.ToProductResponse())
	}
	return &CollectionDetailsResponse{
		CollectionResponse: *c.ToCollectionResponse(),
		Products:           products,
	}
}
func GetCollectionsResponse(collections []Collection) *CollectionsResponse {
	collectionsResponse := []CollectionResponse{}
	for _, collection := range collections {
		collectionsResponse = append(collectionsResponse, *collection.ToCollectionResponse())
	}
	return &CollectionsResponse{
		Collections: collectionsResponse,
	}
}
