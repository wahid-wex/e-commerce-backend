package dto

type CartItemRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
}

type CartItemResponse struct {
	ProductID uint    `json:"productId"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	ImageURL  string  `json:"imageUrl"`
	Quantity  int     `json:"quantity"`
}

type CartResponse struct {
	Items      []CartItemResponse `json:"items"`
	TotalPrice float64            `json:"totalPrice"`
}
