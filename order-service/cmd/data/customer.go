package data

// CustomerRequest and their Method that mapping CustomerRequest into CustomserModel
type CustomerRequest struct {
	Email        uint   `json:"email"  binding:"required"`
	CustomerName string `json:"customer_name"  binding:"required"`
	PhoneNumber  string `json:"phone_number"  binding:"required"`
}

// CustomerResponse and their Method that mapping  CustomserModel into CustomerResponse
type CustomerResponse struct {
	Email        uint   `json:"email"`
	CustomerName string `json:"customer_name"`
	PhoneNumber  string `json:"phone_number" `
}

// CustomerDetailsResponse and their Method that mapping  CustomserModel into CustomerDetailsResponse
type CustomerDetailsResponse struct {
	CustomerInfo CustomerResponse `json:"customer_info"`
	Orders       []OrderResponse  `json:"orders"`
}

func (customerRequest *CustomerRequest) CreateCustomer() *Customer {

	return &Customer{
		Email:        customerRequest.Email,
		CustomerName: customerRequest.CustomerName,
		PhoneNumber:  customerRequest.PhoneNumber,
	}
}

func (customer *Customer) CreateCustomerResponse() *CustomerResponse {
	return &CustomerResponse{
		Email:        customer.Email,
		CustomerName: customer.CustomerName,
		PhoneNumber:  customer.PhoneNumber,
	}
}

func (customer *Customer) CreateCustomerDetailsResponse() *CustomerDetailsResponse {
	var ordersResponse []OrderResponse
	for _, order := range customer.Orders {
		ordersResponse = append(ordersResponse, *order.CreateOrderResponse())
	}
	return &CustomerDetailsResponse{
		CustomerInfo: *customer.CreateCustomerResponse(),
		Orders:       ordersResponse,
	}
}
