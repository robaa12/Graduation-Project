package model

import (
	"time"

	"github.com/robaa12/gatway-service/internal/middleware/auth"
)

type StoreInfo struct {
	StoreName     string `json:"store_name" validate:"required"`
	Description   string `json:"description" validate:"required"`
	BusinessPhone string `json:"business_phone" validate:"required"`
	CategoryID    uint   `json:"category_id" validate:"required"`
	PlanID        uint   `json:"plan_id" validate:"required"`
	StoreCurrency string `json:"store_currency" validate:"required"`
	Href          string `json:"href,omitempty"`
	Slug          string `json:"slug,omitempty"`
}
type StoreRequest struct {
	UserID uint `json:"user_id,omitempty"`
	StoreInfo
}
type Store struct {
	ID uint `json:"id"`
	StoreRequest
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// CreateStoreResponse is the response structure for store creation
type StoreResponse struct {
	Store       Store               `json:"store"`
	AccessToken *auth.TokenResponse `json:"access_token,omitempty"`
}

func (s *Store) GetStoreResponse(accessToken *auth.TokenResponse) *StoreResponse {
	return &StoreResponse{
		Store:       *s,
		AccessToken: accessToken,
	}
}
func (s *Store) ToServiceCreateStoreRequest() ServiceCreateStoreRequest {
	return ServiceCreateStoreRequest{
		ID:   s.ID,
		Name: s.StoreName,
	}
}
