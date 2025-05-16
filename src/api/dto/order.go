package dto

type CheckoutRequest struct {
	ShippingAddress string `json:"shippingAddress" binding:"required"`
}

type CheckoutResponse struct {
	PaymentGatewayUrl string `json:"paymentGatewayUrl" binding:"required"`
}
