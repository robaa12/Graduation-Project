package model

// OrderItemReqest with their Fuction which's mapping OrderItemRequest using orderId as arg Into OrderItemModel
type OrderItemRequest struct {
	SkuID    uint    `json:"sku_id" binding:"required"`
	Price    float64 `json:"price" binding:"required"`
	Quantity uint    `json:"quantity" binding:"required"`
}

// OrderItemResponse with their Fuction that mapping OrderItemModel into OrderItemResponse
type OrderItemResponse struct {
	ID       uint    `json:"id"`
	OrderID  uint    `json:"order_id"`
	SkuID    uint    `json:"sku_id"`
	Price    float64 `json:"price"`
	Quantity uint    `json:"quantity"`
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
		Price:    orderItem.Price,
		Quantity: orderItem.Quantity,
	}
}
