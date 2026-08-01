package main

import (

	"fmt"
	"gopal-sub/auth_service/internal/database"
	"gopal-sub/auth_service/internal/redis"
	"gopal-sub/auth_service/internal/user"
	"gopal-sub/auth_service/internal/auth/otp"
	rlimiter "gopal-sub/auth_service/internal/ratelimiter"
	"log"
	"net/http"
	"os"


	"github.com/joho/godotenv"
)



func main(){
	router := http.NewServeMux()


	err := godotenv.Load()
	if err != nil{
		log.Fatal("env did not load")
	}

	//postgres init

	db, err := database.New()
	
	

	if err != nil {
		log.Fatalf(`database error  wowowo  %v`, err)
	}

	//redis db init
	rdb, err := redis.New()
	if err != nil {
		log.Fatalf(`redis error  wowowo  %v`, err)
	}
	defer rdb.Close()

	
	// userRepo and userservice init
	userRepo := user.NewRepository(db)
	userservice := user.NewService(userRepo)
	userhandler := user.NewHandler(userservice)


	// ratelimiter Repo and svc init

	rateLimiterRepo := rlimiter.NewRepo(rdb)
	rateLimiterSvc := rlimiter.NewService(rateLimiterRepo)

	// otprepo and svc init || SMTP svc and repo init

	otpRepo := otp.NewOTPRepo(rdb)

	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	fromEmail := os.Getenv("SMTP_EMAIL")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	smtpEmailSender := otp.NewSMTPEmailSender(host, port, fromEmail, smtpPassword)


	otpService := otp.NewOTPService(otpRepo, smtpEmailSender)

	

	otpHandler := otp.NewHandler(otpService, rateLimiterSvc, userservice)


	router.HandleFunc("POST /pass-reset/otp/send", otpHandler.PasswordResetSendOTP)
	router.HandleFunc("POST /account-verify/otp/send", otpHandler.AccountVerifySendOTP)
	router.HandleFunc("POST /account-verify/otp/verify",otpHandler.AccountVerifyVerifyOTP)




	router.HandleFunc("POST /user", userhandler.SignUpHandler)
	router.HandleFunc("POST /gettoken", userhandler.SigninHandler)

	userServer := http.Server{
		Addr:    ":3000",
		Handler: router,
	}
	userServer.ListenAndServe()

	
	fmt.Println("connected")
	defer db.Close()
}