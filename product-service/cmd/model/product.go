package model

import (
	"github.com/lib/pq"
	"github.com/robaa12/product-service/cmd/utils"
)

type ProductRequest struct {
	Name         string        `json:"name" binding:"required"`
	Description  string        `json:"description" binding:"required"`
	Published    bool          `json:"published" binding:"required"`
	StartPrice   float64       `json:"startPrice" binding:"required"`
	Slug         string        `json:"slug"`
	MainImageURL string        `json:"main_image_url" binding:"required,url"`
	ImagesURL    []string      `json:"images_url"`
	SKUs         []SKURequest  `json:"skus" binding:"required"`
	Category     *CategoryInfo `json:"category,omitempty" `
}
type ProductResponse struct {
	ID           uint          `json:"id"`
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Slug         string        `json:"slug"`
	Published    bool          `json:"published"`
	StartPrice   float64       `json:"startPrice"`
	MainImageURL string        `json:"main_image_url"`
	ImagesURL    []string      `json:"images_url"`
	Category     *CategoryInfo `json:"category,omitempty"`
}
type ProductDetailsResponse struct {
	ProductResponse

	SKUs             []SKUResponse             `json:"skus"`
	ReviewStatistics *ProductReviewsStatistics `json:"review_statistics,omitempty"`
}

// PaginatedProductsResponse represents a paginated list of products
type PaginatedProductsResponse struct {
	Products    []ProductResponse `json:"products"`
	Total       int64             `json:"total"`
	Limit       int               `json:"limit"`
	Offset      int               `json:"offset"`
	IsPaginated bool              `json:"is_paginated"`
}

func (p *ProductRequest) CreateProduct(storeID uint) *Product {
	mainImageURL := utils.SanitizeURL(p.MainImageURL)
	imagesURL := make(pq.StringArray, 0)
	for _, url := range p.ImagesURL {
		if sanitizedURL := utils.SanitizeURL(url); sanitizedURL != "" {
			imagesURL = append(imagesURL, sanitizedURL)
		}
	}

	return &Product{
		Name:         p.Name,
		Description:  p.Description,
		StoreID:      storeID,
		Published:    p.Published,
		StartPrice:   p.StartPrice,
		Slug:         p.Slug,
		MainImageURL: mainImageURL,
		CategoryID:   &p.Category.ID,
		ImagesURL:    imagesURL,
	}
}

func (p *Product) ToProductResponse() *ProductResponse {
	images := make([]string, len(p.ImagesURL))
	for i, url := range p.ImagesURL {
		images[i] = url
	}

	return &ProductResponse{
		ID:           p.ID,
		Name:         p.Name,
		Slug:         p.Slug,
		Description:  p.Description,
		Published:    p.Published,
		StartPrice:   p.StartPrice,
		MainImageURL: p.MainImageURL,
		Category:     p.Category.ToCategoryInfo(),
		ImagesURL:    images,
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
