package services

import (
	"errors"
	"github.com/wahid-wex/e-commerce-backend/api/dto"
	"github.com/wahid-wex/e-commerce-backend/api/error_handler"
	"github.com/wahid-wex/e-commerce-backend/config"
	"github.com/wahid-wex/e-commerce-backend/data/db"
	"github.com/wahid-wex/e-commerce-backend/data/models"
	logging "github.com/wahid-wex/e-commerce-backend/logs"
	"gorm.io/gorm"
)

type CartService struct {
	base *BaseService[models.Cart, dto.CartItemRequest, dto.CartItemRequest, dto.CartResponse]
}

func NewCartService(cfg *config.Config) *CartService {
	return &CartService{
		base: &BaseService[models.Cart, dto.CartItemRequest, dto.CartItemRequest, dto.CartResponse]{
			Database: db.GetDb(),
			Logger:   logging.NewLogger(cfg),
			Preloads: []preload{},
		},
	}
}

// AddToCart adds a product to the user's cart
func (s *CartService) AddToCart(customerID uint, productID uint) (*dto.CartResponse, error) {
	// Check if product exists
	var product models.Product
	if err := s.base.Database.First(&product, productID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &error_handler.ServiceError{EndUserMessage: "Product not found"}
		}
		return nil, err
	}

	// Get or create cart
	var cart models.Cart
	err := s.base.Database.Where("customer_id = ?", customerID).First(&cart).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new cart
			cart = models.Cart{CustomerID: customerID}
			if err := s.base.Database.Create(&cart).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	// Check if product already exists in cart
	var cartItem models.CartItem
	err = s.base.Database.Where("cart_id = ? AND product_id = ?", cart.ID, productID).First(&cartItem).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new cart item
			cartItem = models.CartItem{
				CartID:    cart.ID,
				ProductID: productID,
				Quantity:  1,
				Price:     product.Price,
			}
			if err := s.base.Database.Create(&cartItem).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		// Update quantity
		cartItem.Quantity++
		if err := s.base.Database.Save(&cartItem).Error; err != nil {
			return nil, err
		}
	}

	return s.getCartResponse(cart.ID)
}

// RemoveFromCart removes a product from the user's cart
func (s *CartService) RemoveFromCart(customerID uint, productID uint) (*dto.CartResponse, error) {
	// Get cart
	var cart models.Cart
	if err := s.base.Database.Where("customer_id = ?", customerID).First(&cart).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &error_handler.ServiceError{EndUserMessage: "Cart not found"}
		}
		return nil, err
	}

	// Get cart item
	var cartItem models.CartItem
	if err := s.base.Database.Where("cart_id = ? AND product_id = ?", cart.ID, productID).First(&cartItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &error_handler.ServiceError{EndUserMessage: "Product not found in cart"}
		}
		return nil, err
	}

	// Update or delete cart item
	if cartItem.Quantity > 1 {
		cartItem.Quantity--
		if err := s.base.Database.Save(&cartItem).Error; err != nil {
			return nil, err
		}
	} else {
		if err := s.base.Database.Delete(&cartItem).Error; err != nil {
			return nil, err
		}
	}

	return s.getCartResponse(cart.ID)
}

// getCartResponse returns the cart response with items and total price
func (s *CartService) getCartResponse(cartID uint) (*dto.CartResponse, error) {
	var cartItems []models.CartItem
	if err := s.base.Database.
		Preload("Product").
		Where("cart_id = ?", cartID).
		Find(&cartItems).Error; err != nil {
		return nil, err
	}

	response := &dto.CartResponse{
		Items:      make([]dto.CartItemResponse, 0, len(cartItems)),
		TotalPrice: 0,
	}

	for _, item := range cartItems {
		response.Items = append(response.Items, dto.CartItemResponse{
			ProductID: item.ProductID,
			Name:      item.Product.Name,
			Price:     item.Price,
			ImageURL:  item.Product.ImageURL,
			Quantity:  item.Quantity,
		})
		response.TotalPrice += item.Price * float64(item.Quantity)
	}

	return response, nil
}

func (s *CartService) GetCart(customerID int) (*dto.CartResponse, error) {
	// Get or create cart
	var cart models.Cart
	err := s.base.Database.
		Preload("CartItems").
		Where("customer_id = ?", customerID).
		First(&cart).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return s.getCartResponse(cart.ID)
}
