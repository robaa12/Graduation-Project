package model

import "time"

type OrderStatusHistoryResonpse struct {
	From      string    `json:"from" `
	To        string    `json:"to" `
	ChangedAt time.Time `json:"changed_at" `
}

func (order *Order) CreateOrderStatusHistory(newStatus string) *OrderStatusHistory {
	return &OrderStatusHistory{
		OrderID: order.ID,
		From:    order.Status,
		To:      newStatus,
	}
}
func (orderStatusHistory *OrderStatusHistory) ToOrderStatusHistoryResponse() *OrderStatusHistoryResonpse {
	return &OrderStatusHistoryResonpse{
		From:      orderStatusHistory.From,
		To:        orderStatusHistory.To,
		ChangedAt: orderStatusHistory.ChangedAt,
	}
}

const (
	StatusPending    = "pending"
	StatusProcessing = "processing"
	StatusShipped    = "shipped"
	StatusDelivered  = "delivered"
	StatusCancelled  = "cancelled"
)

// validate if the inputl status transition is valid
func IsValidStatus(status string) bool {
	switch status {
	case StatusPending, StatusProcessing, StatusShipped, StatusDelivered, StatusCancelled:
		return true
	default:
		return false
	}
}

var allowedStatusTransitions = map[string][]string{
	StatusPending:    {StatusProcessing, StatusCancelled},
	StatusProcessing: {StatusShipped, StatusCancelled},
	StatusShipped:    {StatusDelivered},
	StatusDelivered:  {}, // Final state
	StatusCancelled:  {}, // Final state
}

func CanTransition(from, to string) bool {
	allowed, ok := allowedStatusTransitions[from]
	if !ok {
		return false
	}
	for _, status := range allowed {
		if status == to {
			return true
		}
	}
	return false
}
