package router

import (
	"github.com/gin-gonic/gin"
	handlers "github.com/wahid-wex/e-commerce-backend/api/handler"
	middlewares "github.com/wahid-wex/e-commerce-backend/api/middleware"
	"github.com/wahid-wex/e-commerce-backend/config"
)

const GetByFilterExp string = "/get-by-filter"

func Category(r *gin.RouterGroup, cfg *config.Config) {
	h := handlers.NewCategoryHandler(cfg)

	r.POST("/", middlewares.Authorization([]string{"admin"}), h.Create)
	r.PUT("/:id", middlewares.Authorization([]string{"admin"}), h.Update)
	r.DELETE("/:id", middlewares.Authorization([]string{"admin"}), h.Delete)
	r.GET("/:id", h.GetById)
	r.POST(GetByFilterExp, h.GetByFilter)
}
