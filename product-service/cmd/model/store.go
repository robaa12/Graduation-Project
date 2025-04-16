package model

/// StoreRequest is the request structure for store-related operations
type StoreRequest struct {
	ID uint `json:"id"`
}

/// StoreResponse is the response structure for store-related operations
type StoreResponse struct {
	ID uint `json:"id"`
}

// Store represents a store in the system
func (s *Store) ToStoreResponse() *StoreResponse {
	return &StoreResponse{
		ID: s.ID,
	}
}

/// ToStore converts StoreRequest to Store model
func (sr *StoreRequest) ToStore() *Store {
	return &Store{
		ID: sr.ID,
	}
}
