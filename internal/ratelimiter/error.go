package ratelimiter

import (
	"errors"
)

var OTPAttemptsExeErr = errors.New("number of otp password reset attempts exceeded")