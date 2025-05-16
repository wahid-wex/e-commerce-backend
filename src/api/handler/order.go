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

type OrderHandler struct {
	service *services.OrderService
}

func NewOrderHandler(cfg *config.Config) *OrderHandler {
	return &OrderHandler{service: services.NewOrderService(cfg)}
}

// MakeOrder godoc
// @Summary Finalize Order
// @Description Finalize Order
// @Tags Order
// @Accept json
// @produces json
// @Param Request body dto.CheckoutRequest true "Finalize Order"
// @Success 200 {object} helper.BaseHttpResponse{result=dto.CheckoutResponse} "Product response"
// @Failure 400 {object} helper.BaseHttpResponse "Bad request"
// @Router /v1/order/ [post]
// @Security AuthBearer
func (h *OrderHandler) MakeOrder(c *gin.Context) {
	req := new(dto.CheckoutRequest)
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
	res, err := h.service.Checkout(c.Request.Context(), customerIdNumber, req.ShippingAddress, "Mellat")
	if err != nil {
		c.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err),
			helper.GenerateBaseResponseWithError(nil, false, helper.InternalError, err))
		return
	}
	c.JSON(http.StatusCreated, helper.GenerateBaseResponse(res, true, 0))
}
