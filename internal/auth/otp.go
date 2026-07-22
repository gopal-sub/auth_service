package auth

import (
	"math/rand/v2"
	"strconv"
	"time"
)

type OTPItem struct {
	email string
	otp string
	expiresAt time.Time
}

func (s *OTPItem) SaveOTP(){}
func (g *OTPItem) GetOTP(){}
func (d *OTPItem) DeleteOTP(){}

func GenerateOPT() string{
	random_number := rand.IntN(999999) + 10000
	return strconv.Itoa(random_number)
}


// implement endpoint -> emailsender inteface-> send main 
