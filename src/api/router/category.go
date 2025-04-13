package router

import (
	"github.com/gin-gonic/gin"
	handlers "github.com/wahid-wex/e-commerce-backend/api/handler"
	"github.com/wahid-wex/e-commerce-backend/config"
)

const GetByFilterExp string = "/get-by-filter"

func Category(r *gin.RouterGroup, cfg *config.Config) {
	h := handlers.NewCategoryHandler(cfg)

	r.POST("/", h.Create)
	r.PUT("/:id", h.Update)
	r.DELETE("/:id", h.Delete)
	r.GET("/:id", h.GetById)
	r.POST(GetByFilterExp, h.GetByFilter)
}
