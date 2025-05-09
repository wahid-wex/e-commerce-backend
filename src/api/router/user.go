package router

import (
	"github.com/gin-gonic/gin"
	handlers "github.com/wahid-wex/e-commerce-backend/api/handler"
	"github.com/wahid-wex/e-commerce-backend/config"
)

func User(router *gin.RouterGroup, cfg *config.Config) {
	ch := handlers.NewCustomersHandler(cfg)
	cs := handlers.NewSellersHandler(cfg)

	router.POST("/send-customer-otp", ch.SendCustomerOtp)
	router.POST("/login-customer-by-username", ch.LoginCustomerByUsername)
	router.POST("/register-customer-by-username", ch.RegisterCustomerByUsername)
	router.POST("/login-customer-by-mobile", ch.RegisterLoginCustomerByMobileNumber)
	router.POST("/refresh-customer-token", ch.RefreshToken)

	router.POST("/send-seller-otp", cs.SendSellerOtp)
	router.POST("/login-seller-by-username", cs.LoginSellerByUsername)
	router.POST("/register-seller-by-username", cs.RegisterSellerByUsername)
	router.POST("/login-seller-by-mobile", cs.RegisterLoginSellerByMobileNumber)
	router.POST("/refresh-seller-token", cs.RefreshToken)

}
