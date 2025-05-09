package services

import (
	"github.com/wahid-wex/e-commerce-backend/api/dto"
	"github.com/wahid-wex/e-commerce-backend/api/error_handler"
	"github.com/wahid-wex/e-commerce-backend/config"
	"github.com/wahid-wex/e-commerce-backend/constants"
	logging "github.com/wahid-wex/e-commerce-backend/logs"
	"time"

	"github.com/golang-jwt/jwt"
)

type TokenService struct {
	logger logging.Logger
	cfg    *config.Config
}

type tokenDto struct {
	UserId       uint
	Name         string
	Username     string
	MobileNumber string
	Email        string
	Roles        []string
}

func NewTokenService(cfg *config.Config) *TokenService {
	logger := logging.NewLogger(cfg)
	return &TokenService{
		cfg:    cfg,
		logger: logger,
	}
}

func (s *TokenService) GenerateToken(token *tokenDto) (*dto.TokenDetail, error) {
	td := &dto.TokenDetail{}
	td.AccessTokenExpireTime = time.Now().Add(s.cfg.JWT.AccessTokenExpireDuration * time.Minute).Unix()
	td.RefreshTokenExpireTime = time.Now().
	    Add((s.cfg.JWT.RefreshTokenExpireDuration + s.cfg.JWT.AccessTokenExpireDuration) * time.Minute).
		Unix()

	atc := jwt.MapClaims{}

	atc[constants.UserIdKey] = token.UserId
	atc[constants.NameKey] = token.Name
	atc[constants.UsernameKey] = token.Username
	atc[constants.EmailKey] = token.Email
	atc[constants.MobileNumberKey] = token.MobileNumber
	atc[constants.RolesKey] = token.Roles
	atc[constants.ExpireTimeKey] = td.AccessTokenExpireTime

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atc)

	var err error
	td.AccessToken, err = at.SignedString([]byte(s.cfg.JWT.Secret))

	if err != nil {
		return nil, err
	}

	rtc := jwt.MapClaims{}

	rtc[constants.UserIdKey] = token.UserId
	rtc[constants.ExpireTimeKey] = td.RefreshTokenExpireTime

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtc)

	td.RefreshToken, err = rt.SignedString([]byte(s.cfg.JWT.RefreshSecret))

	if err != nil {
		return nil, err
	}

	return td, nil
}

func (s *TokenService) VerifyToken(token string) (*jwt.Token, error) {
	at, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, &error_handler.ServiceError{EndUserMessage: error_handler.UnExpectedError}
		}
		return []byte(s.cfg.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	return at, nil
}

func (s *TokenService) GetClaims(token string) (claimMap map[string]interface{}, err error) {
	claimMap = map[string]interface{}{}

	verifyToken, err := s.VerifyToken(token)
	if err != nil {
		return nil, err
	}
	claims, ok := verifyToken.Claims.(jwt.MapClaims)
	if ok && verifyToken.Valid {
		for k, v := range claims {
			claimMap[k] = v
		}
		return claimMap, nil
	}
	return nil, &error_handler.ServiceError{EndUserMessage: error_handler.ClaimsNotFound}
}

func (s *TokenService) RefreshToken(refreshToken string) (*dto.TokenDetail, error) {
	// Verify the refresh token
	rt, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, &error_handler.ServiceError{EndUserMessage: error_handler.UnExpectedError}
		}
		return []byte(s.cfg.JWT.RefreshSecret), nil
	})
	if err != nil {
		return nil, &error_handler.ServiceError{EndUserMessage: error_handler.InvalidToken}
	}

	// Check claims
	claims, ok := rt.Claims.(jwt.MapClaims)
	if !ok || !rt.Valid {
		return nil, &error_handler.ServiceError{EndUserMessage: error_handler.InvalidToken}
	}

	// Check expiration time
	expTime, ok := claims[constants.ExpireTimeKey].(float64)
	if !ok {
		return nil, &error_handler.ServiceError{EndUserMessage: error_handler.InvalidToken}
	}

	// Convert to time.Time
	expTimeUnix := int64(expTime)
	if time.Now().Unix() > expTimeUnix {
		return nil, &error_handler.ServiceError{EndUserMessage: error_handler.TokenExpired}
	}

	// Extract user_id from claims
	userIdFloat, ok := claims[constants.UserIdKey].(float64)
	if !ok {
		return nil, &error_handler.ServiceError{EndUserMessage: error_handler.InvalidToken}
	}
	userId := uint(userIdFloat)

	// Now you have the userId, create a new access+refresh token
	tokenDto := tokenDto{
		UserId: userId,
	}

	newTokens, err := s.GenerateToken(&tokenDto)
	if err != nil {
		return nil, err
	}

	return newTokens, nil
}
