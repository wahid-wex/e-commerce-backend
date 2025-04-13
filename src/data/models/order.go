package models

import (
	"gorm.io/gorm"
	"time"
)

// Cart represents a user's shopping cart
type Cart struct {
	gorm.Model
	CustomerID uint `gorm:"not null;uniqueIndex"`

	// Relations - only one direction to avoid circular references
	Customer  Customer   `gorm:"foreignKey:CustomerID"`
	CartItems []CartItem `gorm:"foreignKey:CartID"`
}

// CartItem represents an item in a user's shopping cart
type CartItem struct {
	gorm.Model
	CartID    uint    `gorm:"not null;index"`
	ProductID uint    `gorm:"not null;index"`
	Quantity  int     `gorm:"not null"`
	Price     float64 `gorm:"not null"` // Price at the time of adding to cart

	// Relations - only one direction to avoid circular references
	Cart    Cart    `gorm:"foreignKey:CartID"`
	Product Product `gorm:"foreignKey:ProductID"`
}

// Order represents a user's order
type Order struct {
	gorm.Model
	CustomerID      uint        `gorm:"not null;index"`
	SellerID        uint        `gorm:"not null;index"`
	TotalAmount     float64     `gorm:"not null"`
	Status          OrderStatus `gorm:"size:20;not null;default:'pending'"`
	ShippingAddress string      `gorm:"size:255;not null"`
	TrackingNumber  string      `gorm:"size:100"`
	OrderDate       time.Time   `gorm:"not null"`

	// Relations - only one direction to avoid circular references
	Customer   Customer    `gorm:"foreignKey:CustomerID"`
	Seller     Seller      `gorm:"foreignKey:SellerID"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID"`
	// Removed Payment relation to avoid circular reference
}

// OrderItem represents an item in an order
type OrderItem struct {
	gorm.Model
	OrderID   uint    `gorm:"not null;index"`
	ProductID uint    `gorm:"not null;index"`
	Quantity  int     `gorm:"not null"`
	Price     float64 `gorm:"not null"` // Price at the time of order

	// Relations - only one direction to avoid circular references
	Order   Order   `gorm:"foreignKey:OrderID"`
	Product Product `gorm:"foreignKey:ProductID"`
}
