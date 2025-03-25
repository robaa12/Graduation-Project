package repository

import (
	"order-service/cmd/model"

	"gorm.io/gorm"
)

type OrderItemRepository struct {
	db *gorm.DB
}

func NewOrderItemRepository(db *gorm.DB) *OrderItemRepository {
	return &OrderItemRepository{db: db}
}

// Getters for each Model  using ID or Email

func (r *OrderItemRepository) GetOrderItem(item *model.OrderItem, id string) error {
	return r.db.First(item, id).Error
}
func (r *OrderItemRepository) GetAllOrderItems(orderId string) ([]model.OrderItem, error) {

	var order model.Order
	err := r.db.Preload("OrderItems").First(&order, orderId).Error
	if err != nil {

		return nil, err
	}

	return order.OrderItems, nil
}

// CreateOrderItem Create functions for each model
func CreateOrderItem(orderItem *model.OrderItem, tx *gorm.DB) error {
	return tx.Create(orderItem).Error
}

func (r *OrderItemRepository) AddOrderItem(orderItem *model.OrderItem) error {
	return r.db.Create(orderItem).Error
}

// update functions for each model

func (r *OrderItemRepository) UpdateOrderItem(item *model.OrderItemRequest, id string) error {
	return r.db.Model(&model.OrderItem{}).Where("id = ?", id).Updates(item).Error
}

func (r *OrderItemRepository) FindOrderItem(orderItem *model.OrderItem, id string) (int64, error) {

	result := r.db.Find(orderItem, id)
	return result.RowsAffected, result.Error
}
func (r *OrderItemRepository) DeleteOrderItem(orderItem *model.OrderItem) error {
	return r.db.Unscoped().Delete(orderItem).Error
}
