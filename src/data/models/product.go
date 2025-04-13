package models

import "gorm.io/gorm"

// Category represents a product category
type Category struct {
	gorm.Model
	Name        string `gorm:"size:100;not null"`
	Description string `gorm:"size:500"`
	ImageURL    string `gorm:"size:255"`

	// Relations - only one direction to avoid circular references
	Products []Product `gorm:"foreignKey:CategoryID"`
}

// Product represents a product in the store
type Product struct {
	gorm.Model
	SellerID         uint    `gorm:"not null"`
	CategoryID       uint    `gorm:"not null"`
	Name             string  `gorm:"size:200;not null"`
	Description      string  `gorm:"size:1000"`
	Price            float64 `gorm:"not null"`
	ImageURL         string  `gorm:"size:255"`
	IsActive         bool    `gorm:"default:true"`
	SatisfactionRate float64 `gorm:"default:0"` // Calculated from ProductReactions

	// Relations
	Seller            Seller             `gorm:"foreignKey:SellerID"`
	Category          Category           `gorm:"foreignKey:CategoryID"`
	ProductAttributes []ProductAttribute `gorm:"foreignKey:ProductID"`
	ProductStocks     []ProductStock     `gorm:"foreignKey:ProductID"`
	ProductReactions  []ProductReaction  `gorm:"foreignKey:ProductID"`
	ProductQuestions  []ProductQuestion  `gorm:"foreignKey:ProductID"`
	CartItems         []CartItem         `gorm:"foreignKey:ProductID"`
	OrderItems        []OrderItem        `gorm:"foreignKey:ProductID"`
	Reviews           []Review           `gorm:"foreignKey:ProductID"`
}

// ProductAttribute represents a key-value attribute for a product
type ProductAttribute struct {
	gorm.Model
	ProductID uint   `gorm:"not null;index"`
	Key       string `gorm:"size:100;not null"`
	Value     string `gorm:"size:255;not null"`

	// Relations - only one direction to avoid circular references
	Product Product `gorm:"foreignKey:ProductID"`
}

// ProductStock represents the stock of a product for a specific seller
type ProductStock struct {
	gorm.Model
	ProductID uint `gorm:"not null;index:idx_product_seller,unique"`
	SellerID  uint `gorm:"not null;index:idx_product_seller,unique"`
	Quantity  int  `gorm:"not null"`

	// Relations - only one direction to avoid circular references
	Product Product `gorm:"foreignKey:ProductID"`
	Seller  Seller  `gorm:"foreignKey:SellerID"`
}

// ProductReaction represents a like or dislike from a customer for a product
type ProductReaction struct {
	gorm.Model
	ProductID  uint `gorm:"not null;index:idx_product_customer,unique"`
	CustomerID uint `gorm:"not null;index:idx_product_customer,unique"`
	IsLike     bool `gorm:"not null"` // true for like, false for dislike

	// Relations - only one direction to avoid circular references
	Product  Product  `gorm:"foreignKey:ProductID"`
	Customer Customer `gorm:"foreignKey:CustomerID"`
}

// ProductQuestion represents a question asked by a customer about a product
type ProductQuestion struct {
	gorm.Model
	ProductID  uint   `gorm:"not null;index"`
	CustomerID uint   `gorm:"not null;index"`
	Question   string `gorm:"size:1000;not null"`

	// Relations
	Product        Product         `gorm:"foreignKey:ProductID"`
	Customer       Customer        `gorm:"foreignKey:CustomerID"`
	ProductAnswers []ProductAnswer `gorm:"foreignKey:QuestionID"`
}

// ProductAnswer represents an answer from a seller to a product question
type ProductAnswer struct {
	gorm.Model
	QuestionID uint   `gorm:"not null;index"`
	SellerID   uint   `gorm:"not null;index"`
	Answer     string `gorm:"size:1000;not null"`

	// Relations - only one direction to avoid circular references
	Question ProductQuestion `gorm:"foreignKey:QuestionID"`
	Seller   Seller          `gorm:"foreignKey:SellerID"`
}

// Review represents a review for a product
type Review struct {
	gorm.Model
	CustomerID uint   `gorm:"not null;index"`
	ProductID  uint   `gorm:"not null;index"`
	Content    string `gorm:"size:1000"`
	Rating     int    `gorm:"not null;check:rating >= 1 AND rating <= 5"`

	// Relations - only one direction to avoid circular references
	Customer Customer `gorm:"foreignKey:CustomerID"`
	Product  Product  `gorm:"foreignKey:ProductID"`
}
