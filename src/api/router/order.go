package router

import (
	"github.com/gin-gonic/gin"
	handlers "github.com/wahid-wex/e-commerce-backend/api/handler"
	middlewares "github.com/wahid-wex/e-commerce-backend/api/middleware"
	"github.com/wahid-wex/e-commerce-backend/config"
)

func Order(r *gin.RouterGroup, cfg *config.Config) {
	h := handlers.NewOrderHandler(cfg)

	r.POST("/", middlewares.Authorization([]string{"customer"}), h.MakeOrder)
}
