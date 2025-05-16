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
)

type CartHandler struct {
	service *services.CartService
}

func NewCartHandler(cfg *config.Config) *CartHandler {
	return &CartHandler{service: services.NewCartService(cfg)}
}

// RemoveFromCart godoc
// @Summary Remove From Cart
// @Description Remove From Cart
// @Tags Cart
// @Accept json
// @produces json
// @Param Request body dto.CartItemRequest true "Remove Product to Cart"
// @Success 201 {object} helper.BaseHttpResponse{result=dto.CartResponse} "Product response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v1/cart/ [delete]
// @Security AuthBearer
func (h *CartHandler) RemoveFromCart(c *gin.Context) {
	req := new(dto.CartItemRequest)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationError, err))
		return
	}
	customerID, exists := c.Get(constants.UserIdKey)
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &helper.BaseHttpResponse{
			Result: nil,
			Error:  "Unauthorized",
		})
		return
	}
	customerIdNumber := uint(customerID.(float64))
	res, err := h.service.RemoveFromCart(customerIdNumber, req.ProductID)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}
	c.JSON(http.StatusCreated, helper.GenerateBaseResponse(res, true, 0))
}

// AddToCart godoc
// @Summary Add Product to Cart
// @Description Add Product to Cart
// @Tags Cart
// @Accept json
// @produces json
// @Param Request body dto.CartItemRequest true "Add Product to Cart"
// @Success 201 {object} helper.BaseHttpResponse{result=dto.CartResponse} "Product response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v1/cart/ [post]
// @Security AuthBearer
func (h *CartHandler) AddToCart(c *gin.Context) {
	req := new(dto.CartItemRequest)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			helper.GenerateBaseResponseWithValidationError(nil, false, helper.ValidationError, err))
		return
	}
	customerID, exists := c.Get(constants.UserIdKey)
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &helper.BaseHttpResponse{
			Result: nil,
			Error:  "Unauthorized",
		})
		return
	}
	customerIdNumber := uint(customerID.(float64))
	res, err := h.service.AddToCart(customerIdNumber, req.ProductID)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}
	c.JSON(http.StatusCreated, helper.GenerateBaseResponse(res, true, 0))
}

// GetCart godoc
// @Summary Get Cart
// @Description Get Cart Detail
// @Tags Cart
// @Accept json
// @produces json
// @Success 200 {object} helper.BaseHttpResponse{result=dto.CartResponse} "Cart Response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v1/cart/ [get]
// @Security AuthBearer
func (h *CartHandler) GetCart(c *gin.Context) {
	customerID, exists := c.Get(constants.UserIdKey)
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, &helper.BaseHttpResponse{
			Result: nil,
			Error:  "Unauthorized",
		})
		return
	}
	customerIdNumber := int(customerID.(float64))
	res, err := h.service.GetCart(customerIdNumber)
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}
	c.JSON(http.StatusCreated, helper.GenerateBaseResponse(res, true, 0))
}
