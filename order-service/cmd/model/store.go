package model

// StoreResponse represents the basic response data for a store.
// / StoreRequest is the request structure for store-related operations
type StoreRequest struct {
	ID   uint   `json:"id"`
	Name string `json:"name" gorm:"size:255;not null"`
}

// / StoreResponse is the response structure for store-related operations
type StoreResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// StoreOrderResponse represents a store's response containing a list of orders.
type StoreOrderResponse struct {
	StoreResponse                 // Embedded struct containing the store's basic information.
	Orders        []OrderResponse `json:"orders" binding:"required"` // List of orders associated with the store.
}

// CreateStore creates a Store object from store id parm.

func CreateStore(storeId uint) *Store {
	return &Store{ID: storeId} // Returns a Store object with the ID set to the StoreID from OrderRequestDetails.
}

// CreateStoreResponse creates a StoreResponse object from a Store object.
func (store *Store) CreateStoreResponse() *StoreResponse {
	return &StoreResponse{
		ID: store.ID, // Returns a StoreResponse object with the StoreID set to the ID of the Store.
	}
}

// CreateStoreOrderResponse  creates a StoreOrderResponse object from a Store object.
func (store *Store) CreateStoreOrderResponse() *StoreOrderResponse {
	var orders []OrderResponse
	for _, order := range store.Orders {
		// Creates an OrderResponse for each order and appends it to the list.
		orders = append(orders, *order.CreateOrderResponse())
	}

	return &StoreOrderResponse{
		StoreResponse: *store.CreateStoreResponse(), // Embeds the store's basic information.
		Orders:        orders,                       // Sets the list of orders associated with the store.
	}
}

// Store represents a store in the system
func (s *Store) ToStoreResponse() *StoreResponse {
	return &StoreResponse{
		ID:   s.ID,
		Name: s.Name,
	}
}

// / ToStore converts StoreRequest to Store model
func (sr *StoreRequest) ToStore() *Store {
	return &Store{
		ID:   sr.ID,
		Name: sr.Name,
	}
}
