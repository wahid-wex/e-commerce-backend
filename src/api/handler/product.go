package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/wahid-wex/e-commerce-backend/api/dto"
	_ "github.com/wahid-wex/e-commerce-backend/api/dto"
	"github.com/wahid-wex/e-commerce-backend/api/helper"
	_ "github.com/wahid-wex/e-commerce-backend/api/helper"
	"github.com/wahid-wex/e-commerce-backend/config"
	"github.com/wahid-wex/e-commerce-backend/constants"
	"github.com/wahid-wex/e-commerce-backend/services"
	"net/http"
	"strconv"
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
	req := new(dto.CreateUpdateProductRequest)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationError, err))
		return
	}

	res, err := h.service.Create(c, req)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}
	c.JSON(http.StatusCreated, helper.GenerateBaseResponse(res, true, 0))
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
	id, _ := strconv.Atoi(c.Params.ByName("id"))
	req := new(dto.CreateUpdateProductRequest)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationError, err))
		return
	}

	res, err := h.service.Update(c, id, req)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}
	c.JSON(http.StatusOK, helper.GenerateBaseResponse(res, true, 0))
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
	id, _ := strconv.Atoi(c.Params.ByName("id"))
	if id == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound,
			helper.GenerateBaseResponse(nil, false, helper.ValidationError))
		return
	}
	// Get customer ID from context (assuming it's set by auth middleware)
	customerID, exists := c.Get(constants.UserIdKey)
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &helper.BaseHttpResponse{
			Result: nil,
			Error:  "Unauthorized",
		})
		return
	}
	customerIdNumber := int(customerID.(float64))
	res, err := h.service.GetById(id, customerIdNumber)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}
	c.JSON(http.StatusOK, helper.GenerateBaseResponse(res, true, 0))
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

// AddToFavorite godoc
// @Summary Add a product to favorites
// @Description Add a product to customer's favorites
// @Tags Products
// @Accept json
// @produces json
// @Param Request body dto.AddRemoveToFavoriteRequest true "Add to favorites request"
// @Success 200 {object} helper.BaseHttpResponse{result=dto.FavoriteResponse} "Favorite response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v1/products/favorite [post]
// @Security AuthBearer
func (h *ProductService) AddToFavorite(c *gin.Context) {
	req := new(dto.AddRemoveToFavoriteRequest)
	err := c.ShouldBindJSON(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &helper.BaseHttpResponse{
			Result: nil,
			Error:  err.Error(),
		})
		return
	}

	// Get customer ID from context (assuming it's set by auth middleware)
	customerID, exists := c.Get(constants.UserIdKey)
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &helper.BaseHttpResponse{
			Result: nil,
			Error:  "Unauthorized",
		})
		return
	}

	customerIdNumber := uint(customerID.(float64))

	result, err := h.service.AddToFavorite(c.Request.Context(), customerIdNumber, req.ProductID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &helper.BaseHttpResponse{
			Result: nil,
			Error:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &helper.BaseHttpResponse{
		Result: result,
		Error:  nil,
	})
}

// RemoveFavorite godoc
// @Summary remove a product from favorites
// @Description remove a product from customer's favorites
// @Tags Products
// @Accept json
// @produces json
// @Param Request body dto.AddRemoveToFavoriteRequest true "Add to favorites request"
// @Success 200 {object} helper.BaseHttpResponse{result=dto.FavoriteResponse} "Favorite response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v1/products/favorite [delete]
// @Security AuthBearer
func (h *ProductService) RemoveFavorite(c *gin.Context) {
	req := new(dto.AddRemoveToFavoriteRequest)
	err := c.ShouldBindJSON(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &helper.BaseHttpResponse{
			Result: nil,
			Error:  err.Error(),
		})
		return
	}

	// Get customer ID from context (assuming it's set by auth middleware)
	customerID, exists := c.Get(constants.UserIdKey)
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &helper.BaseHttpResponse{
			Result: nil,
			Error:  "Unauthorized",
		})
		return
	}

	customerIdNumber := uint(customerID.(float64))

	result, err := h.service.RemoveFavorite(c.Request.Context(), customerIdNumber, req.ProductID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &helper.BaseHttpResponse{
			Result: nil,
			Error:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &helper.BaseHttpResponse{
		Result: result,
		Error:  nil,
	})
}

// LeaveReview godoc
// @Summary Add a product to favorites
// @Description Add a product to customer's favorites
// @Tags Products
// @Accept json
// @produces json
// @Param Request body dto.CreateReviewRequest true "Add to favorites request"
// @Success 200 {object} helper.BaseHttpResponse{result=nil} "Favorite response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v1/products/leave-review [post]
// @Security AuthBearer
func (h *ProductService) LeaveReview(c *gin.Context) {
	req := new(dto.CreateReviewRequest)
	err := c.ShouldBindJSON(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &helper.BaseHttpResponse{
			Result: nil,
			Error:  err.Error(),
		})
		return
	}

	// Get customer ID from context (assuming it's set by auth middleware)
	customerID, exists := c.Get(constants.UserIdKey)
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &helper.BaseHttpResponse{
			Result: nil,
			Error:  "Unauthorized",
		})
		return
	}

	customerIdNumber := uint(customerID.(float64))

	result, err := h.service.LeaveReview(c.Request.Context(), *req, customerIdNumber)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &helper.BaseHttpResponse{
			Result: nil,
			Error:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &helper.BaseHttpResponse{
		Result: result,
		Error:  nil,
	})
}

// CreateUpdateProductStock godoc
// @Summary Change Stock
// @Description Change or create stock of product
// @Tags Products
// @Accept json
// @produces json
// @Param Request body dto.CreateUpdateProductStockRequest true "Stock Change"
// @Success 200 {object} helper.BaseHttpResponse{result=nil} "Stock response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v1/products/change-stock [post]
// @Security AuthBearer
func (h *ProductService) CreateUpdateProductStock(c *gin.Context) {
	req := new(dto.CreateUpdateProductStockRequest)
	err := c.ShouldBindJSON(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &helper.BaseHttpResponse{
			Result: nil,
			Error:  err.Error(),
		})
		return
	}

	sellerID, exists := c.Get(constants.UserIdKey)
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &helper.BaseHttpResponse{
			Result: nil,
			Error:  "Unauthorized",
		})
		return
	}

	req.SellerID = uint(sellerID.(float64))
	result, err := h.service.AddUpdateProductQuantity(c.Request.Context(), *req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &helper.BaseHttpResponse{
			Result: nil,
			Error:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &helper.BaseHttpResponse{
		Result: result,
		Error:  nil,
	})
}
