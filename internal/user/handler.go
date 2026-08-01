package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopal-sub/auth_service/internal/httputils"
	"gopal-sub/auth_service/internal/validator"
	"net/http"
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

	err := httputils.DecodeAndWriteError(w, r, &SignupReq)
	if err != nil {
		return
	}

	err = validator.Email(SignupReq.Email)
	if errors.Is(err, validator.EmailRequiredErr){
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": `no email was provided`,
		})
    	return
	}
	if errors.Is(err, validator.InvalidEmailErr){
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": `email provided was invalid`,
		})
    	return
	}


	user, err := h.svc.Signup(SignupReq.Email, SignupReq.Password)
	fmt.Println(err)
	if errors.Is(err, UserExistsConflict){
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{
			"error": `user with email already exists`,
			"email": SignupReq.Email,
		})
    	return
	}
	if err != nil{
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": `There was a server error`,
		})
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

func (h *Handler) SigninHandler(w http.ResponseWriter, r *http.Request){
	var signinRequest SigninRequest


	decorder := json.NewDecoder(r.Body)
	decorder.DisallowUnknownFields()
	err := decorder.Decode(&signinRequest)
	if err != nil{
		return
	}
	token, err := h.svc.Signin(signinRequest.Email, signinRequest.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	var signinResponse SigninResponse

	signinResponse.Token = token
	
	json.NewEncoder(w).Encode(signinResponse)

}