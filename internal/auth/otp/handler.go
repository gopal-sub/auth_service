package otp

import (
	"encoding/json"
	"errors"
	"gopal-sub/auth_service/internal/httputils"
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



func (h *Handler) PasswordResetSendOTP(w http.ResponseWriter, r *http.Request) {
	var SendOTP SendOTP


	err := httputils.DecodeJSONMiddleware(w, r, &SendOTP)
	if err != nil {
		return
	}

	if !h.userSvc.UserExists(SendOTP.Email){
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "user does not exist",
		})
		return
	}

	//rate limiter kickes in
	err = h.rlimiter.RateLimit(SendOTP.Email)
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



	err = h.svc.SendOTP(SendOTP.Email)
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

func (h *Handler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var VerifyOTP VerifyOTP

	err := httputils.DecodeJSONMiddleware(w, r, &VerifyOTP)
	if err != nil {
		return
	}
	

	valid, err := h.svc.VerifyOTP(VerifyOTP.Email, VerifyOTP.Code)

	if errors.Is(err, invalidOTPErr) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusGone)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid otp",
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

	if valid {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "OTP valid",
		})
	}

}




func (h *Handler) AccountVerifySendOTP(w http.ResponseWriter, r *http.Request){
	var otpItem SendOTP

	err := httputils.DecodeJSONMiddleware(w, r, &otpItem)
	if err != nil {
		return
	}
	
	userAccountVerifier, err := h.userSvc.IsUserAccountVerified(otpItem.Email)
	if errors.Is(err, user.UserDoesNotExistErr){
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "user does not exist",
		})
		return
	}
	if !userAccountVerifier{
		//rate limiter kickes in
		err := h.rlimiter.RateLimit(otpItem.Email)
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
		err = h.svc.SendOTP(otpItem.Email)
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
		return

	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusConflict)
	json.NewEncoder(w).Encode(map[string]string{
		"error": "Account is already verified",
	})
	

}

func (h *Handler) AccountVerifyVerifyOTP(w http.ResponseWriter, r *http.Request) {
	var VerifyOTP VerifyOTP

	err := httputils.DecodeJSONMiddleware(w, r, &VerifyOTP)
	if err != nil {
		return
	}

	valid, err := h.svc.VerifyOTP(VerifyOTP.Email, VerifyOTP.Code)


	if errors.Is(err, invalidOTPErr) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusGone)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid otp",
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

	if valid {
		accountVerified, err := h.userSvc.VerifyAnAccount(VerifyOTP.Email)
		if errors.Is(err, user.UserDoesNotExistErr){
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"status": "user account does not exist",
			})
		}
		if err != nil{
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "server down",
			})
			return
		}
		if accountVerified{
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]string{
				"status": "Account has been verified successfully",
			})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "Could not verify account",
		})
	}

}