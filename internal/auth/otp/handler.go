package otp

import (
	"encoding/json"
	"errors"
	rlimiter "gopal-sub/auth_service/internal/ratelimiter"
	"gopal-sub/auth_service/internal/user"
	"net/http"
)

type Handler struct {
	svc *Service
	rlimiter *rlimiter.Service
	userSvc *user.Service
}


func NewHandler (svc *Service, rlimiter *rlimiter.Service, userSvc *user.Service) *Handler{
	return &Handler{
		svc: svc,
		rlimiter: rlimiter,
		userSvc: userSvc,
	}
}



func (s *Handler) SendOTP(w http.ResponseWriter, r *http.Request) {
	var SendOTP SendOTP

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	decoder.Decode(&SendOTP)

	if !s.userSvc.UserExists(SendOTP.Email){
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "user does not exist",
		})
		return
	}

	//rate limiter kickes in
	err := s.rlimiter.RateLimit(SendOTP.Email)
	if errors.Is(err, rlimiter.OTPAttemptsExeErr){
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "too many attenpts try after sometime",
		})
		return
	}
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "server down",
		})
		return
	}



	err = s.svc.SendOTP(SendOTP.Email)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "server down",
		})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "OTP sent",
	})


}

func (s *Handler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var VerifyOTP VerifyOTP

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	decoder.Decode(&VerifyOTP)

	valid, err := s.svc.VerifyOTP(VerifyOTP.Email, VerifyOTP.Code)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "server down",
		})
		return
	}

	if valid {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "OTP valid",
		})
	}

}