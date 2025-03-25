package model

// CustomerRequest represents the request payload for creating or fetching a customer.
type CustomerRequest struct {
	CustomerEmail string `json:"email" binding:"required"` // Email of the customer. Marked as required.
}

// CustomerResponse represents the basic response data for a customer.
type CustomerResponse struct {
	CustomerEmail string `json:"customer_email" binding:"required"` // Email of the customer. Marked as required.
}

// CustomerResponseInfo represents detailed customer information including their ID and email.
type CustomerResponseInfo struct {
	CustomerID       uint `json:"custom_id" binding:"required"` // Unique identifier of the customer. Marked as required.
	CustomerResponse      // Embedded struct containing the customer's email.
}

// CustomerResponseDetails represents comprehensive customer information including their ID, email, and order history.
type CustomerResponseDetails struct {
	CustomerResponseInfo                     // Embedded struct containing the customer's ID and email.
	Orders               []OrderResponseInfo `json:"orders"` // List of orders associated with the customer.
}

// CreateCustomer creates a Customer object from a CustomerRequest object.
func (customer *CustomerRequest) CreateCustomer() *Customer {
	return &Customer{
		Email: customer.CustomerEmail, // Sets the customer's email from the request.
	}
}

// CreateCustomerResponse creates a CustomerResponse object from a Customer object.
func (customer *Customer) CreateCustomerResponse() *CustomerResponse {
	return &CustomerResponse{
		CustomerEmail: customer.Email, // Sets the customer's email in the response.
	}
}

// CreateCustomerResponseInfo creates a CustomerResponseInfo object from a Customer object.
func (customer *Customer) CreateCustomerResponseInfo() *CustomerResponseInfo {
	return &CustomerResponseInfo{
		CustomerID:       customer.ID,                        // Sets the customer's ID.
		CustomerResponse: *customer.CreateCustomerResponse(), // Embeds the customer's email.
	}
}

// CreateCustomerResponseDetails creates a CustomerResponseDetails object from a Customer object.
func (customer *Customer) CreateCustomerResponseDetails() *CustomerResponseDetails {
	var orders []OrderResponseInfo
	for _, order := range customer.Orders {
		// Creates an OrderResponseInfo for each order and appends it to the list.
		orders = append(orders, *order.CreateOrderResponseInfo())
	}
	return &CustomerResponseDetails{
		CustomerResponseInfo: *customer.CreateCustomerResponseInfo(), // Embeds the customer's ID and email.
		Orders:               orders,                                 // Sets the list of orders associated with the customer.
	}
}
