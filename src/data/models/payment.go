package models

import (
	"gorm.io/gorm"
	"time"
)

// Payment represents a payment for an order
type Payment struct {
	gorm.Model
	OrderID       uint          `gorm:"not null;uniqueIndex"`
	Amount        float64       `gorm:"not null"`
	Status        PaymentStatus `gorm:"size:20;not null;default:'pending'"`
	TransactionID string        `gorm:"size:100"`
	PaymentMethod string        `gorm:"size:50;not null"`
	PaymentDate   time.Time

	// Relations - only one direction to avoid circular references
	Order Order `gorm:"foreignKey:OrderID"`
}
