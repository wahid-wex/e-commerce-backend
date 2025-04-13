package models

import (
	"gorm.io/gorm"
	"time"
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"   // در انتظار تایید
	OrderStatusProcessed OrderStatus = "processed" // پردازش شده
	OrderStatusShipped   OrderStatus = "shipped"   // ارسال شده
	OrderStatusDelivered OrderStatus = "delivered" // دریافت شده
	OrderStatusReturned  OrderStatus = "returned"  // بازگشت خورده
	OrderStatusReviewing OrderStatus = "reviewing" // در حال بررسی
	OrderStatusRefunded  OrderStatus = "refunded"  // تایید بازگشت و عملیات بازپرداخت
)

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
)

// Customer represents a customer user
type Customer struct {
	gorm.Model
	ID              uint   `gorm:"uniqueIndex;not null"`
	FirstName       string `gorm:"size:100"`
	LastName        string `gorm:"size:100"`
	Address         string `gorm:"size:255"`
	PostalCode      string `gorm:"size:20"`
	Phone           string `gorm:"size:20"`
	CardNumber      string `gorm:"size:20"` // For refunds
	ShippingAddress string `gorm:"size:255"`
	Username        string `gorm:"size:100;not null;unique"`
	Email           string `gorm:"size:100;not null;unique"`
	Password        string `gorm:"size:100;not null"`
	// Relations

	ProductReactions []ProductReaction `gorm:"foreignKey:CustomerID"`
	ProductQuestions []ProductQuestion `gorm:"foreignKey:CustomerID"`
	// Removed Cart relation to avoid circular reference
	Orders  []Order  `gorm:"foreignKey:CustomerID"`
	Reviews []Review `gorm:"foreignKey:CustomerID"`

	UserRoles []*UserRole `gorm:"foreignKey:CustomerID"`
}

// Seller represents a seller user
type Seller struct {
	gorm.Model
	ID               uint      `gorm:"uniqueIndex;not null"`
	StoreName        string    `gorm:"size:100;not null"`
	BusinessLicense  string    `gorm:"size:100"`
	NationalID       string    `gorm:"size:20"`
	SatisfactionRate float64   `gorm:"default:0"`
	RegistrationDate time.Time `gorm:"not null"`
	IsVerified       bool      `gorm:"default:false"`
	Address          string    `gorm:"size:255"`
	Phone            string    `gorm:"size:20"`
	Description      string    `gorm:"size:1000"`
	Logo             string    `gorm:"size:255"`
	Username         string    `gorm:"size:100;not null;unique"`
	Email            string    `gorm:"size:100;not null;unique"`
	Password         string    `gorm:"size:100;not null"`
	// Relations

	Products       []Product       `gorm:"foreignKey:SellerID"`
	ProductStocks  []ProductStock  `gorm:"foreignKey:SellerID"`
	ProductAnswers []ProductAnswer `gorm:"foreignKey:SellerID"`
	Orders         []Order         `gorm:"foreignKey:SellerID"`
	UserRoles      []*UserRole     `gorm:"foreignKey:SellerID"`
}

// Admin represents an admin user
type Admin struct {
	gorm.Model
	ID         uint   `gorm:"uniqueIndex;not null"`
	FirstName  string `gorm:"size:100"`
	LastName   string `gorm:"size:100"`
	Department string `gorm:"size:100"`
	Username   string `gorm:"size:100;not null;unique"`
	Email      string `gorm:"size:100;not null;unique"`
	Password   string `gorm:"size:100;not null"`
}

// Role represents a user role
type Role struct {
	gorm.Model
	Name        string `gorm:"size:50;not null;unique"`
	Description string `gorm:"size:255"`

	// Relations - only one direction to avoid circular references
	UserRoles       []UserRole       `gorm:"foreignKey:RoleID"`
	RolePermissions []RolePermission `gorm:"foreignKey:RoleID"`
}

// Permission represents a permission
type Permission struct {
	gorm.Model
	Name        string `gorm:"size:100;not null;unique"`
	Description string `gorm:"size:255"`

	// Relations - only one direction to avoid circular references
	RolePermissions []RolePermission `gorm:"foreignKey:PermissionID"`
}

// RolePermission represents the many-to-many relationship between roles and permissions
type RolePermission struct {
	gorm.Model
	RoleID       uint `gorm:"not null;index:idx_role_permission,unique"`
	PermissionID uint `gorm:"not null;index:idx_role_permission,unique"`

	// Relations - only one direction to avoid circular references
	Role       Role       `gorm:"foreignKey:RoleID"`
	Permission Permission `gorm:"foreignKey:PermissionID"`
}

// UserRole represents the many-to-many relationship between users and roles
type UserRole struct {
	gorm.Model
	CustomerID uint `gorm:"index"` // Customer can have roles (nullable)
	SellerID   uint `gorm:"index"` // Seller can have roles (nullable)
	AdminID    uint `gorm:"index"` // Admin can have roles (nullable)
	RoleID     uint `gorm:"not null;index:idx_user_role,unique"`

	// Relations
	Customer *Customer `gorm:"foreignKey:CustomerID;references:ID"`
	Seller   *Seller   `gorm:"foreignKey:SellerID;references:ID"`
	Admin    *Admin    `gorm:"foreignKey:AdminID;references:ID"`
	Role     Role      `gorm:"foreignKey:RoleID"`
}
