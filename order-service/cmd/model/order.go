package model

import (
	"time"
)

// orderRequest and thier fuction which's mapping orderRequest to orderModel
type OrderRequestDetails struct {
	StoreID        uint               `json:"store_id" binding:"required"`
	TotalPrice     float64            `json:"total_price" binding:"required"`
	CustomerEmail  string             `json:"email" binding:"required"`
	CustomerName   string             `json:"customer_name"  binding:"required"`
	PhoneNumber    string             `json:"phone_number"  binding:"required"`
	Address        string             `json:"address" binding:"required"`
	PaymentMethod  string             `json:"payment_method" binding:"required"`
	Note           string             `json:"note" binding:"required"`
	City           string             `json:"city" binding:"required"`
	Governorate    string             `json:"governorate" binding:"required"`
	PostalCode     string             `json:"postal_code" binding:"required"`
	ShippingMethod string             `json:"shipping_method" binding:"required"`
	OrderItems     []OrderItemRequest `json:"order_items" binding:"required"`
}
type OrderRequest struct {
	TotalPrice     float64 `json:"total_price" binding:"required"`
	CustomerEmail  string  `json:"email" binding:"required"`
	CustomerName   string  `json:"customer_name"  binding:"required"`
	PhoneNumber    string  `json:"phone_number"  binding:"required"`
	Address        string  `json:"address" binding:"required"`
	PaymentMethod  string  `json:"payment_method" binding:"required"`
	Note           string  `json:"note" binding:"required"`
	City           string  `json:"city" binding:"required"`
	Governorate    string  `json:"governorate" binding:"required"`
	PostalCode     string  `json:"postal_code" binding:"required"`
	ShippingMethod string  `json:"shipping_method" binding:"required"`
}

// orderResponse with their fuction that mapping OrderModel into OrderResponse
type OrderResponse struct {
	ID             uint      `json:"order_id"`
	StoreID        uint      `json:"store_id"`
	TotalPrice     float64   `json:"total_price"`
	CustomerEmail  string    `json:"email" gorm:"size:255; not null"`
	CustomerName   string    `json:"customer_name"`
	PhoneNumber    string    `json:"phone_number" `
	Address        string    `json:"address"`
	PaymentMethod  string    `json:"payment_method"`
	Note           string    `json:"note"`
	City           string    `json:"city"`
	Governorate    string    `json:"governorate" `
	PostalCode     string    `json:"postal_code" `
	ShippingMethod string    `json:"shipping_method"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// orderDetailsResponse with their fuction that mapping OrderModel using CustomerModel as arg into OrderDetailsResponse
type OrderDetailsResponse struct {
	OrderInfo  OrderResponse       `json:"order_info"`
	OrderItems []OrderItemResponse `json:"order_items"`
}

func (orderRequest *OrderRequestDetails) CreateOrder() *Order {
	return &Order{
		StoreID:        orderRequest.StoreID,
		TotalPrice:     orderRequest.TotalPrice,
		CustomerEmail:  orderRequest.CustomerEmail,
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

func (order *Order) CreateOrderResponse() *OrderResponse {
	return &OrderResponse{
		ID:             order.ID,
		StoreID:        order.StoreID,
		TotalPrice:     order.TotalPrice,
		CustomerEmail:  order.CustomerEmail,
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

// NOTE MAKE SURE IF customerinto is important to be on OrderDetailsResponse ------->
func (order *Order) CreateOrderDetailsResponse() *OrderDetailsResponse {

	var orderItems []OrderItemResponse
	for _, orderItem := range order.OrderItems {
		orderItems = append(orderItems, *orderItem.CreateOrderItemResponse())
	}
	return &OrderDetailsResponse{
		OrderInfo:  *order.CreateOrderResponse(),
		OrderItems: orderItems,
	}
}
