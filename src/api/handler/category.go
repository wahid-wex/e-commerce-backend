package handlers

import (
	"github.com/gin-gonic/gin"
	_ "github.com/wahid-wex/e-commerce-backend/api/dto"
	_ "github.com/wahid-wex/e-commerce-backend/api/helper"
	"github.com/wahid-wex/e-commerce-backend/config"
	"github.com/wahid-wex/e-commerce-backend/services"
)

type CategoryHandler struct {
	service *services.CategoryService
}

func NewCategoryHandler(cfg *config.Config) *CategoryHandler {
	return &CategoryHandler{service: services.NewCategoryService(cfg)}
}

// Create Categories godoc
// @Summary Create a Category
// @Description Create a Category
// @Tags Categories
// @Accept json
// @produces json
// @Param Request body dto.CreateUpdateCategoryRequest true "Create a Category"
// @Success 201 {object} helper.BaseHttpResponse{result=dto.CategoryResponse} "Category response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v1/categories/ [post]
func (h *CategoryHandler) Create(c *gin.Context) {
	Create(c, h.service.Create)
}

// Update Categories godoc
// @Summary Update a Category
// @Description Update a Category
// @Tags Categories
// @Accept json
// @produces json
// @Param id path int true "ID"
// @Param Request body dto.CreateUpdateCategoryRequest true "Update a Category"
// @Success 200 {object} helper.BaseHttpResponse{result=dto.CategoryResponse} "Category response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v1/categories/{id} [put]
func (h *CategoryHandler) Update(c *gin.Context) {
	Update(c, h.service.Update)
}

// Delete Categories godoc
// @Summary Delete a Category
// @Description Delete a Category
// @Tags Categories
// @Accept json
// @produces json
// @Param id path int true "Id"
// @Success 200 {object} helper.BaseHttpResponse "response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v1/categories/{id} [delete]
func (h *CategoryHandler) Delete(c *gin.Context) {
	Delete(c, h.service.Delete)
}

// GetById godoc
// @Summary Get a Category
// @Description Get a Category
// @Tags Categories
// @Accept json
// @produces json
// @Param id path int true "ID"
// @Success 200 {object} helper.BaseHttpResponse{result=dto.CategoryResponse} "Category response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v1/categories/{id} [get]
func (h *CategoryHandler) GetById(c *gin.Context) {
	GetById(c, h.service.GetById)
}

// GetByFilter godoc
// @Summary Get a Category
// @Description Get a Category
// @Tags Categories
// @Accept json
// @produces json
// @Param Request body dto.PaginationInputWithFilter true "Request"
// @Success 200 {object} helper.BaseHttpResponse{result=dto.PagedList[dto.CategoryResponse]} "Category response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v1/categories/get-by-filter [post]
func (h *CategoryHandler) GetByFilter(c *gin.Context) {
	GetByFilter(c, h.service.GetByFilter)
}
