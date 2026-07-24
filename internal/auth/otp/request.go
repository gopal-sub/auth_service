package otp



type SendOTP struct {
	Email string
}

type VerifyOTP struct {
	Email string
	Code string
}