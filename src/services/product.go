package services

import (
	"context"
	"errors"
	"github.com/wahid-wex/e-commerce-backend/api/dto"
	"github.com/wahid-wex/e-commerce-backend/api/error_handler"
	"github.com/wahid-wex/e-commerce-backend/common"
	"github.com/wahid-wex/e-commerce-backend/config"
	"github.com/wahid-wex/e-commerce-backend/data/db"
	"github.com/wahid-wex/e-commerce-backend/data/models"
	logging "github.com/wahid-wex/e-commerce-backend/logs"
	"gorm.io/gorm"
)

type ProductService struct {
	base *BaseService[models.Product, dto.CreateUpdateProductRequest, dto.CreateUpdateProductRequest, dto.ProductResponse]
}

func NewProductService(cfg *config.Config) *ProductService {
	return &ProductService{
		base: &BaseService[models.Product, dto.CreateUpdateProductRequest, dto.CreateUpdateProductRequest, dto.ProductResponse]{
			Database: db.GetDb(),
			Logger:   logging.NewLogger(cfg),
			Preloads: []preload{},
		},
	}
}

// Create
func (s *ProductService) Create(ctx context.Context, req *dto.CreateUpdateProductRequest) (bool, error) {
	model := &models.Product{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		ImageURL:    req.ImageURL,
		IsActive:    true,
	}
	tx := s.base.Database.WithContext(ctx).Begin()
	err := tx.
		Create(model).
		Error
	if err != nil {
		tx.Rollback()
		s.base.Logger.Error(logging.Postgres, logging.Insert, err.Error(), nil)
		return false, err
	}
	tx.Commit()
	return true, nil
}

// Update
func (s *ProductService) Update(ctx context.Context, id int, req *dto.CreateUpdateProductRequest) (bool, error) {

	updateMap := map[string]interface{}{
		"CategoryID":  req.CategoryID,
		"Name":        req.Name,
		"Description": req.Description,
		"Price":       req.Price,
		"ImageURL":    req.ImageURL,
		"IsActive":    true,
	}
	snakeMap := map[string]interface{}{}
	for k, v := range updateMap {
		snakeMap[common.ToSnakeCase(k)] = v
	}
	model := new(models.Product)
	tx := s.base.Database.WithContext(ctx).Begin()
	if err := tx.Model(model).
		Where("id = ?", id).
		Updates(snakeMap).
		Error; err != nil {
		tx.Rollback()
		s.base.Logger.Error(logging.Postgres, logging.Update, err.Error(), nil)
		return false, err
	}
	tx.Commit()
	return true, nil
}

// Delete
func (s *ProductService) Delete(ctx context.Context, id int) error {
	result := s.base.Database.WithContext(ctx).Delete(&models.Product{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		s.base.Logger.Error(logging.Postgres, logging.Delete, error_handler.RecordNotFound, nil)
		return &error_handler.ServiceError{EndUserMessage: error_handler.RecordNotFound}
	}
	return nil
}

// Get By Id
func (s *ProductService) GetById(productId int, customerId int) (*dto.ProductResponse, error) {
	ID := uint(productId)
	product := new(models.Product)
	err := s.base.Database.
		Preload("Category").
		Preload("ProductAttributes").
		Preload("ProductStocks").
		Preload("CartItems").
		Preload("OrderItems").
		Preload("Reviews").
		Where("id = ?", ID).
		First(product).
		Error
	if err != nil {
		return nil, err
	}

	var favorite models.Favorite

	productResponse, err := common.TypeConverter[dto.ProductResponse](product)
	if err != nil {
		return nil, err
	}

	err = s.base.Database.
		Where("product_id = ? AND customer_id = ?", ID, customerId).
		First(&favorite).
		Error

	if err == nil {
		productResponse.IsFavorite = true
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		productResponse.IsFavorite = false
	}
	return productResponse, nil
}

// Get By Filter
func (s *ProductService) GetByFilter(ctx context.Context, req *dto.PaginationInputWithFilter) (*dto.PagedList[dto.ProductResponse], error) {
	return s.base.GetByFilter(ctx, req)
}

// AddToFavorite adds a product to customer's favorites
func (s *ProductService) AddToFavorite(ctx context.Context, customerID uint, productID uint) (bool, error) {
	product := new(models.Product)
	err := s.base.Database.
		Where("id = ?", productID).
		First(product).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, &error_handler.ServiceError{
				EndUserMessage: "record not found",
			}
		} else {
			return false, &error_handler.ServiceError{
				EndUserMessage: "internal server error",
			}
		}
	}

	// Check if already favorite
	var existingFavorite models.Favorite
	result := s.base.Database.WithContext(ctx).
		Where("customer_id = ? AND product_id = ?", customerID, productID).
		First(&existingFavorite)

	if result.Error == nil {
		return false, &error_handler.ServiceError{
			EndUserMessage: "Product is already in favorites",
		}
	}

	// Create new favorite
	favorite := &models.Favorite{
		CustomerID: customerID,
		ProductID:  productID,
	}

	tx := s.base.Database.WithContext(ctx).Begin()
	if err := tx.Create(favorite).Error; err != nil {
		tx.Rollback()
		s.base.Logger.Error(logging.Postgres, logging.Insert, err.Error(), nil)
		return false, err
	}
	tx.Commit()

	return true, nil
}

