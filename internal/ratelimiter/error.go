package ratelimiter

import (
	"errors"
)

var otpAttemptsExeErr = errors.New("number of otp password reset attempts exceeded")