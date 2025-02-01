package repository

import (
	"order-service/cmd/model"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) GetOrder(order *model.Order, id string) error {
	return r.db.First(order, id).Error
}

func (r *OrderRepository) GetOrderDetails(order *model.Order, id string) error {
	return r.db.Preload("OrderItems").First(order, id).Error
}
func (r *OrderRepository) GetAllOrder(id string) ([]model.Order, error) {
	var orders []model.Order
	err := r.db.Where("store_id = ?", id).Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) AddOrder(order *model.Order) error {
	return r.db.Create(order).Error
}
func (r *OrderRepository) UpdateOrder(orderRequest *model.OrderRequest, id uint) error {
	return r.db.Model(&model.Order{}).Where("id = ?", id).Updates(orderRequest).Error
}
func (r *OrderRepository) FindOrder(order *model.Order, id string) (int64, error) {

	result := r.db.Preload("OrderItems").Find(order, id)
	return result.RowsAffected, result.Error
}
func (r *OrderRepository) DeleteOrder(order *model.Order) error {
	return r.db.Unscoped().Delete(order).Error
}
