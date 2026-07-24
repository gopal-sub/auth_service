package otp

type Service struct{
	repo OTPRepository
	smtp SMTPEmailSender
}

func NewOTPService(repo OTPRepository) *Service{
	return &Service{
		repo: repo,
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

	//send mail here
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