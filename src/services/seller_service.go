package services

import (
	"github.com/wahid-wex/e-commerce-backend/api/dto"
	"github.com/wahid-wex/e-commerce-backend/api/error_handler"
	"github.com/wahid-wex/e-commerce-backend/common"
	"github.com/wahid-wex/e-commerce-backend/config"
	"github.com/wahid-wex/e-commerce-backend/data/db"
	"github.com/wahid-wex/e-commerce-backend/data/models"
	logging "github.com/wahid-wex/e-commerce-backend/logs"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type SellerService struct {
	logger       logging.Logger
	cfg          *config.Config
	otpService   *OtpService
	tokenService *TokenService
	database     *gorm.DB
}

func NewSellerService(cfg *config.Config) *SellerService {
	database := db.GetDb()
	logger := logging.NewLogger(cfg)
	return &SellerService{
		cfg:          cfg,
		database:     database,
		logger:       logger,
		otpService:   NewOtpService(cfg),
		tokenService: NewTokenService(cfg),
	}
}

func (s *SellerService) LoginSellerByUsername(req *dto.LoginByUsernameRequest) (*dto.TokenDetail, error) {
	var seller models.Seller
	err := s.database.
		Model(&models.Seller{}).
		Where("username = ?", req.Username).
		Preload("UserRoles.Role").
		First(&seller).Error

	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(seller.Password), []byte(req.Password))
	if err != nil {
		return nil, &error_handler.ServiceError{
			EndUserMessage:   error_handler.WrongPassword,
			TechnicalMessage: "Password comparison failed",
			Err:              err,
		}
	}
	tokenDto := tokenDto{UserId: seller.ID, Name: seller.StoreName,
		Email: seller.Email, MobileNumber: seller.Phone, Username: seller.Username}

	if len(seller.UserRoles) > 0 {
		for _, ur := range seller.UserRoles {
			tokenDto.Roles = append(tokenDto.Roles, ur.Role.Name)
		}
	}

	token, err := s.tokenService.GenerateToken(&tokenDto)
	if err != nil {
		return nil, err
	}
	return token, nil

}

func (s *SellerService) RegisterSellerByUsername(req *dto.RegisterSellerByUsernameRequest) error {
	u := models.Seller{
		Username:    req.Username,
		Email:       req.Email,
		Password:    req.Password,
		StoreName:   req.StoreName,
		NationalID:  req.NationalID,
		Address:     req.Address,
		Phone:       req.Phone,
		Description: req.Description,
		Logo:        req.Logo,
	}

	exists, err := s.existsByEmail(req.Email)
	if err != nil {
		return err
	}
	if exists {
		return &error_handler.ServiceError{EndUserMessage: error_handler.EmailExists}
	}
	exists, err = s.existsByUsername(req.Username)
	if err != nil {
		return err
	}
	if exists {
		return &error_handler.ServiceError{EndUserMessage: error_handler.UsernameExists}
	}

	bp := []byte(req.Password)
	hp, err := bcrypt.GenerateFromPassword(bp, bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error(logging.General, logging.HashPassword, err.Error(), nil)
		return err
	}
	u.Password = string(hp)
	roleId, err := s.getDefaultRole()
	if err != nil {
		s.logger.Error(logging.Postgres, logging.DefaultRoleNotFound, err.Error(), nil)
		return err
	}

	tx := s.database.Begin()
	err = tx.Create(&u).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error(logging.Postgres, logging.Rollback, err.Error(), nil)
		return err
	}
	err = tx.Exec("INSERT INTO user_roles (role_id, seller_id, created_at, updated_at) VALUES (?, ?, NOW(), NOW())", roleId, u.ID).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error(logging.Postgres, logging.Rollback, err.Error(), nil)
		return err
	}
	tx.Commit()
	return nil

}

