package otp

import (
	"fmt"

)

type Service struct{
	repo OTPRepository
	sender EmailSender
	
}

func NewOTPService(repo OTPRepository, sender EmailSender) *Service{
	return &Service{
		repo: repo,
		sender: sender,
	}
}

//svc fns


func (s *Service) SendOTP(email string) error{



	otpItem := OTPItem{
		Email: email,
		Code: GenerateOPT(),
	}
	
	if err := s.repo.SaveOTP(otpItem); err != nil {
		return err
	}
	subject := fmt.Sprintf("YOUR OTP TO RESET FOR ACCOUNT %s", email)
	body := fmt.Sprintf("OTP %s", otpItem.Code)

	//send mail here
	s.sender.SendEmail(email, subject, body)
	return nil
}
func (v *Service) VerifyOTP(email string, otp string) (bool, error){
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