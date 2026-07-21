package auth

import (
	"os"
	"github.com/golang-jwt/jwt/v5"
	"time"
)


func CreateJWT(email string) (string, error){
	secret := os.Getenv("JWT_SECRET");

	var jwtToken = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	"email": email,
	"nbf": time.Now().Unix(),
	});

	

	tokenString, err := jwtToken.SignedString([]byte(secret))
	
	if err != nil {
		return "", err
	}
	return tokenString, nil
}