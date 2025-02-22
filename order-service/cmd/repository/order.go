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
	return r.db.Preload("Customer").First(order, id).Error
}

func (r *OrderRepository) GetOrderDetails(order *model.Order, id string) error {
	return r.db.Preload("OrderItems").Preload("Customer").First(order, id).Error
}
func (r *OrderRepository) GetAllOrder(id string) ([]model.Order, error) {
	var orders []model.Order
	err := r.db.Preload("Customer").Where("store_id = ?", id).Find(&orders).Error
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *OrderRepository) AddOrder(orderRequest *model.OrderRequestDetails) (*model.Order, error) {

	// start transaction
	tx := r.db.Begin()
	// Defer rollback if transaction fails
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// create Store
	store := orderRequest.CreateStore()
	if err := CreateStore(store, tx); err != nil {
		tx.Rollback()
		return nil, err
	}
	// create Customer
	customer := orderRequest.CreateCustomer()
	if err := AddCustomer(customer, tx); err != nil {
		tx.Rollback()
		return nil, err
	}
	// create store customer which's trace orders history for the customer at specific store
	storeCustomer := model.CreateStoreCustmer(store.ID, customer.ID)
	if _, err := AddStoreCustomer(storeCustomer, tx); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create order
	order := orderRequest.CreateOrder(customer.ID)
	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	// Create order items
	for _, item := range orderRequest.OrderItems {
		orderItem := item.CreateOrderItem(order.ID)
		if err := CreateOrderItem(orderItem, tx); err != nil {
			tx.Rollback()
			return nil, err
		}
		order.OrderItems = append(order.OrderItems, *orderItem)
	}

	// Commit transaction
	tx.Commit()
	order.Customer = *customer

	return order, nil
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
