package auth

import (
	"errors"
)

var expiredOTPErr = errors.New("OTP has expired")
var invalidOTPErr = errors.New("OTP is invalid")