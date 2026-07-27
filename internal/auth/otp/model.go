package otp

type OTPItem struct {
	Email string
	Code string
}

var RedisKeyPrefix string = "otp:"