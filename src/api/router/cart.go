package router

import (
	"github.com/gin-gonic/gin"
	handlers "github.com/wahid-wex/e-commerce-backend/api/handler"
	middlewares "github.com/wahid-wex/e-commerce-backend/api/middleware"
	"github.com/wahid-wex/e-commerce-backend/config"
)

func Cart(r *gin.RouterGroup, cfg *config.Config) {
	h := handlers.NewCartHandler(cfg)

	r.POST("/", middlewares.Authorization([]string{"customer"}), h.AddToCart)
	r.DELETE("/", middlewares.Authorization([]string{"customer"}), h.RemoveFromCart)
	r.GET("/", middlewares.Authorization([]string{"customer"}), h.GetCart)
}
