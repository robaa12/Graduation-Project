package model

// / StoreRequest is the request structure for store-related operations
type StoreRequest struct {
	ID   uint   `json:"id"`
	Name string `json:"name" gorm:"size:255;not null"`
	Slug string `json:"slug" gorm:"size:255;not null"`
}

// / StoreResponse is the response structure for store-related operations
type StoreResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// Store represents a store in the system
func (s *Store) ToStoreResponse() *StoreResponse {
	return &StoreResponse{
		ID:   s.ID,
		Name: s.Name,
		Slug: s.Slug,
	}
}

// / ToStore converts StoreRequest to Store model
func (sr *StoreRequest) ToStore() *Store {
	return &Store{
		ID:   sr.ID,
		Name: sr.Name,
		Slug: sr.Slug,
	}
}
