package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopal-sub/auth_service/internal/httputils"
	"gopal-sub/auth_service/internal/validator"
	"net/http"
	"net/url"
	"os"
	"strings"
	"io"
	"github.com/golang-jwt/jwt/v5"
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


func (h *Handler) OauthLogin(w http.ResponseWriter, r *http.Request){

	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	// 
	redirectURI := "http://localhost:3000/callback"


	params := url.Values{}
	params.Add("client_id", clientID)
	params.Add("redirect_uri", redirectURI)
	params.Add("scope", "openid profile email")
	params.Add("response_type", "code")

	authorisationEndPoint := "https://accounts.google.com/o/oauth2/v2/auth?" + params.Encode()


	http.Redirect(w, r, authorisationEndPoint, http.StatusMovedPermanently)


}

func (h *Handler) OauthCallback(w http.ResponseWriter, r *http.Request){
	
	queryParams := r.URL.Query()

	code := queryParams.Get("code")
	err := queryParams.Get("error")
	if err != "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "shit went down",
		})
	}
	tokenEndpoint := "https://oauth2.googleapis.com/token"
	redirectURI := "http://localhost:3000/callback"

	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	clientID := os.Getenv("GOOGLE_CLIENT_ID")

	data := url.Values{}

	data.Add("client_id", clientID)
	data.Add("client_secret", clientSecret)
	data.Add("code", code)
	data.Add("grant_type", "authorization_code")
	data.Add("redirect_uri", redirectURI)

	res, _ := http.Post(tokenEndpoint, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))

	body, _ := io.ReadAll(res.Body)
	type id_token struct{
		AccessToken string `json:"access_token"`  
		ExpiresIn int `json:"expires_in"` 
		Score string `json:"scope"` 
		TokenType string `json:"token_type"` 
		IdToken string `json:"id_token"` 
	}

	id_token_get := id_token{}

	json.Unmarshal(body, &id_token_get)

	claims := jwt.MapClaims{}

	token, _, _ := jwt.NewParser().ParseUnverified(id_token_get.IdToken, claims)
	
	fmt.Println(claims["email"])
	fmt.Println(token)



}
	


