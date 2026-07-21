package auth

import (

	"os"
	"github.com/golang-jwt/jwt/v5"
)



func ValidateToken(tokenString string)(bool, error){

	secret := os.Getenv("JWT_SECRET");


	parsedToken ,err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error){
		return []byte(secret), nil
	})

	if parsedToken.Valid == false {
		return false, nil
	}
	if err != nil{
		return false, err
	}
	return true, nil
}