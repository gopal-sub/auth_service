package otp

type OTPItem struct {
	Email string
	Code string
}

type SendOTP struct {
	Email string
}

type VerifyOTP struct {
	Email string
	Code string
}

var RedisKeyPrefix string = "otp:"