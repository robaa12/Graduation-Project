package model

import (
	"time"
)

type OrderRequestDetails struct {
	StoreID uint `json:"store_id" binding:"required"`
	OrderRequest
	OrderItems []OrderItemRequest `json:"order_items" binding:"required"`
}
type OrderRequest struct {
	CustomerRequest
	CustomerName   string  `json:"customer_name"  binding:"required"`
	PhoneNumber    string  `json:"phone_number"  binding:"required"`
	Address        string  `json:"address" binding:"required"`
	TotalPrice     float64 `json:"total_price" binding:"required"`
	PaymentMethod  string  `json:"payment_method" binding:"required"`
	Note           string  `json:"note" binding:"required"`
	City           string  `json:"city" binding:"required"`
	Governorate    string  `json:"governorate" binding:"required"`
	PostalCode     string  `json:"postal_code" binding:"required"`
	ShippingMethod string  `json:"shipping_method" binding:"required"`
}

// orderResponse with their function that mapping OrderModel into OrderResponse

type OrderResponseInfo struct {
	ID             uint      `json:"order_id"`
	StoreID        uint      `json:"store_id"`
	StoreName      string    `json:"store_name"`
	CustomerName   string    `json:"customer_name"`
	PhoneNumber    string    `json:"phone_number" `
	Address        string    `json:"address"`
	TotalPrice     float64   `json:"total_price"`
	PaymentMethod  string    `json:"payment_method"`
	Note           string    `json:"note"`
	City           string    `json:"city"`
	Governorate    string    `json:"governorate" `
	PostalCode     string    `json:"postal_code" `
	ShippingMethod string    `json:"shipping_method"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
type OrderResponse struct {
	CustomerResponse
	OrderResponseInfo
}

// OrderDetailsResponse  with their function that mapping OrderModel using CustomerModel as arg into OrderDetailsResponse
type OrderDetailsResponse struct {
	OrderResponse
	OrderItems []OrderItemResponse `json:"order_items"`
}

func (orderRequest *OrderRequestDetails) CreateOrder(customerID uint) *Order {
	return &Order{
		StoreID:        orderRequest.StoreID,
		TotalPrice:     orderRequest.TotalPrice,
		CustomerID:     customerID,
		CustomerName:   orderRequest.CustomerName,
		PhoneNumber:    orderRequest.PhoneNumber,
		Address:        orderRequest.Address,
		PaymentMethod:  orderRequest.PaymentMethod,
		Note:           orderRequest.Note,
		City:           orderRequest.City,
		Governorate:    orderRequest.Governorate,
		PostalCode:     orderRequest.PostalCode,
		ShippingMethod: orderRequest.ShippingMethod,
	}
}

func (order *Order) CreateOrderResponseInfo() *OrderResponseInfo {
	return &OrderResponseInfo{
		ID:             order.ID,
		StoreID:        order.StoreID,
		StoreName:      "", // Will be populated by the service layer
		TotalPrice:     order.TotalPrice,
		CustomerName:   order.CustomerName,
		PhoneNumber:    order.PhoneNumber,
		Address:        order.Address,
		PaymentMethod:  order.PaymentMethod,
		Note:           order.Note,
		City:           order.City,
		Governorate:    order.Governorate,
		PostalCode:     order.PostalCode,
		ShippingMethod: order.ShippingMethod,
	}

}

func (order *Order) CreateOrderResponse() *OrderResponse {
	return &OrderResponse{
		CustomerResponse:  *order.Customer.CreateCustomerResponse(),
		OrderResponseInfo: *order.CreateOrderResponseInfo(),
	}

}

// CreateOrderDetailsResponse NOTE MAKE SURE IF customer is important to be on OrderDetailsResponse ------->
func (order *Order) CreateOrderDetailsResponse() *OrderDetailsResponse {

	var orderItems []OrderItemResponse
	for _, orderItem := range order.OrderItems {
		orderItems = append(orderItems, *orderItem.CreateOrderItemResponse())
	}
	return &OrderDetailsResponse{
		OrderResponse: *order.CreateOrderResponse(),
		OrderItems:    orderItems,
	}
}
