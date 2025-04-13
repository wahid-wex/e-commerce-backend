package services

import (
	"context"
	"github.com/wahid-wex/e-commerce-backend/api/dto"
	"github.com/wahid-wex/e-commerce-backend/config"
	"github.com/wahid-wex/e-commerce-backend/data/db"
	"github.com/wahid-wex/e-commerce-backend/data/models"
	logging "github.com/wahid-wex/e-commerce-backend/logs"
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
func (s *ProductService) Create(ctx context.Context, req *dto.CreateUpdateProductRequest) (*dto.ProductResponse, error) {
	return s.base.Create(ctx, req)
}

// Update
func (s *ProductService) Update(ctx context.Context, id int, req *dto.CreateUpdateProductRequest) (*dto.ProductResponse, error) {
	return s.base.Update(ctx, id, req)
}

// Delete
func (s *ProductService) Delete(ctx context.Context, id int) error {
	return s.base.Delete(ctx, id)
}

// Get By Id
func (s *ProductService) GetById(ctx context.Context, id int) (*dto.ProductResponse, error) {
	return s.base.GetById(ctx, id)
}

// Get By Filter
func (s *ProductService) GetByFilter(ctx context.Context, req *dto.PaginationInputWithFilter) (*dto.PagedList[dto.ProductResponse], error) {
	return s.base.GetByFilter(ctx, req)
}
