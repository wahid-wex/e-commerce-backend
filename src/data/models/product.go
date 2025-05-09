package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name        string `gorm:"size:100;not null"`
	Description string `gorm:"size:500"`
	ImageURL    string `gorm:"size:255"`

	Products []Product `gorm:"foreignKey:CategoryID"`
}

type Product struct {
	gorm.Model
	CategoryID       uint    `gorm:"not null"`
	Name             string  `gorm:"size:200;not null"`
	Description      string  `gorm:"size:1000"`
	Price            float64 `gorm:"not null"`
	ImageURL         string  `gorm:"size:255"`
	IsActive         bool    `gorm:"default:true"`
	SatisfactionRate float64 `gorm:"default:0"`

	Category          Category           `gorm:"foreignKey:CategoryID"`
	ProductAttributes []ProductAttribute `gorm:"foreignKey:ProductID"`
	ProductStocks     []ProductStock     `gorm:"foreignKey:ProductID"`
	Favorite          []Favorite         `gorm:"foreignKey:ProductID"`
	CartItems         []CartItem         `gorm:"foreignKey:ProductID"`
	OrderItems        []OrderItem        `gorm:"foreignKey:ProductID"`
	Reviews           []Review           `gorm:"foreignKey:ProductID"`
}

type ProductAttribute struct {
	gorm.Model
	ProductID uint   `gorm:"not null;index"`
	Key       string `gorm:"size:100;not null"`
	Value     string `gorm:"size:255;not null"`

	Product Product `gorm:"foreignKey:ProductID"`
}

type ProductStock struct {
	gorm.Model
	ProductID uint `gorm:"not null;index:idx_product_seller,unique"`
	SellerID  uint `gorm:"not null;index:idx_product_seller,unique"`
	Quantity  int  `gorm:"not null"`

	Product Product `gorm:"foreignKey:ProductID"`
	Seller  Seller  `gorm:"foreignKey:SellerID"`
}

type Review struct {
	gorm.Model
	CustomerID uint   `gorm:"not null;index"`
	ProductID  uint   `gorm:"not null;index"`
	Content    string `gorm:"size:1000"`
	Rating     int    `gorm:"not null;check:rating >= 1 AND rating <= 5"`

	Customer Customer `gorm:"foreignKey:CustomerID"`
	Product  Product  `gorm:"foreignKey:ProductID"`
}

type Favorite struct {
	gorm.Model
	ProductID  uint `gorm:"not null;index:idx_product_customer,unique"`
	CustomerID uint `gorm:"not null;index:idx_product_customer,unique"`

	Product  Product  `gorm:"foreignKey:ProductID"`
	Customer Customer `gorm:"foreignKey:CustomerID"`
}
