package user

import (
	"golang.org/x/crypto/bcrypt"
)


func CreateHash(password string)(string, error){
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	hashedPassword := string(hash)
	return hashedPassword, nil
}

func CompareHashWithPass(hash string, passwordOriginal string) (error){
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwordOriginal))
	return err
}