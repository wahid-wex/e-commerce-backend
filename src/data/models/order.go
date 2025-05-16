package models

import (
	"gorm.io/gorm"
	"time"
)

type Cart struct {
	gorm.Model
	CustomerID uint `gorm:"not null;uniqueIndex"`

	Customer  Customer   `gorm:"foreignKey:CustomerID"`
	CartItems []CartItem `gorm:"foreignKey:CartID"`
}

type CartItem struct {
	gorm.Model
	CartID    uint    `gorm:"not null;index"`
	ProductID uint    `gorm:"not null;index"`
	Quantity  int     `gorm:"not null"`
	Price     float64 `gorm:"not null"`

	Cart    Cart    `gorm:"foreignKey:CartID"`
	Product Product `gorm:"foreignKey:ProductID"`
}

type Order struct {
	gorm.Model
	CustomerID      uint        `gorm:"not null;index"`
	SellerID        uint        `gorm:"not null;index"`
	TotalAmount     float64     `gorm:"not null"`
	Status          OrderStatus `gorm:"size:20;not null;default:'pending'"`
	ShippingAddress string      `gorm:"size:255;not null"`
	TrackingNumber  string      `gorm:"size:100"`
	OrderDate       time.Time   `gorm:"not null"`

	Customer   Customer    `gorm:"foreignKey:CustomerID"`
	Seller     Seller      `gorm:"foreignKey:SellerID"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	gorm.Model
	OrderID   uint    `gorm:"not null;index"`
	ProductID uint    `gorm:"not null;index"`
	Quantity  int     `gorm:"not null"`
	Price     float64 `gorm:"not null"`

	Order   Order   `gorm:"foreignKey:OrderID"`
	Product Product `gorm:"foreignKey:ProductID"`
}
