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

const countFilterExp string = "count(*) > 0"

type CustomerService struct {
	logger       logging.Logger
	cfg          *config.Config
	otpService   *OtpService
	tokenService *TokenService
	database     *gorm.DB
}

func NewCustomerService(cfg *config.Config) *CustomerService {
	database := db.GetDb()
	logger := logging.NewLogger(cfg)
	return &CustomerService{
		cfg:          cfg,
		database:     database,
		logger:       logger,
		otpService:   NewOtpService(cfg),
		tokenService: NewTokenService(cfg),
	}
}

// Login by username
func (s *CustomerService) LoginByUsername(req *dto.LoginByUsernameRequest) (*dto.TokenDetail, error) {
	var customer models.Customer
	err := s.database.
		Model(&models.Customer{}).
		Where("username = ?", req.Username).
		Preload("UserRoles.Role").
		First(&customer).Error
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(req.Password))
	if err != nil {
		return nil, &error_handler.ServiceError{EndUserMessage: error_handler.WrongPassword}
	}
	tokenDto := tokenDto{UserId: customer.ID, Name: customer.FirstName + customer.LastName,
		Email: customer.Email, MobileNumber: customer.Phone}

	if len(customer.UserRoles) > 0 {
		for _, ur := range customer.UserRoles {
			tokenDto.Roles = append(tokenDto.Roles, ur.Role.Name)
		}
	}

	token, err := s.tokenService.GenerateToken(&tokenDto)
	if err != nil {
		return nil, err
	}
	return token, nil

}

func (s *CustomerService) RefreshToken(r string) (*dto.TokenDetail, error) {
	tokenDto, err := s.tokenService.RefreshToken(r)
	if err != nil {
		s.logger.Error(logging.General, logging.InvalidToken, err.Error(), nil)
		return nil, err
	}
	return tokenDto, nil
}

// Register by username
func (s *CustomerService) RegisterByUsername(req *dto.RegisterCustomerByUsernameRequest) error {
	u := models.Customer{
		Username:        req.Username,
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		Email:           req.Email,
		PostalCode:      req.PostalCode,
		Phone:           req.Phone,
		CardNumber:      req.CardNumber,
		ShippingAddress: req.ShippingAddress,
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
	roleId, err := s.getCustomerRoleId()
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
	err = tx.Create(&models.UserRole{RoleID: roleId, CustomerID: &u.ID}).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error(logging.Postgres, logging.Rollback, err.Error(), nil)
		return err
	}
	tx.Commit()
	return nil

}

func (s *CustomerService) RegisterLoginByMobileNumber(req *dto.RegisterLoginByMobileRequest) (*dto.TokenDetail, error) {
	err := s.otpService.ValidateOtp(req.MobileNumber, req.Otp)
	if err != nil {
		return nil, err
	}
	exists, err := s.existsByMobileNumber(req.MobileNumber)
	u := models.Customer{Phone: req.MobileNumber}
	if exists {
		return s.loginByMobileNumber(req.MobileNumber)
	}

	// Register and login
	bp := []byte(common.GeneratePassword())
	hp, err := bcrypt.GenerateFromPassword(bp, bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error(logging.General, logging.HashPassword, err.Error(), nil)
		return nil, err
	}
	u.Password = string(hp)
	roleId, err := s.getCustomerRoleId()
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
	err = tx.Create(&models.UserRole{RoleID: roleId, CustomerID: &u.ID}).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error(logging.Postgres, logging.Rollback, err.Error(), nil)
		return nil, err
	}
	tx.Commit()

	return s.loginByMobileNumber(u.Phone)

}

func (s *CustomerService) loginByMobileNumber(mobile string) (*dto.TokenDetail, error) {
	var customer models.Customer
	err := s.database.
		Model(&models.Customer{}).
		Where("Phone = ?", mobile).
		Preload("UserRoles.Role").
		Find(&customer).Error
	if err != nil {
		return nil, err
	}
	tokenDto := tokenDto{UserId: customer.ID, Name: customer.FirstName + customer.LastName,
		Email: customer.Email, MobileNumber: customer.Phone}

	if len(customer.UserRoles) > 0 {
		for _, ur := range customer.UserRoles {
			tokenDto.Roles = append(tokenDto.Roles, ur.Role.Name)
		}
	}

	token, err := s.tokenService.GenerateToken(&tokenDto)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (s *CustomerService) SendOtp(req *dto.GetOtpRequest) error {
	otp := common.GenerateOtp()
	err := s.otpService.SetOtp(req.MobileNumber, otp)
	if err != nil {
		return err
	}
	return nil
}

func (s *CustomerService) existsByEmail(email string) (bool, error) {
	var exists bool
	if err := s.database.Model(&models.Customer{}).
		Select(countFilterExp).
		Where("email = ?", email).
		Find(&exists).
		Error; err != nil {
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		return false, err
	}
	return exists, nil
}

func (s *CustomerService) existsByUsername(username string) (bool, error) {
	var exists bool
	if err := s.database.Model(&models.Customer{}).
		Select(countFilterExp).
		Where("username = ?", username).
		Find(&exists).
		Error; err != nil {
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		return false, err
	}
	return exists, nil
}

func (s *CustomerService) existsByMobileNumber(mobileNumber string) (bool, error) {
	var exists bool
	if err := s.database.Model(&models.Customer{}).
		Select(countFilterExp).
		Where("phone = ?", mobileNumber).
		Find(&exists).
		Error; err != nil {
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		return false, err
	}
	return exists, nil
}

func (s *CustomerService) getCustomerRoleId() (roleId uint, err error) {

	if err = s.database.Model(&models.Role{}).
		Select("id").
		Where("name = ?", "customer").
		First(&roleId).Error; err != nil {
		return 0, err
	}
	return roleId, nil
}