func (s *ProductService) RemoveFavorite(ctx context.Context, customerID uint, productID uint) (bool, error) {
	var existingFavorite models.Favorite

	err := s.base.Database.WithContext(ctx).
		Unscoped().
		Where("customer_id = ? AND product_id = ?", customerID, productID).
		First(&existingFavorite).Error

	if err != nil {
		s.base.Logger.Error(logging.Postgres, logging.Delete, error_handler.RecordNotFound, nil)
		return false, err
	}

	err = s.base.Database.WithContext(ctx).
		Unscoped().
		Where("customer_id = ? AND product_id = ?", customerID, productID).
		Delete(&models.Favorite{}).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *ProductService) LeaveReview(ctx context.Context, reviewInputs dto.CreateReviewRequest, CustomerID uint) (bool, error) {
	review := &models.Review{
		CustomerID: CustomerID,
		ProductID:  reviewInputs.ProductID,
		Content:    reviewInputs.Content,
		Rating:     reviewInputs.Rating,
	}

	tx := s.base.Database.WithContext(ctx).Begin()
	if err := tx.Create(review).Error; err != nil {
		tx.Rollback()
		s.base.Logger.Error(logging.Postgres, logging.Insert, err.Error(), nil)
		return false, err
	}
	tx.Commit()
	return true, nil
}

func (s *ProductService) AddUpdateProductQuantity(ctx context.Context, req dto.CreateUpdateProductStockRequest) (bool, error) {

	// check exist product
	var productCount int64
	err := s.base.Database.
		Model(&models.Product{}).
		Where("id = ?", req.ProductID).
		Count(&productCount).
		Error

	if err != nil {
		s.base.Logger.Error(logging.Postgres, logging.Update, err.Error(), nil)
		return false, err
	}

	if productCount == 0 {
		return false, &error_handler.ServiceError{EndUserMessage: error_handler.RecordNotFound}
	}

	// check exist product stock
	var productStock models.ProductStock
	err = s.base.Database.
		WithContext(ctx).
		Where("product_id = ?", req.ProductID).
		First(&productStock).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { // stock doesnt exist
			newStock := models.ProductStock{
				ProductID: req.ProductID,
				Quantity:  req.Quantity,
				SellerID:  req.SellerID,
			}
			err = s.base.Database.
				WithContext(ctx).
				Create(&newStock).
				Error

			if err != nil {
				s.base.Logger.Error(logging.Postgres, logging.Update, err.Error(), nil)
				return false, err
			}

			return true, nil
		}

		// another error
		s.base.Logger.Error(logging.Postgres, logging.Update, err.Error(), nil)
		return false, err
	}

	// plus one quantity
	productStock.Quantity = req.Quantity
	err = s.base.Database.
		WithContext(ctx).
		Save(&productStock).
		Error

	if err != nil {
		s.base.Logger.Error(logging.Postgres, logging.Update, err.Error(), nil)
		return false, err
	}

	return true, nil
}
