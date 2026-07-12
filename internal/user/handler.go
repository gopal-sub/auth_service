package user

import (
	"encoding/json"
	"fmt"
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

	decoder := json.NewDecoder(r.Body)
	//reject req if there are junk fields
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&SignupReq)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid request body",
		})
		return
	}
	fmt.Println(SignupReq)

	user, err := h.svc.Signup(SignupReq.Email, SignupReq.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{
			"error": `user with email already exists`,
			"email": SignupReq.Email,
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
	
	json.NewEncoder(w).Encode(token)

}