package user

import (
	"database/sql"
	"errors"
	"gopal-sub/auth_service/internal/auth"
)

type Service struct{
	repo UserRepository
}

func NewService(repo UserRepository) *Service{
	return &Service{
		repo: repo,
	}
}

func (s *Service) Signup(email string, password string) (User, error){
	_, err := s.repo.FindUserByEmail(email)

	if errors.Is(err, sql.ErrNoRows){
		hashedPassword, err := CreateHash(password)
		if err != nil{
			return User{}, err
		}

		newUser := NewUser(email, hashedPassword)
		err = s.repo.Create(newUser)
		
		
		if err != nil {
			return User{}, err
		}
		return newUser, nil
	}
	if err == nil{
		return User{}, UserExistsConflict
	}

	return User{}, err



}


func (s *Service) Signin(email string, password string) (token string, err error){
	user, err := s.repo.FindUserByEmail(email)
	if errors.Is(err, sql.ErrNoRows) {
		// tasks -> if  user does not exist 
		return "", ErrInvalidCredentials
	}

	if err != nil{
		return "", err
	}

	err = CompareHashWithPass(user.PasswordHash, password)


	if err != nil {
		return "", ErrInvalidCredentials
	}

	
	tokenString, err := auth.CreateJWT(user.Email)

	if err != nil {
		return "", err
	}





	//make jwt imp
	return tokenString, nil
	
}

func (s *Service) UserExists(email string) bool {
	_, err := s.repo.FindUserByEmail(email)
	if errors.Is(err, sql.ErrNoRows){
		return false
	}
	return true
}

func (s *Service) IsUserAccountVerified(email string) (bool, error) {
	user, err := s.repo.FindUserByEmail(email)
	if errors.Is(err, sql.ErrNoRows){
		return false, UserDoesNotExistErr
	}
	if(user.IsVerified){
		return true, nil
	}
	return false, nil
}

func (s *Service) VerifyAnAccount(email string)(accountVerified bool, err error){
	accountVer, err := s.repo.VerifyAccount(email)
	if errors.Is(err, UserDoesNotExistErr){
		return false, UserDoesNotExistErr
	}
	if err != nil {
		return false, err
	}
	if(accountVer == true){
		return true, nil
	}
	return false, nil
}
	

