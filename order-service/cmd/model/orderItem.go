package model

// OrderItemRequest  with their Function which map OrderItemRequest using orderId as arg Into OrderItemModel
type OrderItemRequest struct {
	SkuID    uint    `json:"sku_id" binding:"required"`
	Price    float64 `json:"price" binding:"required"`
	Quantity uint    `json:"quantity" binding:"required"`
}

// OrderItemResponse with their Function that mapping OrderItemModel into OrderItemResponse
type OrderItemResponse struct {
	ID         uint    `json:"id"`
	OrderID    uint    `json:"order_id"`
	SkuID      uint    `json:"sku_id"`
	SkuName    string  `json:"sku_name"`
	ProductID  uint    `json:"product_id,omitempty"`
	ProductName string `json:"product_name,omitempty"`
	ImageURL   string  `json:"image_url,omitempty"`
	Price      float64 `json:"price"`
	Quantity   uint    `json:"quantity"`
	Subtotal   float64 `json:"subtotal"`
}

func (orderItemRequest *OrderItemRequest) CreateOrderItem(orderID uint) *OrderItem {
	return &OrderItem{
		OrderID:  orderID,
		SkuID:    orderItemRequest.SkuID,
		Price:    orderItemRequest.Price,
		Quantity: orderItemRequest.Quantity,
	}
}

func (orderItem *OrderItem) CreateOrderItemResponse() *OrderItemResponse {
	return &OrderItemResponse{
		ID:       orderItem.ID,
		OrderID:  orderItem.OrderID,
		SkuID:    orderItem.SkuID,
		SkuName:  "", // Will be populated by the service layer
		Price:    orderItem.Price,
		Quantity: orderItem.Quantity,
		Subtotal: orderItem.Price * float64(orderItem.Quantity),
	}
}
