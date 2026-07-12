package user

import (
	"net/http"
	"encoding/json"
	"fmt"
	"time"
)


type Handler struct{
	svc *Service
}

func NewHandler(svc *Service) *Handler{
	return &Handler{
		svc: svc,
	}
}



func (h *Handler) SignUpHandler(w http.ResponseWriter, r *http.Request){
	var SignupReq SignUpRequest

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&SignupReq)
	if err != nil {
		// return code http
		return
	}
	fmt.Println(SignupReq)

	user, err := h.svc.Signup(SignupReq.Email, SignupReq.Password)
	if err != nil {
    // return appropriate HTTP response
    	return
	}

	var userResponse SignUpResponse

	userResponse.id = user.ID
	userResponse.Email = user.Email



	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	encoder := json.NewEncoder(w)
	encoder.Encode(userResponse)


}