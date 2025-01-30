package data

import (
	"log"
	"time"
)

// orderRequest and thier fuction which's mapping orderRequest to orderModel
type OrderRequest struct {
	StoreID        uint               `json:"store_id" binding:"required"`
	TotalPrice     float64            `json:"total_price" binding:"required"`
	Customers      CustomerRequest    `json:"customer" binding:"required"`
	Address        string             `json:"address" binding:"required"`
	PaymentMethod  string             `json:"payment_method" binding:"required"`
	Note           string             `json:"note" binding:"required"`
	City           string             `json:"city" binding:"required"`
	Governorate    string             `json:"governorate" binding:"required"`
	PostalCode     string             `json:"postal_code" binding:"required"`
	ShippingMethod string             `json:"shipping_method" binding:"required"`
	OrderItems     []OrderItemRequest `json:"order_items" binding:"required"`
}

// orderResponse with their fuction that mapping OrderModel into OrderResponse
type OrderResponse struct {
	ID             uint      `json:"order_id"`
	StoreID        uint      `json:"store_id"`
	TotalPrice     float64   `json:"total_price"`
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
	CustomerInfo CustomerResponse    `json:"customer_info"`
	OrderInfo    OrderResponse       `json:"order_info"`
	OrderItems   []OrderItemResponse `json:"order_items"`
}

func (orderRequest *OrderRequest) CreateOrder(customerID uint) *Order {
	return &Order{
		StoreID:        orderRequest.StoreID,
		TotalPrice:     orderRequest.TotalPrice,
		CustomerID:     customerID,
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
		Address:        order.Address,
		PaymentMethod:  order.PaymentMethod,
		Note:           order.Note,
		City:           order.City,
		Governorate:    order.Governorate,
		PostalCode:     order.PostalCode,
		ShippingMethod: order.ShippingMethod,
		CreatedAt:      order.CreatedAt,
		UpdatedAt:      order.UpdatedAt,
	}

}

// NOTE MAKE SURE IF customerinto is important to be on OrderDetailsResponse ------->
func (order *Order) CreateOrderDetailsResponse(customer *Customer) *OrderDetailsResponse {
	if customer.ID != order.CustomerID {
		// change this it very dirty
		log.Println("invaild Customer , order is done by anouther customer")
		return nil
	}

	var orderItems []OrderItemResponse
	for _, orderItem := range order.OrderItems {
		orderItems = append(orderItems, *orderItem.CreateOrderItemResponse())
	}
	return &OrderDetailsResponse{
		CustomerInfo: *customer.CreateCustomerResponse(),
		OrderInfo:    *order.CreateOrderResponse(),
		OrderItems:   orderItems,
	}
}