func (s *SellerService) RegisterLoginSellerByMobileNumber(req *dto.RegisterLoginByMobileRequest) (*dto.TokenDetail, error) {
	err := s.otpService.ValidateOtp(req.MobileNumber, req.Otp)
	if err != nil {
		return nil, err
	}
	exists, err := s.existsSellerByMobileNumber(req.MobileNumber)
	if err != nil {
		return nil, err
	}

	u := models.Seller{Phone: req.MobileNumber, Username: req.MobileNumber}

	if exists {
		return s.loginSellerByMobileNumber(u.Username)
	}

	// Register and login
	bp := []byte(common.GeneratePassword())
	hp, err := bcrypt.GenerateFromPassword(bp, bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error(logging.General, logging.HashPassword, err.Error(), nil)
		return nil, err
	}
	u.Password = string(hp)
	roleId, err := s.getDefaultRole()
	if err != nil {
		s.logger.Error(logging.Postgres, logging.DefaultRoleNotFound, err.Error(), nil)
		return nil, err
	}

	tx := s.database.Begin()
	err = tx.Create(&u).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error(logging.Postgres, logging.Rollback, err.Error(), nil)
		return nil, err
	}
	err = tx.Create(&models.UserRole{RoleID: roleId, SellerID: &u.ID}).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error(logging.Postgres, logging.Rollback, err.Error(), nil)
		return nil, err
	}
	tx.Commit()

	return s.loginSellerByMobileNumber(u.Username)

}

func (s *SellerService) loginSellerByMobileNumber(phone string) (*dto.TokenDetail, error) {
	var seller models.Seller
	err := s.database.
		Model(&models.Seller{}).
		Where("phone = ?", phone).
		Preload("UserRoles", func(tx *gorm.DB) *gorm.DB {
			return tx.Preload("Role")
		}).
		Find(&seller).Error
	if err != nil {
		return nil, err
	}
	tokenDto := tokenDto{UserId: seller.ID, Name: seller.StoreName,
		Email: seller.Email, MobileNumber: seller.Phone, Username: seller.Username}

	if len(seller.UserRoles) > 0 {
		for _, ur := range seller.UserRoles {
			tokenDto.Roles = append(tokenDto.Roles, ur.Role.Name)
		}
	}

	token, err := s.tokenService.GenerateToken(&tokenDto)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (s *SellerService) SendOtpToSeller(req *dto.GetOtpRequest) error {
	otp := common.GenerateOtp()
	err := s.otpService.SetOtp(req.MobileNumber, otp)
	if err != nil {
		return err
	}
	return nil
}

func (s *SellerService) existsByEmail(email string) (bool, error) {
	var exists bool
	if err := s.database.Model(&models.Seller{}).
		Select(countFilterExp).
		Where("email = ?", email).
		Find(&exists).
		Error; err != nil {
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		return false, err
	}
	return exists, nil
}

func (s *SellerService) existsByUsername(username string) (bool, error) {
	var exists bool
	if err := s.database.Model(&models.Seller{}).
		Select(countFilterExp).
		Where("username = ?", username).
		Find(&exists).
		Error; err != nil {
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		return false, err
	}
	return exists, nil
}

func (s *SellerService) existsSellerByMobileNumber(mobileNumber string) (bool, error) {
	var exists bool
	if err := s.database.Model(&models.Seller{}).
		Select(countFilterExp).
		Where("phone = ?", mobileNumber).
		Find(&exists).
		Error; err != nil {
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		return false, err
	}
	return exists, nil
}

func (s *SellerService) getDefaultRole() (roleId uint, err error) {

	if err = s.database.Model(&models.Role{}).
		Select("id").
		Where("name = ?", "seller").
		First(&roleId).Error; err != nil {
		return 0, err
	}
	return roleId, nil
}

func (s *SellerService) RefreshToken(r string) (*dto.TokenDetail, error) {
	tokenDto, err := s.tokenService.RefreshToken(r)
	if err != nil {
		s.logger.Error(logging.General, logging.InvalidToken, err.Error(), nil)
		return nil, err
	}
	return tokenDto, nil
}
