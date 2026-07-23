package auth

import (
	"time"
	"math/rand/v2"
	"strconv"
	"github.com/redis/go-redis/v9"
)

type OTPItem struct {
	Email string
	Code string
	ExpiresAt time.Time
}
//repo 
type Repository struct {
	repo *redis.Client
}

type OTPRepository interface {
	SaveOTP(item OTPItem) error
	GetOTP(email string) (OTPItem, error)
 	DeleteOTP(email string) error
}



//svc

type OTPService struct{
	repo OTPRepository
}

func NewOTPService(repo OTPRepository) *OTPService{
	return &OTPService{
		repo: repo,
	}
}

//svc fns


func (s *OTPService) SendOTP(email string) error{

	otpItem := OTPItem{
		Email: email,
		Code: GenerateOPT(),
		ExpiresAt: time.Now().Add(10*time.Minute),
	}
	
	if err := s.repo.SaveOTP(otpItem); err != nil {
		return err
	}
	//send mail here
	return nil
}
func (v *OTPService) VerifyOTP(email string, otp string) (bool, error){
	otpItem, err := v.repo.GetOTP(email)
	//this error code are from reddis
	if err != nil {
		return false, err
	}
	//otp errrors
	if !IsOTPExpired(otpItem){
		return false, expiredOTPErr
	}

	if otpItem.Code != otp{
		return false, invalidOTPErr
	}
	if err := v.repo.DeleteOTP(email); err !=nil {
		return false, nil
	}
	return true, nil


}

//helpers

func GenerateOPT() string{
	random_number := rand.IntN(999999) + 10000
	return strconv.Itoa(random_number)
}

func IsOTPExpired(otpitem OTPItem) (bool){
	t1 := time.Now()
	if t1.After(otpitem.ExpiresAt){
		return true
	}
	return false
}

// implement endpoint -> emailsender inteface-> send main 
