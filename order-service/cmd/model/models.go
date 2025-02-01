package model

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID             uint           `json:"_" gorm:"primaryKey"`
	StoreID        uint           `json:"store_id" gorm:"not null; constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	TotalPrice     float64        `json:"total_price" gorm:"not null"`
	CustomerEmail  string         `json:"email" gorm:"size:255; not null"`
	CustomerName   string         `json:"customer_name" gorm:"size:255; not null"`
	PhoneNumber    string         `json:"phone_number" gorm:"size:255; not null"`
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

type OrderItem struct {
	ID       uint    `json:"id" gorm:"primaryKey"`
	OrderID  uint    `json:"order_id" gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	SkuID    uint    `json:"sku_id" gorm:"not null; constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Price    float64 `json:"price" gorm:"not null"`
	Quantity uint    `json:"quantity" gorm:"not null"`
}
