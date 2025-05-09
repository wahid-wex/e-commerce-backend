package router

import (
	"github.com/gin-gonic/gin"
	handlers "github.com/wahid-wex/e-commerce-backend/api/handler"
	middlewares "github.com/wahid-wex/e-commerce-backend/api/middleware"
	"github.com/wahid-wex/e-commerce-backend/config"
)

func Product(r *gin.RouterGroup, cfg *config.Config) {
	h := handlers.NewProductService(cfg)

	r.POST("/", middlewares.Authorization([]string{"seller"}), h.Create)
	r.PUT("/:id", middlewares.Authorization([]string{"seller"}), h.Update)
	r.DELETE("/:id", middlewares.Authorization([]string{"seller"}), h.Delete)
	r.POST("/change-stock", middlewares.Authorization([]string{"seller"}), h.CreateUpdateProductStock)
	r.GET("/:id", h.GetById)
	r.POST("/get-by-filter", h.GetByFilter)
	r.POST("/favorite", middlewares.Authorization([]string{"customer"}), h.AddToFavorite)
	r.DELETE("/favorite", middlewares.Authorization([]string{"customer"}), h.RemoveFavorite)
	r.POST("/leave-review", middlewares.Authorization([]string{"customer"}), h.LeaveReview)

}
