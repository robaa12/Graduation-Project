package messaging

import "time"

// Common message types
const (
	InventoryVerificationRequestType  = "inventory.verification.request"
	InventoryVerificationResponseType = "inventory.verification.response"
	InventoryUpdateRequestType        = "inventory.update.request"
	OrderCreatedEventType             = "order.created"
)

// Common exchanges
const (
	OrderExchange     = "orders"
	ProductExchange   = "products"
	InventoryExchange = "inventory"
)

type BaseMessage struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	CorrelID  string    `json:"correlation_id,omitempty"` // For request/response correlation
}

type InventoryItem struct {
	SkuID    uint    `json:"sku_id"`
	Quantity uint    `json:"quantitiy"`
	Price    float64 `json:"price"`
}

// InventoryVerificationRequest is sent to verify inventory before order creation
type InventoryVerificationRequest struct {
	BaseMessage
	StoreID uint            `json:"store_id"`
	Items   []InventoryItem `json:"items"`
}

// VerifiedItem contains verification results for a single item
type VerifiedItem struct {
	SkuID    uint     `json:"sku_id"`
	Valid    bool     `json:"valid"`
	InStock  bool     `json:"in_stock"`
	Price    float64  `json:"actual_price"`
	Messages []string `json:"messages,omitempty"`
}

// InventoryVerificationResponse contains verification results
type InventoryVerificationResponse struct {
	BaseMessage
	Valid   bool           `json:"valid"`
	Message string         `json:"message,omitempty"`
	Items   []VerifiedItem `json:"items"`
}

// InventoryUpdateRequest is sent to update inventory after order confirmation
type InventoryUpdateRequest struct {
	BaseMessage
	OrderID uint            `json:"order_id"`
	StoreID uint            `json:"store_id"`
	Items   []InventoryItem `json:"items"`
}

// OrderCreatedEvent is sent when a new order is created
type OrderCreatedEvent struct {
	BaseMessage
	OrderID  uint            `json:"order_id"`
	StoreID  uint            `json:"store_id"`
	Customer CustomerInfo    `json:"customer"`
	Items    []OrderItemInfo `json:"items"`
	Total    float64         `json:"total"`
}

type CustomerInfo struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// OrderItemInfo contains basic order item information for events
type OrderItemInfo struct {
	SkuID    uint    `json:"sku_id"`
	Quantity uint    `json:"quantity"`
	Price    float64 `json:"price"`
}
