package data

import (
	"time"

	"gorm.io/gorm"
)

func New() Models {
	return Models{
		Order:     Order{},
		OrderItem: OrderItem{},
		Customer:  Customer{},
	}

}

type Models struct {
	Order     Order
	OrderItem OrderItem
	Customer  Customer
}
type Order struct {
	ID             uint           `json:"_" gorm:"primaryKey"`
	StoreID        uint           `json:"store_id" gorm:"not null; constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	TotalPrice     float64        `json:"total_price" gorm:"not null"`
	CustomerID     uint           `json:"customer_id" gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Address        string         `json:"address" gorm:"type:text;not null"`
	PaymentMethod  string         `json:"payment_method" gorm:"size:255;not null" `
	Note           string         `json:"note" gorm:"type:text"`
	City           string         `json:"city" gorm:"size:255;not null"`
	Governorate    string         `json:"governorate" gorm:"size:255"`
	PostalCode     string         `json:"postal_code" gorm:"size:255"`
	ShippingMethod string         `json:"shipping_method" gorm:"size:255;not null"`
	OrderItems     []OrderItem    `json:"order_items" gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

// you should add index constraint
type Customer struct {
	ID           uint    `json:"_" gorm:"primaryKey"`
	Email        uint    `json:"email" gorm:"size:255; not null"`
	CustomerName string  `json:"customer_name" gorm:"size:255; not null"`
	PhoneNumber  string  `json:"phone_number" gorm:"size:255; not null"`
	Orders       []Order `json:"orders" gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
type OrderItem struct {
	ID       uint    `json:"id" gorm:"primaryKey"`
	OrderID  uint    `json:"order_id" gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	SkuID    uint    `json:"sku_id" gorm:"not null; constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Price    float64 `json:"price" gorm:"not null"`
	Quantity uint    `json:"quantity" gorm:"not null"`
}

// Getters for each Model  using ID or Email

func (item *OrderItem) GetOrderItem(id string, db *gorm.DB) error {
	return db.First(item, id).Error
}
func (order *Order) GetOrder(id string, db *gorm.DB) error {
	return db.First(order, id).Error
}

// MAKE SURE FROM GetCustomer ------------------->
func (customer *Customer) GetCustomer(email string, db *gorm.DB) error {
	return db.Where("Email=?", email).First(&customer).Error
}

// update functions for each model

func (item *OrderItem) UpdateOrderItem(db *gorm.DB) error {
	return db.Save(item).Error
}

func (order *Order) UpdateOrder(db *gorm.DB) error {
	return db.Save(order).Error
}

func (customer *Customer) UpdateCustomer(db *gorm.DB) error {
	return db.Save(customer).Error
}
