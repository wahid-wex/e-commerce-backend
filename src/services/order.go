package services

import (
	"context"
	"github.com/wahid-wex/e-commerce-backend/api/dto"
	"github.com/wahid-wex/e-commerce-backend/api/error_handler"
	"github.com/wahid-wex/e-commerce-backend/config"
	"github.com/wahid-wex/e-commerce-backend/data/db"
	"github.com/wahid-wex/e-commerce-backend/data/models"
	logging "github.com/wahid-wex/e-commerce-backend/logs"
	"time"

	"gorm.io/gorm"
)

type OrderService struct {
	DB     *gorm.DB
	Logger logging.Logger
}

func NewOrderService(cfg *config.Config) *OrderService {
	return &OrderService{
		DB:     db.GetDb(),
		Logger: logging.NewLogger(cfg),
	}
}

func (s *OrderService) Checkout(ctx context.Context, customerID uint, shippingAddress string, paymentMethod string) (dto.CheckoutResponse, error) {
	var cart models.Cart
	if err := s.DB.Preload("CartItems.Product.ProductStocks").Where("customer_id = ?", customerID).First(&cart).Error; err != nil {
		return dto.CheckoutResponse{}, &error_handler.ServiceError{EndUserMessage: error_handler.RecordNotFound}
	}

	if len(cart.CartItems) == 0 {
		return dto.CheckoutResponse{}, &error_handler.ServiceError{EndUserMessage: "Empty Cart"}
	}

	totalAmount := 0.0
	for _, item := range cart.CartItems {
		var stock models.ProductStock
		err := s.DB.Where("product_id = ?", item.ProductID).First(&stock).Error
		if err != nil {
			return dto.CheckoutResponse{}, &error_handler.ServiceError{EndUserMessage: "Product stock not found"}
		}
		if stock.Quantity < item.Quantity {
			return dto.CheckoutResponse{}, &error_handler.ServiceError{EndUserMessage: "not enough stock for product: " + item.Product.Name}
		}
		totalAmount += float64(item.Quantity) * item.Price
	}

	order := models.Order{
		CustomerID:      customerID,
		SellerID:        cart.CartItems[0].Product.ProductStocks[0].SellerID, // فرض کردیم همه محصولات سبد برای یک فروشنده هستن
		TotalAmount:     totalAmount,
		Status:          models.OrderStatusPending,
		ShippingAddress: shippingAddress,
		OrderDate:       time.Now(),
	}

	tx := s.DB.WithContext(ctx).Begin()
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		s.Logger.Error(logging.Postgres, logging.Insert, err.Error(), nil)
		return dto.CheckoutResponse{}, err
	}

	for _, item := range cart.CartItems {
		orderItem := models.OrderItem{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			s.Logger.Error(logging.Postgres, logging.Insert, err.Error(), nil)
			return dto.CheckoutResponse{}, err
		}

		if err := tx.Model(&models.ProductStock{}).
			Where("product_id = ? AND seller_id = ?", item.ProductID, order.SellerID).
			Update("quantity", gorm.Expr("quantity - ?", item.Quantity)).Error; err != nil {
			tx.Rollback()
			s.Logger.Error(logging.Postgres, logging.Insert, err.Error(), nil)
			return dto.CheckoutResponse{}, err
		}
	}

	payment := models.Payment{
		OrderID:       order.ID,
		Amount:        totalAmount,
		Status:        models.PaymentStatusPending,
		PaymentMethod: paymentMethod,
		PaymentDate:   time.Now(),
	}
	if err := tx.Create(&payment).Error; err != nil {
		tx.Rollback()
		s.Logger.Error(logging.Postgres, logging.Insert, err.Error(), nil)
		return dto.CheckoutResponse{}, err
	}

	if err := tx.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err != nil {
		tx.Rollback()
		s.Logger.Error(logging.Postgres, logging.Insert, err.Error(), nil)
		return dto.CheckoutResponse{}, err
	}

	tx.Commit()
	return dto.CheckoutResponse{PaymentGatewayUrl: "https://mock.com/payment/gateway"}, nil
}
