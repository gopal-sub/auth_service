package auth

import (
	"time"
	"math/rand/v2"
	"strconv"
	"github.com/redis/go-redis/v9"
	"context"
)

type OTPItem struct {
	Email string
	Code string
}
//repo 
type Repository struct {
	repo *redis.Client
}

type OTPRepository interface {
	SaveOTP(item OTPItem) error
	GetOTP(email string) (string, error)
 	DeleteOTP(email string) error
}

func NewOTPRepo(client *redis.Client) *Repository{
	return &Repository{
		repo: client,
	}
}

func (s *Repository)SaveOTP(item OTPItem) error{
	ctx := context.Background()

	err := s.repo.Set(ctx, item.Email, item, 10*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}

func (g *Repository)GetOTP(email string) (string, error){
	ctx := context.Background()
	val, err := g.repo.Get(ctx, email).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (d *Repository)DeleteOTP(email string) error{
	ctx := context.Background()
	err := d.repo.Del(ctx, email).Err()
	if err != nil {
		return err
	}
	return nil
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
	}
	
	if err := s.repo.SaveOTP(otpItem); err != nil {
		return err
	}
	//send mail here
	return nil
}
func (v *OTPService) VerifyOTP(email string, otp string) (bool, error){
	code, err := v.repo.GetOTP(email)
	//this error code are from reddis
	if err != nil {
		return false, err
	}
	//otp errrors

	if code != otp{
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

// implement endpoint -> emailsender inteface-> send main 
