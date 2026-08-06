package auth

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)



func ValidateToken(tokenString string)(bool, error){

	secret := os.Getenv("JWT_SECRET");


	parsedToken ,err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error){
		return []byte(secret), nil
	})

	if err != nil{
		return false, err
	}
	
	if parsedToken.Valid == false {
		return false, nil
	}
	return true, nil
}

func Middleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "No token was passed",
			})
			return 
		}
		valid, _ := ValidateToken(token)
		if !valid {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "bad token",
			})
			return 
		}
		next.ServeHTTP(w, r)

	})
}