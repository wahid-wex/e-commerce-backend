package dto

type CreateUpdateCategoryRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=20"`
	Description string `json:"description" binding:"required,min=10,max=300"`
	ImageURL    string `json:"imageUrl" binding:"required,min=3,max=100"`
}

type CategoryResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"imageUrl"`
}

type CreateUpdateProductRequest struct {
	Name                string  `json:"name" binding:"required,alpha,min=3,max=20"`
	Description         string  `json:"description" binding:"required,alpha,min=30,max=200"`
	Price               float64 `json:"price" binding:"required,numeric"`
	ImageURL            string  `json:"imageUrl" binding:"required,alpha,min=3,max=100"`
	ProductAttributesID []int   `json:"productAttributes" binding:"required"`
	ProductStocksID     []int   `json:"productStocks" binding:"required"`
	SellerID            int     `json:"seller" binding:"required,numeric"`
	CategoryID          int     `json:"category" binding:"required,numeric"`
}

type ProductResponse struct {
	Name             string  `json:"name"`
	Description      string  `json:"description"`
	Price            float64 `json:"price"`
	ImageURL         string  `json:"imageUrl"`
	IsActive         bool    `json:"isActive"`
	SatisfactionRate float64 `json:"satisfactionRate"`
}
