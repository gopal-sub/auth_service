package user

import (
	"encoding/json"
	"errors"
	"gopal-sub/auth_service/internal/auth/oauth"
	"gopal-sub/auth_service/internal/auth/oauth/google"
	"gopal-sub/auth_service/internal/httputils"
	"gopal-sub/auth_service/internal/validator"
	"net/http"
	// "net/url"
	// "os"
	// "strings"
	// "github.com/golang-jwt/jwt/v5"
)




type Handler struct{
	svc UserService
}

func NewHandler(svc UserService) *Handler{
	return &Handler{
		svc: svc,
	}
}



func (h *Handler) SignUpHandler(w http.ResponseWriter, r *http.Request){
	var SignupReq SignUpRequest

	queryValues := r.URL.Query()

	if len(queryValues) > 0 {
		val, ok := queryValues["oauthKind"]
		if ok {
			if val[0]== "google"{
				googleProvider := google.NewGoogleProvider()
				handler := oauth.NewHandler(*googleProvider)
				handler.OauthLogin(w, r)
				return 
			}
		}
		
		
	}
	

	err := httputils.DecodeAndWriteError(w, r, &SignupReq)
	if err != nil {
		return
	}

	if SignupReq.Password == ""{
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": `no password was provided`,
		})
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

func (h *Handler) PasswordReset(w http.ResponseWriter, r *http.Request) {
	var userPasswordUpdate UpdateUserRequest
	decorder := json.NewDecoder(r.Body)
	decorder.DisallowUnknownFields()
	err := decorder.Decode(&userPasswordUpdate)
	if err != nil{
		return
	}
	if userPasswordUpdate.OldPassword == userPasswordUpdate.NewPassword{
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Bad request",
		})
		return
	}

	_, err = h.svc.Signin(userPasswordUpdate.Email, userPasswordUpdate.OldPassword)
	if errors.Is(err, ErrInvalidCredentials){
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid password",
		})
		return
	}
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "server error",
		})
		return
	}
	err = h.svc.PasswordReset(userPasswordUpdate.Email, userPasswordUpdate.NewPassword)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "password reset failed",
		})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"success": "password reset success",
	})
}

func (h *Handler) Profile(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"User": "welcome to profile",
	})
}
	


