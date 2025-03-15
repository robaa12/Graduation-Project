package model

type VariantRequest struct {
	Name  string `json:"name" binding:"required"`
	Value string `json:"value" binding:"required"`
}

type VariantResponse struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (v *VariantRequest) CreateVariant() *Variant {
	return &Variant{
		Name: v.Name,
	}
}

func CreateSkuVariant(skuID uint, variantID uint, value string) *SKUVariant {
	return &SKUVariant{
		SkuID:     skuID,
		VariantID: variantID,
		Value:     value,
	}
}
