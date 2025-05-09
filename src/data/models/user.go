package models

import (
	"gorm.io/gorm"
	"time"
)

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

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
)

type Customer struct {
	gorm.Model
	ID               uint       `gorm:"uniqueIndex;not null"`
	FirstName        string     `gorm:"size:100"`
	LastName         string     `gorm:"size:100"`
	Address          string     `gorm:"size:255"`
	PostalCode       string     `gorm:"size:20"`
	Phone            string     `gorm:"size:20"`
	CardNumber       string     `gorm:"size:20"` // For refunds
	ShippingAddress  string     `gorm:"size:255"`
	Username         string     `gorm:"size:100;not null;unique"`
	Email            string     `gorm:"size:100;not null;unique"`
	Password         string     `gorm:"size:100;not null"`
	ProductReactions []Favorite `gorm:"foreignKey:CustomerID"`
	Orders           []Order    `gorm:"foreignKey:CustomerID"`
	Reviews          []Review   `gorm:"foreignKey:CustomerID"`

	UserRoles []*UserRole `gorm:"foreignKey:CustomerID"`
}

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
	ProductStocks []ProductStock `gorm:"foreignKey:SellerID;references:ID"`
	Orders        []Order        `gorm:"foreignKey:SellerID"`
	UserRoles     []*UserRole    `gorm:"foreignKey:SellerID"`
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
	CustomerID *uint `gorm:"index:idx_user_role_combo"`
	SellerID   *uint `gorm:"index:idx_user_role_combo"`
	AdminID    *uint `gorm:"index:idx_user_role_combo"`
	RoleID     uint  `gorm:"not null;index:idx_user_role_combo"`

	// Relations
	Customer *Customer `gorm:"foreignKey:CustomerID;references:ID"`
	Seller   *Seller   `gorm:"foreignKey:SellerID;references:ID"`
	Admin    *Admin    `gorm:"foreignKey:AdminID;references:ID"`
	Role     Role      `gorm:"foreignKey:RoleID"`
}
