package data

type ProductRequest struct {
	StoreID     uint         `json:"store_id" binding:"required"`
	Name        string       `json:"name" binding:"required"`
	Description string       `json:"description" binding:"required"`
	SKUs        []SKURequest `json:"skus" binding:"required"`
}

type SKURequest struct {
	Stock    int              `json:"stock" binding:"required"`
	Price    float64          `json:"price" binding:"required"`
	Variants []VariantRequest `json:"variants" binding:"required"`
}

type VariantRequest struct {
	Name  string `json:"name" binding:"required"`
	Value string `json:"value" binding:"required"`
}
