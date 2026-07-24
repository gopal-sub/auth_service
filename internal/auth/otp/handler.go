package otp

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	svc *Service
}


func NewHandler (svc *Service) *Handler{
	return &Handler{
		svc: svc,
	}
}



func (s *Handler) SendOTP(w http.ResponseWriter, r *http.Request) {
	var SendOTP SendOTP

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	decoder.Decode(&SendOTP)

	err := s.svc.SendOTP(SendOTP.Email)
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