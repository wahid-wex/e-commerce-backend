package dto

type GetOtpRequest struct {
	MobileNumber string `json:"mobileNumber" binding:"required,mobile,min=11,max=11"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type TokenDetail struct {
	AccessToken            string `json:"accessToken"`
	RefreshToken           string `json:"refreshToken"`
	AccessTokenExpireTime  int64  `json:"accessTokenExpireTime"`
	RefreshTokenExpireTime int64  `json:"refreshTokenExpireTime"`
}

type RegisterCustomerByUsernameRequest struct {
	FirstName       string `json:"firstName" binding:"required,min=3"`
	LastName        string `json:"lastName" binding:"required,min=3"`
	Username        string `json:"username" binding:"required,min=5"`
	Email           string `json:"email" binding:"min=6,email"`
	Password        string `json:"password" binding:"required,password,min=6"`
	PostalCode      string `json:"postalCode" binding:"required"`
	Phone           string `json:"phone" binding:"required"`
	CardNumber      string `json:"cardNumber" binding:"required"`
	ShippingAddress string `json:"shippingAddress" binding:"required"`
}

type RegisterSellerByUsernameRequest struct {
	Username    string `json:"username" binding:"required,min=5"`
	Email       string `json:"email" binding:"min=6,email"`
	Password    string `json:"password" binding:"required,password,min=6"`
	StoreName   string `json:"storeName" binding:"required,min=3"`
	NationalID  string `json:"nationalId" binding:"required"`
	Address     string `json:"address" binding:"required,min=5"`
	Phone       string `json:"phone" binding:"required,min=5"`
	Description string `json:"description"`
	Logo        string `json:"logo"`
}

type RegisterLoginByMobileRequest struct {
	MobileNumber string `json:"mobileNumber" binding:"required,mobile,min=11,max=11"`
	Otp          string `json:"otp" binding:"required,min=6,max=6"`
}

type LoginByUsernameRequest struct {
	Username string `json:"username" binding:"required,min=5"`
	Password string `json:"password" binding:"required,min=6"`
}
