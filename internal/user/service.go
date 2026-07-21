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
	if err != nil{
		return User{}, dberr
	}

	return User{}, userExistsConflict



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

