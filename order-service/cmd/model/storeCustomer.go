package model

// StoreCustomerItem represents a customer's information along with their order history for a specific store.
type StoreCustomerItem struct {
	CustomerResponseInfo         // Embedded struct containing customer details.
	NumberOfOrders       uint    `json:"number_of_orders"` // Total number of orders placed by the customer at the store.
	TotalSpent           float64 `json:"total_spent"`      // Total amount spent by the customer at the store.
}

// CreateStoreCustmer creates a StoreCustomer object from a StoreID, CustomerID
func CreateStoreCustmer(storeID, customerID uint) *StoreCustomer {
	return &StoreCustomer{
		StoreID:    storeID,
		CustomerID: customerID,
	}
}

// CreateStoreCustomerItem creates a StoreCustomerItem object from a StoreCustomer
func (storeCustomer *StoreCustomer) CreateStoreCustomerItem() *StoreCustomerItem {
	return &StoreCustomerItem{
		CustomerResponseInfo: *storeCustomer.Customer.CreateCustomerResponseInfo(), // Embeds customer details.
	}
}
