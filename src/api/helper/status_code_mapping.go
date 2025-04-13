package helper

import (
	"github.com/wahid-wex/e-commerce-backend/api/error_handler"
	"net/http"
)

var StatusCodeMapping = map[string]int{

	// OTP
	error_handler.OptExists:   409,
	error_handler.OtpUsed:     409,
	error_handler.OtpNotValid: 400,

	// User
	error_handler.EmailExists:      409,
	error_handler.UsernameExists:   409,
	error_handler.RecordNotFound:   404,
	error_handler.PermissionDenied: 403,
}

func TranslateErrorToStatusCode(err error) int {
	value, ok := StatusCodeMapping[err.Error()]
	if !ok {
		return http.StatusInternalServerError
	}
	return value
}
