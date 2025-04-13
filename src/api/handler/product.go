package handlers

import (
	"github.com/gin-gonic/gin"
	_ "github.com/wahid-wex/e-commerce-backend/api/dto"
	_ "github.com/wahid-wex/e-commerce-backend/api/helper"
	"github.com/wahid-wex/e-commerce-backend/config"
	"github.com/wahid-wex/e-commerce-backend/services"
)

type ProductService struct {
	service *services.ProductService
}

func NewProductService(cfg *config.Config) *ProductService {
	return &ProductService{service: services.NewProductService(cfg)}
}

// Create Products godoc
// @Summary Create a Product
// @Description Create a Product
// @Tags Products
// @Accept json
// @produces json
// @Param Request body dto.CreateUpdateProductRequest true "Create a Product"
// @Success 201 {object} helper.BaseHttpResponse{result=dto.ProductResponse} "Product response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v1/products/ [post]
// @Security AuthBearer
func (h *ProductService) Create(c *gin.Context) {
	Create(c, h.service.Create)
}

// Update Products godoc
// @Summary Update a Product
// @Description Update a Product
// @Tags Products
// @Accept json
// @produces json
// @Param id path int true "ID"
// @Param Request body dto.CreateUpdateProductRequest true "Update a Product"
// @Success 200 {object} helper.BaseHttpResponse{result=dto.ProductResponse} "Product response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v1/products/{id} [put]
// @Security AuthBearer
func (h *ProductService) Update(c *gin.Context) {
	Update(c, h.service.Update)
}

// Delete Products godoc
// @Summary Delete a Product
// @Description Delete a Product
// @Tags Products
// @Accept json
// @produces json
// @Param id path int true "Id"
// @Success 200 {object} helper.BaseHttpResponse "response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v1/products/{id} [delete]
// @Security AuthBearer
func (h *ProductService) Delete(c *gin.Context) {
	Delete(c, h.service.Delete)
}

// GetById godoc
// @Summary Get a Product
// @Description Get a Product
// @Tags Products
// @Accept json
// @produces json
// @Param id path int true "ID"
// @Success 200 {object} helper.BaseHttpResponse{result=dto.ProductResponse} "Product response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v1/products/{id} [get]
// @Security AuthBearer
func (h *ProductService) GetById(c *gin.Context) {
	GetById(c, h.service.GetById)
}

// GetByFilter godoc
// @Summary Get a Product
// @Description Get a Product
// @Tags Products
// @Accept json
// @produces json
// @Param Request body dto.PaginationInputWithFilter true "Request"
// @Success 200 {object} helper.BaseHttpResponse{result=dto.PagedList[dto.ProductResponse]} "Product response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v1/products/get-by-filter [post]
// @Security AuthBearer
func (h *ProductService) GetByFilter(c *gin.Context) {
	GetByFilter(c, h.service.GetByFilter)
}
