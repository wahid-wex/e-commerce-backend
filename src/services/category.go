package services

import (
	"context"
	"github.com/wahid-wex/e-commerce-backend/api/dto"
	"github.com/wahid-wex/e-commerce-backend/config"
	"github.com/wahid-wex/e-commerce-backend/data/db"
	"github.com/wahid-wex/e-commerce-backend/data/models"
	logging "github.com/wahid-wex/e-commerce-backend/logs"
)

type CategoryService struct {
	base *BaseService[models.Category, dto.CreateUpdateCategoryRequest, dto.CreateUpdateCategoryRequest, dto.CategoryResponse]
}

func NewCategoryService(cfg *config.Config) *CategoryService {
	return &CategoryService{
		base: &BaseService[models.Category, dto.CreateUpdateCategoryRequest, dto.CreateUpdateCategoryRequest, dto.CategoryResponse]{
			Database: db.GetDb(),
			Logger:   logging.NewLogger(cfg),
			Preloads: []preload{},
		},
	}
}

// Create
func (s *CategoryService) Create(ctx context.Context, req *dto.CreateUpdateCategoryRequest) (*dto.CategoryResponse, error) {
	return s.base.Create(ctx, req)
}

// Update
func (s *CategoryService) Update(ctx context.Context, id int, req *dto.CreateUpdateCategoryRequest) (*dto.CategoryResponse, error) {
	return s.base.Update(ctx, id, req)
}

// Delete
func (s *CategoryService) Delete(ctx context.Context, id int) error {
	return s.base.Delete(ctx, id)
}

// Get By Id
func (s *CategoryService) GetById(ctx context.Context, id int) (*dto.CategoryResponse, error) {
	return s.base.GetById(ctx, id)
}

// Get By Filter
func (s *CategoryService) GetByFilter(ctx context.Context, req *dto.PaginationInputWithFilter) (*dto.PagedList[dto.CategoryResponse], error) {
	return s.base.GetByFilter(ctx, req)
}
