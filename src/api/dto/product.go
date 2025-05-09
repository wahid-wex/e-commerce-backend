package dto

type CreateUpdateCategoryRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=20"`
	Description string `json:"description" binding:"required,min=10,max=300"`
	ImageURL    string `json:"imageUrl" binding:"required,min=3,max=100"`
}

type ProductAttributeDTO struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
}

type CreateUpdateProductRequest struct {
	CategoryID        uint                  `json:"category_id" binding:"required"`
	Name              string                `json:"name" binding:"required,min=1,max=200"`
	Description       string                `json:"description,omitempty"`
	Price             float64               `json:"price" binding:"required,gt=0"`
	ImageURL          string                `json:"image_url,omitempty"`
	IsActive          *bool                 `json:"is_active,omitempty"`
	ProductAttributes []ProductAttributeDTO `json:"attributes,omitempty"`
	ProductStocks     int                   `json:"stocks,omitempty"`
}

type ProductResponse struct {
	Name              string                     `json:"name"`
	Description       string                     `json:"description"`
	Price             float64                    `json:"price"`
	ImageURL          string                     `json:"imageUrl"`
	IsActive          bool                       `json:"isActive"`
	SatisfactionRate  float64                    `json:"satisfactionRate"`
	Category          CategoryResponse           `json:"category"`
	ProductAttributes []ProductAttributeResponse `json:"productAttributes"`
	ProductStocks     []ProductStockResponse     `json:"productStocks"`
	IsFavorite        bool                       `json:"isFavorite"`
	CartItems         []CartItemResponse         `json:"cartItems"`
	Reviews           []ReviewResponse           `json:"reviews"`
}

// AddRemoveToFavoriteRequest represents the request to add a product to favorites
type AddRemoveToFavoriteRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
}

type CreateReviewRequest struct {
	ProductID uint   `json:"product-id" binding:"required"`
	Content   string `json:"content" binding:"required,min=1,max=1000"`
	Rating    int    `json:"rating,omitempty"`
}

type CreateUpdateProductStockRequest struct {
	ProductID uint `json:"productId" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,gt=0"`
	SellerID  uint `json:"-" swaggerignore:"true"`
}

type CategoryResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"imageUrl"`
}

type ProductAttributeResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ProductStockResponse struct {
	SellerID uint `json:"sellerId"`
	Quantity int  `json:"quantity"`
}

type ReviewResponse struct {
	CustomerID uint   `json:"customerId"`
	Content    string `json:"content"`
	Rating     int    `json:"rating"`
}

type FavoriteResponse struct {
	CustomerID uint `json:"customerId"`
}
