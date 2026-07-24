package main

import (

	"fmt"
	"gopal-sub/auth_service/internal/database"
	"gopal-sub/auth_service/internal/redis"
	"gopal-sub/auth_service/internal/user"
	"gopal-sub/auth_service/internal/auth/otp"
	"log"
	"net/http"


	"github.com/joho/godotenv"
)



func main(){
	router := http.NewServeMux()


	err := godotenv.Load()
	if err != nil{
		log.Fatal("env did not load")
	}

	db, err := database.New()
	
	

	if err != nil {
		log.Fatalf(`database error  wowowo  %v`, err)
	}
	rdb, err := redis.New()
	if err != nil {
		log.Fatalf(`redis error  wowowo  %v`, err)
	}
	defer rdb.Close()

	
	
	newRepo := user.NewRepository(db)
	service := user.NewService(newRepo)
	handler := user.NewHandler(service)

	otpRepo := otp.NewOTPRepo(rdb)
	otpService := otp.NewOTPService(otpRepo)
	otpHandler := otp.NewHandler(otpService)


	router.HandleFunc("POST /otp/send", otpHandler.SendOTP)
	router.HandleFunc("POST /otp/verify", otpHandler.VerifyOTP)




	router.HandleFunc("POST /user", handler.SignUpHandler)
	router.HandleFunc("POST /gettoken", handler.SigninHandler)

	userServer := http.Server{
		Addr:    ":3000",
		Handler: router,
	}
	userServer.ListenAndServe()

	
	fmt.Println("connected")
	defer db.Close()
}