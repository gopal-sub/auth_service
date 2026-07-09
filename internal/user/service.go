package user

import (
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type Service struct{
	repo *Repository
}

func NewService(repo *Repository) *Service{
	return &Service{
		repo: repo,
	}
}

func (s *Service) Signup(email string, password string) (User, error){
	_, err := s.repo.FindUserByEmail(email)

	if errors.Is(err, sql.ErrNoRows){
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return User{}, err
		}

		hashedPassword := string(hash)

		newUser := NewUser(email, hashedPassword)
		err = s.repo.Create(newUser)
		
		
		if err != nil {
			return User{}, err
		}
		return newUser, nil
	}
	if err != nil{
		return User{}, err
	}

	return User{}, errors.New("email already exists")



}

