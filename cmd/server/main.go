package main

import (
	"fmt"
	"gopal-sub/auth_service/internal/auth"
	"gopal-sub/auth_service/internal/auth/oauth"
	"gopal-sub/auth_service/internal/auth/oauth/google"
	"gopal-sub/auth_service/internal/auth/otp"
	"gopal-sub/auth_service/internal/database"
	rlimiter "gopal-sub/auth_service/internal/ratelimiter"
	"gopal-sub/auth_service/internal/redis"
	"gopal-sub/auth_service/internal/user"
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

	//oauth init
	var googleProvider oauth.OauthProvider = "google"
	googelProvider := google.NewGoogleProvider()
	oauthsvc := oauth.NewService(googelProvider, userservice)
	url := oauthsvc.GetAuthLoginURL(googleProvider)
	fmt.Println(url)

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
	router.HandleFunc("POST /resetPassword", userhandler.PasswordReset)

	// router.HandleFunc("GET /oauth", userhandler.OauthLogin)
	// router.HandleFunc("GET /callback", userhandler.OauthCallback)

	router.Handle("POST /profile", auth.Middleware(http.HandlerFunc(userhandler.Profile)))

	userServer := http.Server{
		Addr:    ":3000",
		Handler: router,
	}
	userServer.ListenAndServe()

	
	fmt.Println("connected")
	defer db.Close()
}