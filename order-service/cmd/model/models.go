package model

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel to reduce code repetition
type BaseModel struct {
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Order related with Store 'one to many', Customer 'one to many' , OrderItem 'one to many' , StoreCustomer 'one to many'
type Order struct {
	ID             uint                 `json:"id" gorm:"primaryKey"`
	StoreID        uint                 `json:"store_id" gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // Foreign key for Store
	TotalPrice     float64              `json:"total_price" gorm:"not null;index;"`
	CustomerID     uint                 `json:"customer_id" gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // Foreign key for Cutomer
	Customer       Customer             `json:"customer"`
	CustomerName   string               `json:"customer_name" gorm:"size:255; not null"`
	PhoneNumber    string               `json:"phone_number" gorm:"size:255; not null"`
	Address        string               `json:"address" gorm:"type:text;not null"`
	PaymentMethod  string               `json:"payment_method" gorm:"size:255;not null" `
	Note           string               `json:"note" gorm:"type:text"`
	City           string               `json:"city" gorm:"size:255;not null"`
	Governorate    string               `json:"governorate" gorm:"size:255"`
	PostalCode     string               `json:"postal_code" gorm:"size:255"`
	ShippingMethod string               `json:"shipping_method" gorm:"size:255;not null"`
	Status         string               `json:"status" gorm:"type:varchar(50);default:'pending';not null"`
	OrderItems     []OrderItem          `json:"order_items" gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	StatusHistory  []OrderStatusHistory `json:"status_history" gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // One-to-many relationship with OrderStatusHistory
	BaseModel
	Store Store `json:"store"`
}

// OrderItem related with Order 'one to many'
type OrderItem struct {
	ID       uint    `json:"id" gorm:"primaryKey"`
	OrderID  uint    `json:"order_id" gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	SkuID    uint    `json:"sku_id" gorm:"not null; constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Price    float64 `json:"price" gorm:"not null"`
	Quantity uint    `json:"quantity" gorm:"not null"`
}

// Customer related with Store 'many to many' via  StoreCustomer Table 'one to many', Order 'one to many'
type Customer struct {
	ID             uint            `json:"id" gorm:"primaryKey"`
	Email          string          `json:"email" gorm:"size:255;not null;uniqueIndex"`
	Orders         []Order         `json:"orders" gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`          // One-to-many relationship with Orders
	StoreCustomers []StoreCustomer `json:"store_customers" gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // One-to-many relationship with StoreCustomer
	BaseModel
}

// StoreCustomer use to Link many to many relationships between Customer and Store
type StoreCustomer struct {
	StoreID    uint     `gorm:"primaryKey;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // Foreign key for Store
	CustomerID uint     `gorm:"primaryKey;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // Foreign key for Customer
	Customer   Customer `json:"customer" gorm:"foreignKey:CustomerID"`
	BaseModel
}

// Store related with Customer 'many to many' via  StoreCustomer Table 'one to many', Order 'one to many'
type Store struct {
	ID             uint            `json:"id" gorm:"primaryKey;autoIncrement:false"`                                                // Disable auto-increment
	Name           string          `json:"name" gorm:"size:255"`                                                                    // Store name
	Slug           string          `json:"slug" gorm:"size:255"`                                                                    // Store slug
	Orders         []Order         `json:"orders" gorm:"foreignKey:StoreID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`           // One-to-many relationship with Orders
	Customers      []Customer      `json:"customers" gorm:"many2many:store_customers;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // Many-to-many with Customers
	StoreCustomers []StoreCustomer `json:"store_customers" gorm:"foreignKey:StoreID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`  // One-to-many relationship with StoreCustomer
	BaseModel
}
type OrderStatusHistory struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	OrderID   uint      `json:"order_id" gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	From      string    `json:"from" gorm:"size:50;not null"`
	To        string    `json:"to" gorm:"size:50;not null"`
	ChangedAt time.Time `json:"changed_at" gorm:"autoCreateTime"`
}
